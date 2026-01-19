package model

type PushProtectionBypass struct {
	ID               int `json:"id"`
	Repository       Repository
	SecretType       string `json:"secret_type"`
	RulesetName      string `json:"ruleset_name"`
	CreatedAt        string `json:"created_at"`
	Reviewer         User   `json:"reviewer"`
	Status           string `json:"status"` // e.g., "approved", "denied"
	Requester        User   `json:"requester"`
	RequesterComment string `json:"requester_comment"`
}
