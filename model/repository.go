package model

type Repository struct {
	ID                  int
	NodeID              string `json:"node_id"`
	Name                string
	FullName            string `json:"full_name"`
	Owner               Organization
	Private             bool
	HtmlUrl             string `json:"html_url"`
	Description         string
	URL                 string
	Homepage            string
	Language            string
	Topics              []string
	Archived            bool
	Disabled            bool
	Visibility          string
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	SecurityAndAnalysis SecurityAndAnalysis
}
