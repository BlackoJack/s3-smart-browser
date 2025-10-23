package types

type FileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	IsDirectory bool   `json:"is_directory"`
}

type DirectoryListing struct {
	Path  string     `json:"path"`
	Files []FileInfo `json:"files"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
