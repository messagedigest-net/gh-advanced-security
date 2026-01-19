package model

type Organization struct {
	Login                                                 string
	Id                                                    int
	URL                                                   string
	DependencyGraphEnabledForNewRepositories              bool   `json:"dependency_graph_enabled_for_new_repositories"`
	DependabotAlertsEnabledForNewRepositories             bool   `json:"dependabot_alerts_enabled_for_new_repositories"`
	DependabotSecurityUpdatesEnabledForNewRepositories    bool   `json:"dependabot_security_updates_enabled_for_new_repositories"`
	AdvancedSecurityEnabledForNewRepositories             bool   `json:"advanced_security_enabled_for_new_repositories"`
	SecretScanningEnabledForNewRepositories               bool   `json:"secret_scanning_enabled_for_new_repositories"`
	SecretScanningPushProtectionEnabledForNewRepositories bool   `json:"secret_scanning_push_protection_enabled_for_new_repositories"`
	SecretScanningPushProtectionCustomLink                string `json:"secret_scanning_push_protection_custom_link"`
	SecretScanningPushProtectionCustomLinkEnabled         bool   `json:"secret_scanning_push_protection_custom_link_enabled"`
}
