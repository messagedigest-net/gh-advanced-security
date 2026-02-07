package model

// RepoUpdateRequest maps to the PATCH /repos/{owner}/{repo} body
type RepoUpdateRequest struct {
	SecurityAndAnalysis *SecurityAndAnalysisReq `json:"security_and_analysis,omitempty"`
}

type SecurityAndAnalysisReq struct {
	AdvancedSecurity                  *StatusReq `json:"advanced_security,omitempty"`
	SecretScanning                    *StatusReq `json:"secret_scanning,omitempty"`
	SecretScanningNonProviderPatterns *StatusReq `json:"secret_scanning_non_provider_patterns,omitempty"`
	PushProtection                    *StatusReq `json:"secret_scanning_push_protection,omitempty"`
	DependabotSecurityUpdates         *StatusReq `json:"dependabot_security_updates,omitempty"`
}

// OrgUpdateRequest mapeia para o corpo do PATCH /orgs/{org}
// Usamos *bool para diferenciar "false" (desativar) de "nil" (ignorar)
type OrgUpdateRequest struct {
	AdvancedSecurityEnabledForNewRepos             *bool `json:"advanced_security_enabled_for_new_repositories,omitempty"`
	SecretScanningEnabledForNewRepos               *bool `json:"secret_scanning_enabled_for_new_repositories,omitempty"`
	SecretScanningPushProtectionEnabledForNewRepos *bool `json:"secret_scanning_push_protection_enabled_for_new_repositories,omitempty"`
	DependabotAlertsEnabledForNewRepos             *bool `json:"dependabot_alerts_enabled_for_new_repositories,omitempty"`
	DependabotSecurityUpdatesEnabledForNewRepos    *bool `json:"dependabot_security_updates_enabled_for_new_repositories,omitempty"`
	DependencyGraphEnabledForNewRepos              *bool `json:"dependency_graph_enabled_for_new_repositories,omitempty"`
}

type StatusReq struct {
	Status string `json:"status"`
}
