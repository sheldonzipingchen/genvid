package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/genvid/backend/internal/config"
	"github.com/genvid/backend/internal/model"
	"github.com/genvid/backend/internal/repository"
	"github.com/genvid/backend/internal/video"
	"github.com/genvid/backend/internal/zhipu"
	"github.com/genvid/backend/pkg/auth"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailExists         = errors.New("email already registered")
	ErrInsufficientCredits = errors.New("insufficient credits")
)

type AuthService struct {
	profileRepo *repository.ProfileRepository
	jwtService  *auth.JWTService
	cfg         *config.Config
}

func NewAuthService(profileRepo *repository.ProfileRepository, jwtService *auth.JWTService, cfg *config.Config) *AuthService {
	return &AuthService{
		profileRepo: profileRepo,
		jwtService:  jwtService,
		cfg:         cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	existing, _ := s.profileRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailExists
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	_ = passwordHash

	profile := &model.Profile{
		ID:    uuid.New().String(),
		Email: req.Email,
	}

	if req.FullName != "" {
		profile.FullName = &req.FullName
	}

	if err := s.profileRepo.Create(ctx, profile); err != nil {
		return nil, err
	}

	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(
		profile.ID,
		profile.Email,
		profile.SubscriptionTier,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(profile.ID)
	if err != nil {
		return nil, err
	}

	_ = s.profileRepo.UpdateLastLogin(ctx, profile.ID)

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         *profile,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	profile, err := s.profileRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.profileRepo.UpdateLastLogin(ctx, profile.ID); err != nil {
		return nil, err
	}

	accessToken, expiresIn, err := s.jwtService.GenerateAccessToken(
		profile.ID,
		profile.Email,
		profile.SubscriptionTier,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(profile.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         *profile,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	accessToken, expiresIn, err := s.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(profile.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         *profile,
	}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.Profile, error) {
	return s.profileRepo.GetByID(ctx, userID)
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) (*model.Profile, error) {
	profile, err := s.profileRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if fullName, ok := updates["full_name"].(string); ok {
		profile.FullName = &fullName
	}
	if companyName, ok := updates["company_name"].(string); ok {
		profile.CompanyName = &companyName
	}
	if preferredLang, ok := updates["preferred_language"].(string); ok {
		profile.PreferredLanguage = preferredLang
	}

	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *AuthService) ValidateAccessToken(token string) (*auth.Claims, error) {
	return s.jwtService.ValidateToken(token)
}

func (s *AuthService) CheckCredits(ctx context.Context, userID string) (bool, error) {
	profile, err := s.profileRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return profile.CreditsRemaining > 0 || profile.SubscriptionTier == "enterprise", nil
}

func (s *AuthService) UseCredit(ctx context.Context, userID string) error {
	return s.profileRepo.DecrementCredits(ctx, userID)
}

func (s *AuthService) RefundCredit(ctx context.Context, userID string) error {
	return s.profileRepo.IncrementCredits(ctx, userID, 1)
}

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	profileRepo *repository.ProfileRepository
	authService *AuthService
	zhipuClient *zhipu.Client
	cfg         *config.Config
}

func NewProjectService(projectRepo *repository.ProjectRepository, profileRepo *repository.ProfileRepository, authService *AuthService, zhipuClient *zhipu.Client, cfg *config.Config) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		profileRepo: profileRepo,
		authService: authService,
		zhipuClient: zhipuClient,
		cfg:         cfg,
	}
}

func (s *ProjectService) Create(ctx context.Context, userID string, req *model.CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		UserID:             userID,
		ProductName:        &req.ProductName,
		ProductDescription: req.ProductDescription,
		ProductURL:         req.ProductURL,
		ProductImageURL:    req.ProductImageURL,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id, userID string) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, repository.ErrUnauthorized
	}

	return project, nil
}

func (s *ProjectService) ListByUser(ctx context.Context, userID string, page, limit int) ([]model.Project, int, error) {
	offset := (page - 1) * limit
	return s.projectRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *ProjectService) Delete(ctx context.Context, id, userID string) error {
	return s.projectRepo.Delete(ctx, id, userID)
}

func (s *ProjectService) GenerateVideo(ctx context.Context, projectID, userID string, req *model.GenerateVideoRequest) (*model.Project, error) {
	hasCredits, err := s.authService.CheckCredits(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !hasCredits {
		return nil, ErrInsufficientCredits
	}

	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, repository.ErrUnauthorized
	}

	if err := s.authService.UseCredit(ctx, userID); err != nil {
		return nil, err
	}

	project.Script = &req.Script
	project.Language = req.Language
	if project.Language == "" {
		project.Language = "zh"
	}
	project.Format = model.VideoFormat(req.Format)
	project.VideoDuration = req.VideoDuration
	if project.VideoDuration == 0 {
		project.VideoDuration = 5
	}

	if err := s.projectRepo.UpdateStatus(ctx, projectID, model.ProjectStatusQueued, 0); err != nil {
		_ = s.authService.RefundCredit(ctx, userID)
		return nil, err
	}

	project.Status = model.ProjectStatusQueued

	go func() {
		bgCtx := context.Background()
		s.processVideoGeneration(bgCtx, project)
	}()

	return project, nil
}

func (s *ProjectService) processVideoGeneration(ctx context.Context, project *model.Project) {
	_ = s.projectRepo.UpdateStatus(ctx, project.ID, model.ProjectStatusProcessing, 5)

	duration := project.VideoDuration
	if duration == 0 {
		duration = 5
	}

	segments := 1
	if duration > 10 {
		segments = (duration + 9) / 10
	}

	prompt := ""
	if project.ProductName != nil {
		prompt = "Product: " + *project.ProductName + "\n\n"
	}
	if project.Script != nil {
		prompt += *project.Script
	}

	size := s.getVideoSize(string(project.Format))

	var imageURL string
	if project.ProductImageURL != nil && *project.ProductImageURL != "" {
		imageData, err := s.loadImageAsBase64(*project.ProductImageURL)
		if err == nil {
			imageURL = imageData
		}
	}

	var videoURLs []string
	var lastThumbnailURL string

	for i := 0; i < segments; i++ {
		progress := 10 + (i * 60 / segments)
		_ = s.projectRepo.UpdateStatus(ctx, project.ID, model.ProjectStatusProcessing, progress)

		segmentPrompt := prompt
		if segments > 1 {
			segmentPrompts := video.SplitScript(prompt, segments)
			if i < len(segmentPrompts) {
				segmentPrompt = segmentPrompts[i]
			}
		}

		req := zhipu.VideoGenerationRequest{
			Model:   s.cfg.External.Zhipu.Model,
			Prompt:  segmentPrompt,
			Quality: "speed",
			Size:    size,
		}

		if imageURL != "" {
			req.ImageURL = imageURL
			req.Prompt = "Strictly preserve the exact appearance of the product in the image: maintain identical shape, size, colors, textures, materials, branding, logos, labels, and all visual details. Do not modify, distort, or alter the product in any way. Only animate the scene around the product. " + segmentPrompt
		}

		resp, err := s.zhipuClient.GenerateVideo(req)
		if err != nil {
			s.handleVideoFailure(ctx, project, fmt.Sprintf("Segment %d failed: %s", i+1, err.Error()))
			return
		}

		taskProgress := progress + (30 / segments)
		_ = s.projectRepo.UpdateStatus(ctx, project.ID, model.ProjectStatusProcessing, taskProgress)

		result, err := s.zhipuClient.WaitForCompletion(resp.ID, 10*time.Minute)
		if err != nil {
			s.handleVideoFailure(ctx, project, fmt.Sprintf("Segment %d completion failed: %s", i+1, err.Error()))
			return
		}

		if result.VideoResult != nil && result.VideoResult.URL != "" {
			videoURLs = append(videoURLs, result.VideoResult.URL)
			if i == 0 && result.VideoResult.CoverURL != "" {
				lastThumbnailURL = result.VideoResult.CoverURL
			}
		}
	}

	_ = s.projectRepo.UpdateStatus(ctx, project.ID, model.ProjectStatusProcessing, 90)

	var finalVideoURL string
	if len(videoURLs) == 0 {
		s.handleVideoFailure(ctx, project, "No videos generated")
		return
	} else if len(videoURLs) == 1 {
		finalVideoURL = videoURLs[0]
	} else {
		mergedPath, err := s.mergeVideos(videoURLs, project.ID)
		if err != nil {
			finalVideoURL = videoURLs[0]
		} else {
			finalVideoURL = mergedPath
		}
	}

	if err := s.projectRepo.SetCompleted(ctx, project.ID, finalVideoURL, lastThumbnailURL); err != nil {
		_ = s.projectRepo.SetFailed(ctx, project.ID, err.Error())
		_ = s.authService.RefundCredit(ctx, project.UserID)
	}
}

func (s *ProjectService) mergeVideos(videoURLs []string, projectID string) (string, error) {
	if err := video.CheckFFmpeg(); err != nil {
		return "", err
	}

	merger := video.NewMerger("./temp_videos")
	outputPath := merger.GetOutputPath(projectID)

	finalPath, err := merger.MergeVideos(videoURLs, outputPath)
	if err != nil {
		return "", err
	}

	return "/temp_videos/" + filepath.Base(finalPath), nil
}

func (s *ProjectService) loadImageAsBase64(imagePath string) (string, error) {
	if strings.HasPrefix(imagePath, "data:image") {
		return imagePath, nil
	}

	filePath := strings.TrimPrefix(imagePath, "/uploads/")
	filePath = filepath.Join("uploads", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	mimeType := "image/jpeg"
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	}

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data)), nil
}

func (s *ProjectService) handleVideoFailure(ctx context.Context, project *model.Project, errMsg string) {
	_ = s.projectRepo.SetFailed(ctx, project.ID, errMsg)
	_ = s.authService.RefundCredit(ctx, project.UserID)
}

func (s *ProjectService) getVideoSize(format string) string {
	sizes := map[string]string{
		"9:16": "1080x1920",
		"1:1":  "1024x1024",
		"16:9": "1920x1080",
	}

	if size, ok := sizes[format]; ok {
		return size
	}
	return "1080x1920"
}
