package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadDir string
	baseURL   string
}

func NewUploadHandler(uploadDir, baseURL string) *UploadHandler {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
	return &UploadHandler{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	maxSize := int64(10 << 20)
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	if err := r.ParseMultipartForm(maxSize); err != nil {
		respondError(w, http.StatusBadRequest, "FILE_TOO_LARGE", "File size exceeds 10MB limit", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_FILE", "No file provided", nil)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		respondError(w, http.StatusBadRequest, "INVALID_TYPE", "Only image files (jpg, png, gif, webp) are allowed", nil)
		return
	}

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		respondError(w, http.StatusBadRequest, "READ_ERROR", "Failed to read file", nil)
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		respondError(w, http.StatusInternalServerError, "SEEK_ERROR", "Failed to process file", nil)
		return
	}

	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		respondError(w, http.StatusBadRequest, "INVALID_CONTENT", "File is not a valid image", nil)
		return
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	userDir := filepath.Join(h.uploadDir, userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, "CREATE_DIR_ERROR", "Failed to create upload directory", nil)
		return
	}

	filePath := filepath.Join(userDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "SAVE_ERROR", "Failed to save file", nil)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, http.StatusInternalServerError, "SAVE_ERROR", "Failed to save file", nil)
		return
	}

	imageURL := fmt.Sprintf("%s/uploads/%s/%s", h.baseURL, userID, filename)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"url":      imageURL,
			"filename": filename,
		},
	})
}

func (h *UploadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	filePath := filepath.Join(h.uploadDir, userID, req.Filename)
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "File not found", nil)
			return
		}
		respondError(w, http.StatusInternalServerError, "DELETE_ERROR", "Failed to delete file", nil)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "File deleted",
	})
}
