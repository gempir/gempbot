package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/google/uuid"
)

type Asset struct {
	ID         string `json:"id"`
	IsAnimated bool   `json:"isAnimated"`
	IsVideo    bool   `json:"isVideo"`
	MimeType   string `json:"mimeType"`
	URL        string `json:"url"`
}

var animatedMimeTypes = []string{
	"image/gif",
	"image/webp",
	"video/webm",
	"video/mp4",
}

var videoMimeTypes = []string{
	"video/webm",
	"video/mp4",
	"video/avi",
}

func (a *Api) AssetCreationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		a.AssetHandler(w, r)
		return
	}

	_, _, apiErr := a.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	assetID := uuid.New().String()
	fileType := http.DetectContentType(fileContent)

	isAnimated := false
	for _, mimeType := range animatedMimeTypes {
		if mimeType == fileType {
			isAnimated = true
			break
		}
	}
	isVideo := false
	for _, mimeType := range videoMimeTypes {
		if mimeType == fileType {
			isVideo = true
			break
		}
	}

	a.db.CreateAsset(assetID, fileType, fileContent)

	asset := Asset{
		ID:         assetID,
		IsAnimated: isAnimated,
		IsVideo:    isVideo,
		MimeType:   fileType,
		URL:        a.cfg.ApiBaseUrl + "/api/asset?id=" + assetID,
	}

	api.WriteJson(w, asset, http.StatusCreated)
}

func (a *Api) AssetHandler(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("id")
	if assetID == "" {
		http.Error(w, "No asset id provided", http.StatusBadRequest)
		return
	}

	asset := a.db.GetAsset(assetID)
	if asset == nil {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", asset.MimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(asset.Blob)))
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(asset.Blob)
}
