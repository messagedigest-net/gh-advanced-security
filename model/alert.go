package model

type Alert struct {
	Numer              int
	CreatedAt          string `json:"created_at"`
	URL                string
	HtmlUrl            string `json:"html_url"`
	State              string
	DismissedBy        User   `json:"dismissed_by"`
	DismissedAt        string `json:"dismissed_at"`
	DismissedReason    string `json:"dismissed_reason"`
	DismissedComment   string `json:"dismissed_comment"`
	Rule               Rule
	Tool               Tool
	MostRecentInstance Instance `json:"most_recent_instance"`
	InstancesUrl       string   `json:"instances_url"`
	Repository         Repository
}
