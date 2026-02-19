package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/genvid/backend/internal/config"
	"github.com/genvid/backend/internal/heygen"
	"github.com/genvid/backend/internal/repository"
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
	heygenClient *heygen.Client
	projectRepo  *repository.ProjectRepository
	profileRepo  *repository.ProfileRepository
	cfg          *config.Config
}

func NewVideoWorker(
	heygenClient *heygen.Client,
	projectRepo *repository.ProjectRepository,
	profileRepo *repository.ProfileRepository,
	cfg *config.Config,
) *VideoWorker {
	return &VideoWorker{
		heygenClient: heygenClient,
		projectRepo:  projectRepo,
		profileRepo:  profileRepo,
		cfg:          cfg,
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

	req := heygen.VideoGenerationRequest{
		AvatarID: job.AvatarID,
		VoiceID:  w.getVoiceID(job.Language),
		Text:     job.Script,
		Ratio:    job.Format,
	}

	resp, err := w.heygenClient.GenerateAvatarVideo(req)
	if err != nil {
		w.handleFailure(ctx, job, err.Error())
		return fmt.Errorf("failed to generate video: %w", err)
	}

	if err := w.projectRepo.SetProcessing(ctx, job.ProjectID, resp.VideoID, "heygen"); err != nil {
		return fmt.Errorf("failed to set processing: %w", err)
	}

	_ = w.projectRepo.UpdateStatus(ctx, job.ProjectID, "processing", 30)

	status, err := w.heygenClient.WaitForCompletion(resp.VideoID, 5*time.Minute)
	if err != nil {
		w.handleFailure(ctx, job, err.Error())
		return fmt.Errorf("video generation failed: %w", err)
	}

	_ = w.projectRepo.UpdateStatus(ctx, job.ProjectID, "processing", 80)

	videoURL := status.URL
	thumbnailURL := ""

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

func (w *VideoWorker) getVoiceID(language string) string {
	voices := map[string]string{
		"en": "en-US-JennyNeural",
		"es": "es-ES-ElviraNeural",
		"fr": "fr-FR-DeniseNeural",
		"de": "de-DE-KatjaNeural",
		"it": "it-IT-ElsaNeural",
		"pt": "pt-BR-FranciscaNeural",
		"ja": "ja-JP-NanamiNeural",
		"ko": "ko-KR-SunHiNeural",
		"zh": "zh-CN-XiaoxiaoNeural",
	}

	if voiceID, ok := voices[language]; ok {
		return voiceID
	}
	return voices["en"]
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
