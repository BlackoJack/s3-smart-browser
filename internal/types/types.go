package types

import "time"

type FileInfo struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	IsDirectory  bool      `json:"is_directory"`
	ETag         string    `json:"etag"`
}

type DirectoryListing struct {
	Path    string     `json:"path"`
	Files   []FileInfo `json:"files"`
	HasMore bool       `json:"has_more"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
