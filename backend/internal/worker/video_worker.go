package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/genvid/backend/internal/config"
	"github.com/genvid/backend/internal/repository"
	"github.com/genvid/backend/internal/zhipu"
)

type VideoGenerationJob struct {
	ProjectID   string `json:"project_id"`
	UserID      string `json:"user_id"`
	AvatarID    string `json:"avatar_id"`
	Script      string `json:"script"`
	Language    string `json:"language"`
	Format      string `json:"format"`
	ProductName string `json:"product_name"`
}

type VideoWorker struct {
	zhipuClient *zhipu.Client
	projectRepo *repository.ProjectRepository
	profileRepo *repository.ProfileRepository
	cfg         *config.Config
}

func NewVideoWorker(
	zhipuClient *zhipu.Client,
	projectRepo *repository.ProjectRepository,
	profileRepo *repository.ProfileRepository,
	cfg *config.Config,
) *VideoWorker {
	return &VideoWorker{
		zhipuClient: zhipuClient,
		projectRepo: projectRepo,
		profileRepo: profileRepo,
		cfg:         cfg,
	}
}

func (w *VideoWorker) ProcessJob(ctx context.Context, jobData []byte) error {
	var job VideoGenerationJob
	if err := json.Unmarshal(jobData, &job); err != nil {
		return fmt.Errorf("failed to unmarshal job: %w", err)
	}

	log.Printf("Processing video generation job for project %s", job.ProjectID)

	if err := w.projectRepo.UpdateStatus(ctx, job.ProjectID, "processing", 10); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Build prompt from script and product info
	prompt := job.Script
	if job.ProductName != "" {
		prompt = fmt.Sprintf("Product: %s\n\n%s", job.ProductName, job.Script)
	}

	// Determine video size based on format
	size := w.getVideoSize(job.Format)

	// Generate video using Zhipu CogVideoX
	req := zhipu.VideoGenerationRequest{
		Model:   w.cfg.External.Zhipu.Model,
		Prompt:  prompt,
		Quality: "speed",
		Size:    size,
	}

	resp, err := w.zhipuClient.GenerateVideo(req)
	if err != nil {
		w.handleFailure(ctx, job, err.Error())
		return fmt.Errorf("failed to generate video: %w", err)
	}

	if err := w.projectRepo.SetProcessing(ctx, job.ProjectID, resp.ID, "zhipu"); err != nil {
		return fmt.Errorf("failed to set processing: %w", err)
	}

	_ = w.projectRepo.UpdateStatus(ctx, job.ProjectID, "processing", 30)

	// Wait for video completion (5 minute timeout)
	result, err := w.zhipuClient.WaitForCompletion(resp.ID, 5*time.Minute)
	if err != nil {
		w.handleFailure(ctx, job, err.Error())
		return fmt.Errorf("video generation failed: %w", err)
	}

	_ = w.projectRepo.UpdateStatus(ctx, job.ProjectID, "processing", 80)

	var videoURL, thumbnailURL string
	if result.VideoResult != nil {
		videoURL = result.VideoResult.URL
		thumbnailURL = result.VideoResult.CoverURL
	}

	if err := w.projectRepo.SetCompleted(ctx, job.ProjectID, videoURL, thumbnailURL); err != nil {
		return fmt.Errorf("failed to mark completed: %w", err)
	}

	log.Printf("Video generation completed for project %s", job.ProjectID)
	return nil
}

func (w *VideoWorker) handleFailure(ctx context.Context, job VideoGenerationJob, errMsg string) {
	if err := w.projectRepo.SetFailed(ctx, job.ProjectID, errMsg); err != nil {
		log.Printf("Failed to mark project as failed: %v", err)
	}

	if err := w.profileRepo.IncrementCredits(ctx, job.UserID, 1); err != nil {
		log.Printf("Failed to refund credit: %v", err)
	}
}

func (w *VideoWorker) getVideoSize(format string) string {
	sizes := map[string]string{
		"9:16": "1080x1920",
		"1:1":  "1024x1024",
		"16:9": "1920x1080",
	}

	if size, ok := sizes[format]; ok {
		return size
	}
	return "1080x1920" // Default to vertical video for TikTok/Reels
}

func (w *VideoWorker) Start(ctx context.Context, jobs <-chan []byte) {
	log.Println("Video worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Video worker stopped")
			return
		case jobData := <-jobs:
			if err := w.ProcessJob(ctx, jobData); err != nil {
				log.Printf("Job failed: %v", err)
			}
		}
	}
}
