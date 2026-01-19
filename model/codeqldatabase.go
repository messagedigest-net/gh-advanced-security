package model

type CodeQLDatabase struct {
	ID          int
	Name        string
	Language    string
	Uploader    User
	ContentType string `json:"content_type"`
	Size        int
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	URL         string
	CommitID    string `json:"commit_id"`
}
