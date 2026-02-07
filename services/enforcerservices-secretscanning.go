package services

import (
	"fmt"

	"github.com/messagedigest-net/gh-advanced-security/model"
)

// EnableSecretScanning enables ONLY Secret Scanning (without Push Protection)
func (e *EnforcerServices) EnableSecretScanning(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)

	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			SecretScanning: &model.StatusReq{Status: "enabled"},
			// We intentionally omit PushProtection here
		},
	}

	return patch(path, payload)
}

func (e *EnforcerServices) DisableSecretScanning(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)
	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			SecretScanning: &model.StatusReq{Status: "disabled"},
		},
	}
	return patch(path, payload)
}

// BulkEnableSecretScanning agora usa a API de Organização (O(1))
func (e *EnforcerServices) BulkEnableSecretScanning(org string) error {
	fmt.Printf("Enabling Secret Scanning for organization '%s'...\n", org)

	// 1. Atualizar Repositórios EXISTENTES (Async Bulk API)
	if err := e.SetOrgSecurityFeature(org, "secret_scanning", "enable_all"); err != nil {
		return fmt.Errorf("failed to enable for existing repos: %w", err)
	}
	fmt.Println("- Existing Repos: Enabling initiated.")

	// 2. Atualizar Política para NOVOS Repositórios (Org Settings)
	settings := model.OrgUpdateRequest{
		SecretScanningEnabledForNewRepos: boolPtr(true),
	}
	if err := e.UpdateOrgSettings(org, settings); err != nil {
		return fmt.Errorf("failed to update org policy: %w", err)
	}
	fmt.Println("- New Repos Policy: Updated to Enabled.")

	return nil
}

// BulkDisableSecretScanning com ordem de dependência corrigida (Push Protection -> Secret Scanning)
func (e *EnforcerServices) BulkDisableSecretScanning(org string) error {
	fmt.Printf("Disabling Secret Scanning for organization '%s'...\n", org)

	// 1. Política para NOVOS Repositórios (Importante fazer isso ANTES ou DEPOIS?
	// Geralmente tanto faz, mas desligar a política primeiro garante que nenhum repo novo seja criado com a feature enquanto limpamos)
	settings := model.OrgUpdateRequest{
		SecretScanningEnabledForNewRepos:               boolPtr(false),
		SecretScanningPushProtectionEnabledForNewRepos: boolPtr(false),
	}
	if err := e.UpdateOrgSettings(org, settings); err != nil {
		return fmt.Errorf("failed to update org policy: %w", err)
	}
	fmt.Println("- New Repos Policy: Disabled.")

	// 2. Repositórios EXISTENTES (Lembrando da ordem de dependência)
	// Desativar Push Protection primeiro
	if err := e.SetOrgSecurityFeature(org, "secret_scanning_push_protection", "disable_all"); err != nil {
		return fmt.Errorf("failed to disable PP for existing repos: %w", err)
	}
	// Desativar Secret Scanning
	if err := e.SetOrgSecurityFeature(org, "secret_scanning", "disable_all"); err != nil {
		return fmt.Errorf("failed to disable SS for existing repos: %w", err)
	}

	return nil
}
