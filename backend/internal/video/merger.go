package video

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Merger struct {
	tempDir string
}

func NewMerger(tempDir string) *Merger {
	os.MkdirAll(tempDir, 0755)
	return &Merger{tempDir: tempDir}
}

func (m *Merger) DownloadVideo(url, filename string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download video: status %d", resp.StatusCode)
	}

	filePath := filepath.Join(m.tempDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save video: %w", err)
	}

	return filePath, nil
}

func (m *Merger) MergeVideos(videoURLs []string, outputPath string) (string, error) {
	if len(videoURLs) == 0 {
		return "", fmt.Errorf("no videos to merge")
	}

	if len(videoURLs) == 1 {
		return videoURLs[0], nil
	}

	localFiles := make([]string, len(videoURLs))
	for i, url := range videoURLs {
		filename := fmt.Sprintf("segment_%d.mp4", i)
		localPath, err := m.DownloadVideo(url, filename)
		if err != nil {
			m.cleanup(localFiles)
			return "", fmt.Errorf("failed to download segment %d: %w", i, err)
		}
		localFiles[i] = localPath
	}

	listFile := filepath.Join(m.tempDir, "concat_list.txt")
	listContent := ""
	for _, f := range localFiles {
		listContent += fmt.Sprintf("file '%s'\n", f)
	}
	if err := os.WriteFile(listFile, []byte(listContent), 0644); err != nil {
		m.cleanup(localFiles)
		return "", fmt.Errorf("failed to create concat list: %w", err)
	}

	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c", "copy",
		"-y",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.cleanup(localFiles)
		os.Remove(listFile)
		return "", fmt.Errorf("ffmpeg failed: %w, output: %s", err, string(output))
	}

	m.cleanup(localFiles)
	os.Remove(listFile)

	return outputPath, nil
}

func (m *Merger) cleanup(files []string) {
	for _, f := range files {
		os.Remove(f)
	}
}

func (m *Merger) GetOutputPath(projectID string) string {
	return filepath.Join(m.tempDir, fmt.Sprintf("%s_merged.mp4", projectID))
}

func CheckFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg not installed or not in PATH")
	}
	return nil
}

func SplitScript(script string, segments int) []string {
	if segments <= 1 {
		return []string{script}
	}

	sentences := splitIntoSentences(script)
	if len(sentences) < segments {
		segments = len(sentences)
		if segments == 0 {
			return []string{script}
		}
	}

	result := make([]string, segments)
	perSegment := len(sentences) / segments
	extra := len(sentences) % segments

	idx := 0
	for i := 0; i < segments; i++ {
		count := perSegment
		if i < extra {
			count++
		}
		end := idx + count
		if end > len(sentences) {
			end = len(sentences)
		}
		result[i] = strings.Join(sentences[idx:end], " ")
		idx = end
	}

	return result
}

func splitIntoSentences(text string) []string {
	replacements := map[string]string{
		"!":   "!|",
		"?":   "?|",
		".":   ".|",
		"!\n": "!|\n",
		"?\n": "?|\n",
		".\n": ".|\n",
	}

	delimited := text
	for old, new := range replacements {
		delimited = strings.ReplaceAll(delimited, old, new)
	}

	parts := strings.Split(delimited, "|")
	var sentences []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			sentences = append(sentences, p)
		}
	}

	return sentences
}
