package model

type FileInfo struct {
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
	IsFolder     bool   `json:"isFolder"`
	LastModified int64  `json:"lastModified"`
}
