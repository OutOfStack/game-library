package model

// UploadImagesResponse represents the response for image upload
type UploadImagesResponse struct {
	Files []UploadedFileInfo `json:"files"`
}

// UploadedFileInfo represents information about an uploaded file
type UploadedFileInfo struct {
	FileName string `json:"fileName"`
	FileID   string `json:"fileId"`
	FileURL  string `json:"fileUrl"`
	Type     string `json:"type"` // "cover" / "screenshot"
}
