package model

type SecurityAndAnalysis struct {
	AdvancedSecurity                  Status `json:"advanced_security"`
	SecretScanning                    Status `json:"secret_scanning"`
	SecretScanningNonProviderPatterns Status `json:"secret_scanning_non_provider_patterns"`
	SecretScanningValidityChecks      Status `json:"secret_scanning_validity_checks"`
	SecretScanningPushProtection      Status `json:"secret_scanning_push_protection"`
	DependabotSecurityUpdates         Status `json:"dependabot_security_updates"`
}
