package services

import (
	"fmt"

	"github.com/messagedigest-net/gh-advanced-security/model"
)

// EnableSecretScanning enables ONLY Secret Scanning (without Push Protection)
func (e *EnforcerServices) EnableSecretScanningNonProviderPatterns(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)

	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			SecretScanningNonProviderPatterns: &model.StatusReq{Status: "enabled"},
		},
	}

	return patch(path, payload)
}

func (e *EnforcerServices) DisableSecretScanningNonProviderPatterns(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)
	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			SecretScanningNonProviderPatterns: &model.StatusReq{Status: "disabled"},
		},
	}
	return patch(path, payload)
}
