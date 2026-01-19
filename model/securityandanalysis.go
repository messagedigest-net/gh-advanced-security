package model

type SecurityAndAnalysis struct {
	AdvancedSecurity             Status `json:"advanced_security"`
	SecretScanning               Status `json:"secret_scanning"`
	SecretScanningPushProtection Status `json:"secret_scanning_push_protection"`
	SecretScanningValidityChecks Status `json:"secret_scanning_validity_checks"`
}
