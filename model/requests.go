package model

// RepoUpdateRequest maps to the PATCH /repos/{owner}/{repo} body
type RepoUpdateRequest struct {
	SecurityAndAnalysis *SecurityAndAnalysisReq `json:"security_and_analysis,omitempty"`
}

type SecurityAndAnalysisReq struct {
	AdvancedSecurity *StatusReq `json:"advanced_security,omitempty"`
	SecretScanning   *StatusReq `json:"secret_scanning,omitempty"`
	PushProtection   *StatusReq `json:"secret_scanning_push_protection,omitempty"`
}

type StatusReq struct {
	Status string `json:"status"`
}
