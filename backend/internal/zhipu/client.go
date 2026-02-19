package zhipu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client for ZhipuAI CogVideoX API
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// VideoGenerationRequest represents the request to generate a video
type VideoGenerationRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	Quality   string `json:"quality,omitempty"`    // "speed" or "quality"
	WithAudio bool   `json:"with_audio,omitempty"` // Generate AI sound effects
	Size      string `json:"size,omitempty"`       // Resolution: 720x480, 1080x1920, etc.
	FPS       int    `json:"fps,omitempty"`        // 30 or 60
	RequestID string `json:"request_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
}

// VideoGenerationResponse represents the response from video generation
type VideoGenerationResponse struct {
	RequestID  string `json:"request_id"`
	ID         string `json:"id"` // Task order number for querying
	Model      string `json:"model"`
	TaskStatus string `json:"task_status"` // PROCESSING, SUCCESS, FAILED
}

// VideoResultResponse represents the result of video generation
type VideoResultResponse struct {
	RequestID   string       `json:"request_id"`
	ID          string       `json:"id"`
	Model       string       `json:"model"`
	TaskStatus  string       `json:"task_status"` // PROCESSING, SUCCESS, FAILED
	VideoResult *VideoResult `json:"video_result,omitempty"`
}

// AsyncResultResponse represents the response from async-result endpoint
type AsyncResultResponse struct {
	Model       string        `json:"model"`
	TaskStatus  string        `json:"task_status"` // PROCESSING, SUCCESS, FAIL
	VideoResult []VideoResult `json:"video_result,omitempty"`
}

// VideoResult contains the generated video details
type VideoResult struct {
	URL      string  `json:"url"`
	Duration float64 `json:"duration"`
	CoverURL string  `json:"cover_url,omitempty"`
}

// NewClient creates a new ZhipuAI client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://open.bigmodel.cn/api/paas/v4",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// doRequest performs an HTTP request to the ZhipuAI API
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GenerateVideo generates a video from text prompt
func (c *Client) GenerateVideo(req VideoGenerationRequest) (*VideoGenerationResponse, error) {
	// Set default model if not specified
	if req.Model == "" {
		req.Model = "cogvideox-3"
	}

	// Set default quality if not specified
	if req.Quality == "" {
		req.Quality = "speed"
	}

	respBody, err := c.doRequest("POST", "/videos/generations", req)
	if err != nil {
		return nil, err
	}

	var response struct {
		RequestID  string `json:"request_id"`
		ID         string `json:"id"`
		Model      string `json:"model"`
		TaskStatus string `json:"task_status"`
		Code       string `json:"code"`
		Message    string `json:"message"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Code != "" && response.Code != "200" && response.Code != "100" {
		return nil, fmt.Errorf("API returned error: %s - %s", response.Code, response.Message)
	}

	return &VideoGenerationResponse{
		RequestID:  response.RequestID,
		ID:         response.ID,
		Model:      response.Model,
		TaskStatus: response.TaskStatus,
	}, nil
}

// GetVideoResult queries the status and result of a video generation task
func (c *Client) GetVideoResult(taskID string) (*VideoResultResponse, error) {
	endpoint := fmt.Sprintf("/async-result/%s", taskID)

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response AsyncResultResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := &VideoResultResponse{
		ID:         taskID,
		Model:      response.Model,
		TaskStatus: response.TaskStatus,
	}

	if len(response.VideoResult) > 0 {
		result.VideoResult = &response.VideoResult[0]
	}

	return result, nil
}

// WaitForCompletion polls until video generation is complete or timeout
func (c *Client) WaitForCompletion(taskID string, timeout time.Duration) (*VideoResultResponse, error) {
	start := time.Now()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := c.GetVideoResult(taskID)
			if err != nil {
				return nil, err
			}

			switch result.TaskStatus {
			case "SUCCESS":
				return result, nil
			case "FAILED", "FAIL":
				return nil, fmt.Errorf("video generation failed")
			}

			if time.Since(start) > timeout {
				return nil, fmt.Errorf("timeout waiting for video completion")
			}
		}
	}
}

// GenerateTextToVideo generates video from text prompt
func (c *Client) GenerateTextToVideo(prompt string, opts ...VideoOption) (*VideoGenerationResponse, error) {
	req := VideoGenerationRequest{
		Model:  "cogvideox-3",
		Prompt: prompt,
	}

	for _, opt := range opts {
		opt(&req)
	}

	return c.GenerateVideo(req)
}

// GenerateImageToVideo generates video from an image
func (c *Client) GenerateImageToVideo(imageURL, prompt string, opts ...VideoOption) (*VideoGenerationResponse, error) {
	req := VideoGenerationRequest{
		Model:    "cogvideox-3",
		ImageURL: imageURL,
		Prompt:   prompt,
	}

	for _, opt := range opts {
		opt(&req)
	}

	return c.GenerateVideo(req)
}

// VideoOption is a function that modifies VideoGenerationRequest
type VideoOption func(*VideoGenerationRequest)

// WithQuality sets the quality mode
func WithQuality(quality string) VideoOption {
	return func(r *VideoGenerationRequest) {
		r.Quality = quality
	}
}

// WithSize sets the video resolution
func WithSize(size string) VideoOption {
	return func(r *VideoGenerationRequest) {
		r.Size = size
	}
}

// WithAudio enables or disables AI sound effects
func WithAudio(enabled bool) VideoOption {
	return func(r *VideoGenerationRequest) {
		r.WithAudio = enabled
	}
}

// WithModel sets the model to use
func WithModel(model string) VideoOption {
	return func(r *VideoGenerationRequest) {
		r.Model = model
	}
}
