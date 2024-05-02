package dtos

type Claim struct {
	ID        string               `json:"id"`
	UserID    string               `json:"user_id"`
	Ammount   int                  `json:"ammount"`
	Documents []UploadFileResponse `json:"documents"`
}
