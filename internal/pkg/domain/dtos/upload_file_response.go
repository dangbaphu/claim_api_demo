package dtos

type UploadFileResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
}
