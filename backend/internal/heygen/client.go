package heygen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type VideoGenerationRequest struct {
	AvatarID   string `json:"avatar_id"`
	VoiceID    string `json:"voice_id"`
	Text       string `json:"text"`
	Background string `json:"background,omitempty"`
	Ratio      string `json:"ratio,omitempty"`
}

type VideoGenerationResponse struct {
	VideoID string `json:"video_id"`
	Status  string `json:"status"`
}

type VideoStatusResponse struct {
	VideoID  string `json:"video_id"`
	Status   string `json:"status"`
	URL      string `json:"url,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	Error    string `json:"error,omitempty"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.heygen.com/v2",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

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

	req.Header.Set("X-Api-Key", c.apiKey)
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

func (c *Client) GenerateAvatarVideo(req VideoGenerationRequest) (*VideoGenerationResponse, error) {
	payload := map[string]interface{}{
		"video_inputs": []map[string]interface{}{
			{
				"character": map[string]interface{}{
					"avatar_id": req.AvatarID,
				},
				"voice": map[string]interface{}{
					"type":     "text",
					"voice_id": req.VoiceID,
					"input":    req.Text,
				},
			},
		},
		"dimension": map[string]interface{}{
			"width":  1080,
			"height": 1920,
		},
	}

	respBody, err := c.doRequest("POST", "/video/generate", payload)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code int    `json:"code"`
		Data struct {
			VideoID string `json:"video_id"`
		} `json:"data"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Code != 100 && response.Code != 200 {
		return nil, fmt.Errorf("API returned error code %d: %s", response.Code, response.Message)
	}

	return &VideoGenerationResponse{
		VideoID: response.Data.VideoID,
		Status:  "pending",
	}, nil
}

func (c *Client) GetVideoStatus(videoID string) (*VideoStatusResponse, error) {
	endpoint := fmt.Sprintf("/video/status?video_id=%s", videoID)
	
	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code int    `json:"code"`
		Data struct {
			VideoID  string  `json:"video_id"`
			Status   string  `json:"status"`
			URL      string  `json:"video_url,omitempty"`
			Duration float64 `json:"duration,omitempty"`
			Error    string  `json:"error,omitempty"`
		} `json:"data"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &VideoStatusResponse{
		VideoID:  response.Data.VideoID,
		Status:   response.Data.Status,
		URL:      response.Data.URL,
		Duration: response.Data.Duration,
		Error:    response.Data.Error,
	}, nil
}

func (c *Client) WaitForCompletion(videoID string, timeout time.Duration) (*VideoStatusResponse, error) {
	start := time.Now()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			status, err := c.GetVideoStatus(videoID)
			if err != nil {
				return nil, err
			}

			switch status.Status {
			case "completed", "success":
				return status, nil
			case "failed", "error":
				return nil, fmt.Errorf("video generation failed: %s", status.Error)
			}

			if time.Since(start) > timeout {
				return nil, fmt.Errorf("timeout waiting for video completion")
			}
		}
	}
}

func (c *Client) ListAvatars() ([]map[string]interface{}, error) {
	respBody, err := c.doRequest("GET", "/avatars", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code int                      `json:"code"`
		Data []map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Data, nil
}
