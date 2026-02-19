package service

import (
	"context"
	"errors"
	"time"

	"github.com/genvid/backend/internal/config"
	"github.com/genvid/backend/internal/model"
	"github.com/genvid/backend/internal/repository"
	"github.com/genvid/backend/pkg/auth"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already registered")
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
}

func NewProjectService(projectRepo *repository.ProjectRepository, profileRepo *repository.ProfileRepository, authService *AuthService) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		profileRepo: profileRepo,
		authService: authService,
	}
}

func (s *ProjectService) Create(ctx context.Context, userID string, req *model.CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		UserID:             userID,
		ProductName:        &req.ProductName,
		ProductDescription: req.ProductDescription,
		ProductURL:         req.ProductURL,
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

	project.AvatarID = &req.AvatarID
	project.Script = &req.Script
	project.Language = req.Language
	project.Format = model.VideoFormat(req.Format)

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
	_ = s.projectRepo.UpdateStatus(ctx, project.ID, model.ProjectStatusProcessing, 10)

	time.Sleep(2 * time.Second)
	_ = s.projectRepo.UpdateStatus(ctx, project.ID, model.ProjectStatusProcessing, 50)

	time.Sleep(2 * time.Second)
	_ = s.projectRepo.UpdateStatus(ctx, project.ID, model.ProjectStatusProcessing, 80)

	time.Sleep(1 * time.Second)

	videoURL := "https://storage.example.com/videos/" + project.ID + ".mp4"
	thumbnailURL := "https://storage.example.com/thumbnails/" + project.ID + ".jpg"

	if err := s.projectRepo.SetCompleted(ctx, project.ID, videoURL, thumbnailURL); err != nil {
		_ = s.projectRepo.SetFailed(ctx, project.ID, err.Error())
		_ = s.authService.RefundCredit(ctx, project.UserID)
	}
}
