package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"s3-browser/internal/s3"
	"s3-browser/internal/types"
)

type Handlers struct {
	s3Client *s3.Client
}

func NewHandlers(s3Client *s3.Client) *Handlers {
	return &Handlers{
		s3Client: s3Client,
	}
}

func (h *Handlers) ListDirectory(w http.ResponseWriter, r *http.Request) {
	directoryPath := r.URL.Query().Get("path")
	if directoryPath == "" {
		directoryPath = "/"
	}

	listing, err := h.s3Client.ListDirectory(r.Context(), directoryPath)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, listing)
}

func (h *Handlers) DirectDownload(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		h.sendError(w, "File path is required", http.StatusBadRequest)
		return
	}

	url, err := h.s3Client.GetFileURL(r.Context(), filePath)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на presigned URL
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		h.sendError(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Для скачивания через браузер также используем presigned URL
	url, err := h.s3Client.GetFileURL(r.Context(), filePath)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки для скачивания
	fileName := path.Base(filePath)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	// Перенаправляем на presigned URL
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handlers) StreamDownload(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		h.sendError(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле для установки правильного Content-Length
	ctx := r.Context()

	// Сначала получим информацию о файле
	fileInfo, err := h.s3Client.GetFileInfo(ctx, filePath)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := path.Base(filePath)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size, 10))

	// Стримим файл напрямую из S3
	err = h.s3Client.StreamFile(ctx, filePath, w)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) ServeUI(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

func (h *Handlers) sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handlers) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(types.ErrorResponse{Error: message})
}
