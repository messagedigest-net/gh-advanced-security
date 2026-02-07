package services

import (
	"fmt"

	"github.com/messagedigest-net/gh-advanced-security/model"
)

// EnablePushProtection checks and enables feature for a single repo
func (e *EnforcerServices) EnablePushProtection(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)

	// Construct payload: Secret Scanning MUST be enabled to enable Push Protection
	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			SecretScanning: &model.StatusReq{Status: "enabled"},
			PushProtection: &model.StatusReq{Status: "enabled"},
		},
	}

	// We use client.Request directly for PATCH as Go-GH might not have a helper for this specific body
	// Note: You'll need to update services/core.go to expose a Patch helper or access client
	// For now, assuming you add a Patch method to core.go similar to Get

	// Simulated Patch call (since we need to add Patch to core.go):
	// return client.Patch(path, payload, nil)
	// See "Missing Piece" below for the core.go update
	return patch(path, payload)
}

func (e *EnforcerServices) DisablePushProtection(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)
	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			PushProtection: &model.StatusReq{Status: "disabled"},
		},
	}
	return patch(path, payload)
}

// BulkEnablePushProtection agora usa a API de Organização (O(1))
func (e *EnforcerServices) BulkEnablePushProtection(org string) error {
	fmt.Printf("Enabling Push Protection for organization '%s'...\n", org)

	// 1. Existing Repos
	fmt.Println("- Ensuring Secret Scanning is enabled...")
	if err := e.SetOrgSecurityFeature(org, "secret_scanning", "enable_all"); err != nil {
		return err
	}
	if err := e.SetOrgSecurityFeature(org, "secret_scanning_push_protection", "enable_all"); err != nil {
		return err
	}

	// 2. New Repos Policy
	// Nota: Push Protection requer Secret Scanning. Ativamos ambos na política.
	settings := model.OrgUpdateRequest{
		SecretScanningEnabledForNewRepos:               boolPtr(true),
		SecretScanningPushProtectionEnabledForNewRepos: boolPtr(true),
	}
	if err := e.UpdateOrgSettings(org, settings); err != nil {
		return err
	}
	fmt.Println("- New Repos Policy: Updated (SS & Push Protection Enabled).")

	return nil
}

// BulkDisablePushProtection: Atualiza política (Futuro) + Desativa em massa (Presente)
func (e *EnforcerServices) BulkDisablePushProtection(org string) error {
	fmt.Printf("Disabling Push Protection for organization '%s'...\n", org)

	// 1. Atualizar Política para NOVOS Repositórios
	// Define que novos repos nascerão com Push Protection DESLIGADO
	settings := model.OrgUpdateRequest{
		SecretScanningPushProtectionEnabledForNewRepos: boolPtr(false),
	}
	if err := e.UpdateOrgSettings(org, settings); err != nil {
		return fmt.Errorf("failed to update org policy: %w", err)
	}
	fmt.Println("- New Repos Policy: Disabled.")

	// 2. Desativar em Repositórios EXISTENTES
	// Usa a API Bulk para desligar a feature nos repos atuais
	if err := e.SetOrgSecurityFeature(org, "secret_scanning_push_protection", "disable_all"); err != nil {
		return fmt.Errorf("failed to disable for existing repos: %w", err)
	}
	fmt.Println("- Existing Repos: Disabling initiated.")

	return nil
}
