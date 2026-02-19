package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/genvid/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrDuplicate     = errors.New("record already exists")
	ErrUnauthorized  = errors.New("unauthorized")
)

type ProfileRepository struct {
	db *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(ctx context.Context, profile *model.Profile) error {
	query := `
		INSERT INTO profiles (id, email, full_name, credits_remaining, subscription_tier, subscription_status)
		VALUES ($1, $2, $3, 3, 'free', 'inactive')
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		profile.ID,
		profile.Email,
		profile.FullName,
	).Scan(&profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		return err
	}

	profile.CreditsRemaining = 3
	profile.SubscriptionTier = "free"
	profile.SubscriptionStatus = "inactive"

	return nil
}

func (r *ProfileRepository) GetByID(ctx context.Context, id string) (*model.Profile, error) {
	profile := &model.Profile{}
	query := `
		SELECT id, email, full_name, avatar_url, company_name,
		       credits_remaining, credits_used_total, subscription_tier, subscription_status,
		       preferred_language, email_notifications, created_at, updated_at, last_login_at
		FROM profiles
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, profile, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) GetByEmail(ctx context.Context, email string) (*model.Profile, error) {
	profile := &model.Profile{}
	query := `
		SELECT id, email, full_name, avatar_url, company_name,
		       credits_remaining, credits_used_total, subscription_tier, subscription_status,
		       preferred_language, email_notifications, created_at, updated_at, last_login_at
		FROM profiles
		WHERE email = $1
	`

	err := r.db.GetContext(ctx, profile, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) Update(ctx context.Context, profile *model.Profile) error {
	query := `
		UPDATE profiles
		SET full_name = $2, avatar_url = $3, company_name = $4,
		    preferred_language = $5, email_notifications = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		profile.ID,
		profile.FullName,
		profile.AvatarURL,
		profile.CompanyName,
		profile.PreferredLanguage,
		profile.EmailNotifications,
	).Scan(&profile.UpdatedAt)

	return err
}

func (r *ProfileRepository) UpdateLastLogin(ctx context.Context, id string) error {
	query := `UPDATE profiles SET last_login_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ProfileRepository) DecrementCredits(ctx context.Context, id string) error {
	query := `
		UPDATE profiles
		SET credits_remaining = credits_remaining - 1, credits_used_total = credits_used_total + 1
		WHERE id = $1 AND credits_remaining > 0
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no credits remaining")
	}

	return nil
}

func (r *ProfileRepository) IncrementCredits(ctx context.Context, id string, amount int) error {
	query := `UPDATE profiles SET credits_remaining = credits_remaining + $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, amount)
	return err
}

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) error {
	query := `
		INSERT INTO projects (id, user_id, product_name, product_description, product_url, status)
		VALUES ($1, $2, $3, $4, $5, 'draft')
		RETURNING id, created_at, updated_at
	`

	if project.ID == "" {
		project.ID = uuid.New().String()
	}

	err := r.db.QueryRowxContext(
		ctx,
		query,
		project.ID,
		project.UserID,
		project.ProductName,
		project.ProductDescription,
		project.ProductURL,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	return err
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
	project := &model.Project{}
	query := `
		SELECT id, user_id, avatar_id, title, product_name, product_description, product_url,
		       script, language, format, status, progress_percent, error_message,
		       external_task_id, external_provider, video_url, thumbnail_url,
		       created_at, updated_at, started_at, completed_at
		FROM projects
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, project, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]model.Project, int, error) {
	var projects []model.Project
	var total int

	countQuery := `SELECT COUNT(*) FROM projects WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, user_id, avatar_id, title, product_name, product_description, product_url,
		       script, language, format, status, progress_percent, error_message,
		       external_task_id, external_provider, video_url, thumbnail_url,
		       created_at, updated_at, started_at, completed_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &projects, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *ProjectRepository) UpdateStatus(ctx context.Context, id string, status model.ProjectStatus, progress int) error {
	query := `
		UPDATE projects
		SET status = $2, progress_percent = $3, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, status, progress)
	return err
}

func (r *ProjectRepository) SetProcessing(ctx context.Context, id string, taskID, provider string) error {
	query := `
		UPDATE projects
		SET status = 'processing', external_task_id = $2, external_provider = $3,
		    started_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, taskID, provider)
	return err
}

func (r *ProjectRepository) SetCompleted(ctx context.Context, id string, videoURL, thumbnailURL string) error {
	query := `
		UPDATE projects
		SET status = 'completed', video_url = $2, thumbnail_url = $3,
		    progress_percent = 100, completed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, videoURL, thumbnailURL)
	return err
}

func (r *ProjectRepository) SetFailed(ctx context.Context, id string, errMsg string) error {
	query := `
		UPDATE projects
		SET status = 'failed', error_message = $2, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, errMsg)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id, userID string) error {
	query := `DELETE FROM projects WHERE id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
