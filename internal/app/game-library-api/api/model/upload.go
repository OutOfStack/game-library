package model

// UploadImagesResponse represents the response for image upload
type UploadImagesResponse struct {
	Files []UploadedFileInfo `json:"files"`
}

// UploadedFileInfo represents information about an uploaded file
type UploadedFileInfo struct {
	FileName string `json:"file_name"`
	FileID   string `json:"file_id"`
	FileURL  string `json:"file_url"`
	Type     string `json:"type"`
}
