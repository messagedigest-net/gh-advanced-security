package model

type SecretScanningAlert struct {
	Number                   int    `json:"number"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
	Url                      string `json:"url"`
	HtmlUrl                  string `json:"html_url"`
	State                    string `json:"state"`
	SecretType               string `json:"secret_type"`
	SecretTypeDisplayName    string `json:"secret_type_display_name"`
	Secret                   string `json:"secret"`
	Resolution               string `json:"resolution"`
	ResolvedBy               User   `json:"resolved_by"`
	ResolvedAt               string `json:"resolved_at"`
	PushProtectionBypassed   bool   `json:"push_protection_bypassed"`
	PushProtectionBypassedBy User   `json:"push_protection_bypassed_by"`
	PushProtectionBypassedAt string `json:"push_protection_bypassed_at"`
}
