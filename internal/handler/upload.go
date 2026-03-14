package handler

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mfb27/luban/internal/model"
	"github.com/minio/minio-go/v7"
)

type uploadResp struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func (a *App) upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	tp := sniffType(header)
	if tp != "image" && tp != "video" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only image/video allowed"})
		return
	}

	objectKey := buildObjectKey(header.Filename)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	// PutObject needs size if possible; -1 works with streaming but uses multipart.
	_, err = a.minio.Client.PutObject(ctx, a.minioBucket, objectKey, file, -1, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	urlStr, err := a.minio.PublicURL(objectKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	att := model.Attachment{
		ID:        uuid.NewString(),
		Bucket:    a.minioBucket,
		ObjectKey: objectKey,
		URL:       urlStr,
		Type:      tp,
		CreatedAt: time.Now(),
	}
	if err := a.db.Create(&att).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, uploadResp{ID: att.ID, URL: att.URL})
}

func sniffType(header *multipart.FileHeader) string {
	ct := strings.ToLower(header.Header.Get("Content-Type"))
	switch {
	case strings.HasPrefix(ct, "image/"):
		return "image"
	case strings.HasPrefix(ct, "video/"):
		return "video"
	}

	// fallback by extension
	ext := strings.ToLower(path.Ext(header.Filename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp":
		return "image"
	case ".mp4", ".webm", ".mov", ".mkv", ".avi":
		return "video"
	default:
		return ""
	}
}

func buildObjectKey(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	if ext == "" {
		ext = ".bin"
	}
	return fmt.Sprintf("uploads/%s/%s%s", time.Now().Format("2006-01-02"), uuid.NewString(), ext)
}

// In case some clients provide no content-type and we want to inspect the first bytes.
func _unusedSniffByMagic(r io.Reader) string { // keep for future use
	buf := make([]byte, 512)
	n, _ := io.ReadFull(r, buf)
	_ = n
	return ""
}

