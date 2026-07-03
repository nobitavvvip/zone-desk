package model

type FileItem struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	Mode    string `json:"mode"`
	ModTime string `json:"modTime"`
	Hidden  bool   `json:"hidden"`
}

type FileListResult struct {
	Path   string     `json:"path"`
	Parent *string    `json:"parent"`
	Items  []FileItem `json:"items"`
}

type FileReadResult struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}
