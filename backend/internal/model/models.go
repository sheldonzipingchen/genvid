package model

import (
	"time"
)

type Profile struct {
	ID                 string     `json:"id" db:"id"`
	Email              string     `json:"email" db:"email"`
	FullName           *string    `json:"full_name,omitempty" db:"full_name"`
	AvatarURL          *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CompanyName        *string    `json:"company_name,omitempty" db:"company_name"`
	CreditsRemaining   int        `json:"credits_remaining" db:"credits_remaining"`
	CreditsUsedTotal   int        `json:"credits_used_total" db:"credits_used_total"`
	SubscriptionTier   string     `json:"subscription_tier" db:"subscription_tier"`
	SubscriptionStatus string     `json:"subscription_status" db:"subscription_status"`
	PreferredLanguage  string     `json:"preferred_language" db:"preferred_language"`
	EmailNotifications bool       `json:"email_notifications" db:"email_notifications"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type ProjectStatus string

const (
	ProjectStatusDraft      ProjectStatus = "draft"
	ProjectStatusQueued     ProjectStatus = "queued"
	ProjectStatusProcessing ProjectStatus = "processing"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusFailed     ProjectStatus = "failed"
	ProjectStatusCanceled   ProjectStatus = "canceled"
)

type VideoFormat string

const (
	VideoFormat916 VideoFormat = "9:16"
	VideoFormat11  VideoFormat = "1:1"
	VideoFormat169 VideoFormat = "16:9"
)

type Project struct {
	ID                 string        `json:"id" db:"id"`
	UserID             string        `json:"user_id" db:"user_id"`
	AvatarID           *string       `json:"avatar_id,omitempty" db:"avatar_id"`
	Title              *string       `json:"title,omitempty" db:"title"`
	ProductName        *string       `json:"product_name,omitempty" db:"product_name"`
	ProductDescription *string       `json:"product_description,omitempty" db:"product_description"`
	ProductURL         *string       `json:"product_url,omitempty" db:"product_url"`
	ProductImageURL    *string       `json:"product_image_url,omitempty" db:"product_image_url"`
	Script             *string       `json:"script,omitempty" db:"script"`
	Language           string        `json:"language" db:"language"`
	Format             VideoFormat   `json:"format" db:"format"`
	VideoDuration      int           `json:"video_duration" db:"video_duration"`
	Status             ProjectStatus `json:"status" db:"status"`
	ProgressPercent    int           `json:"progress_percent" db:"progress_percent"`
	ErrorMessage       *string       `json:"error_message,omitempty" db:"error_message"`
	ExternalTaskID     *string       `json:"external_task_id,omitempty" db:"external_task_id"`
	ExternalProvider   *string       `json:"external_provider,omitempty" db:"external_provider"`
	VideoURL           *string       `json:"video_url,omitempty" db:"video_url"`
	ThumbnailURL       *string       `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	CreatedAt          time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" db:"updated_at"`
	StartedAt          *time.Time    `json:"started_at,omitempty" db:"started_at"`
	CompletedAt        *time.Time    `json:"completed_at,omitempty" db:"completed_at"`
}

type Avatar struct {
	ID              string   `json:"id" db:"id"`
	Name            string   `json:"name" db:"name"`
	DisplayName     *string  `json:"display_name,omitempty" db:"display_name"`
	Gender          *string  `json:"gender,omitempty" db:"gender"`
	AgeRange        *string  `json:"age_range,omitempty" db:"age_range"`
	Ethnicity       *string  `json:"ethnicity,omitempty" db:"ethnicity"`
	Style           string   `json:"style" db:"style"`
	Languages       []string `json:"languages" db:"languages"`
	PreviewVideoURL *string  `json:"preview_video_url,omitempty" db:"preview_video_url"`
	ThumbnailURL    *string  `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	IsPremium       bool     `json:"is_premium" db:"is_premium"`
	UsageCount      int      `json:"usage_count" db:"usage_count"`
}

type ScriptTemplate struct {
	ID         string `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Category   string `json:"category" db:"category"`
	Template   string `json:"template" db:"template_text"`
	Language   string `json:"language" db:"language"`
	IsPremium  bool   `json:"is_premium" db:"is_premium"`
	UsageCount int    `json:"usage_count" db:"usage_count"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int64   `json:"expires_in"`
	User         Profile `json:"user"`
}

type CreateProjectRequest struct {
	ProductName        string  `json:"product_name" validate:"required"`
	ProductDescription *string `json:"product_description,omitempty"`
	ProductURL         *string `json:"product_url,omitempty"`
	ProductImageURL    *string `json:"product_image_url,omitempty"`
}

type GenerateVideoRequest struct {
	Script        string `json:"script" validate:"required,min=10,max=5000"`
	Language      string `json:"language" validate:"required,len=2"`
	Format        string `json:"format" validate:"required,oneof=9:16 1:1 16:9"`
	VideoDuration int    `json:"video_duration" validate:"oneof=5 10 30"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func SuccessResponseWithMeta(data interface{}, meta *Meta) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

func ErrorResponse(code, message string, details map[string]string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
