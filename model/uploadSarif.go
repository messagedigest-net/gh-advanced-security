package model

type UploadSarif struct {
	CommitSHA   string `json:"commit_sha"`
	Ref         string
	Sarif       string
	CheckoutUri string `json:"checkout_uri"`
	StartedAt   string `json:"started_at"`
	ToolName    string `json:"tool_name"`
	Validate    bool
}
