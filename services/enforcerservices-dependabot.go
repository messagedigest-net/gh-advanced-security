package services

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/messagedigest-net/gh-advanced-security/model"
)

// EnableDependabotAlerts atua no endpoint /vulnerability-alerts
func (e *EnforcerServices) EnableDependabotAlerts(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s/vulnerability-alerts", owner, repo)

	// O endpoint PUT /vulnerability-alerts não requer corpo para ativar
	// Precisamos de um helper 'put' no core.go ou usar client.Request diretamente.
	// Assumindo que você tem acesso ao client ou criará um helper 'put':

	resp, err := client.Request("PUT", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 { // 204 No Content é o sucesso padrão aqui
		return nil
	}
	return api.HandleHTTPError(resp)
}

func (e *EnforcerServices) DisableDependabotAlerts(owner, repo string) error {
	// A API usa DELETE para desativar alertas
	path := fmt.Sprintf("repos/%s/%s/vulnerability-alerts", owner, repo)
	resp, err := client.Request("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return api.HandleHTTPError(resp)
}

// EnableDependabotSecurityUpdates usa o PATCH no security_and_analysis
func (e *EnforcerServices) EnableDependabotSecurityUpdates(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)

	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			DependabotSecurityUpdates: &model.StatusReq{Status: "enabled"},
		},
	}

	return patch(path, payload)
}

func (e *EnforcerServices) DisableDependabotSecurityUpdates(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)
	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			DependabotSecurityUpdates: &model.StatusReq{Status: "disabled"},
		},
	}
	return patch(path, payload)
}

// BulkEnableDependabot agora usa a API de Organização (Muito mais rápido!)
func (e *EnforcerServices) BulkEnableDependabot(org string) error {
	fmt.Printf("Enabling Dependabot features for organization '%s'...\n", org)

	// 1. Existing Repos
	if err := e.SetOrgSecurityFeature(org, "dependency_graph", "enable_all"); err != nil {
		return err
	}
	if err := e.SetOrgSecurityFeature(org, "dependabot_alerts", "enable_all"); err != nil {
		return err
	}
	if err := e.SetOrgSecurityFeature(org, "dependabot_security_updates", "enable_all"); err != nil {
		return err
	}
	fmt.Println("- Existing Repos: Enabling initiated.")

	// 2. New Repos Policy
	settings := model.OrgUpdateRequest{
		DependencyGraphEnabledForNewRepos:           boolPtr(true),
		DependabotAlertsEnabledForNewRepos:          boolPtr(true),
		DependabotSecurityUpdatesEnabledForNewRepos: boolPtr(true),
	}
	if err := e.UpdateOrgSettings(org, settings); err != nil {
		return fmt.Errorf("failed to update org policy: %w", err)
	}
	fmt.Println("- New Repos Policy: Updated.")

	return nil
}

// BulkDisableDependabot: Desativa Updates, Alerts e Graph (nesta ordem)
func (e *EnforcerServices) BulkDisableDependabot(org string) error {
	fmt.Printf("Disabling Dependabot features for organization '%s'...\n", org)

	// 1. Atualizar Política para NOVOS Repositórios
	// Desliga todas as flags de Dependabot/Graph para novos projetos
	settings := model.OrgUpdateRequest{
		DependabotSecurityUpdatesEnabledForNewRepos: boolPtr(false),
		DependabotAlertsEnabledForNewRepos:          boolPtr(false),
		DependencyGraphEnabledForNewRepos:           boolPtr(false),
	}
	if err := e.UpdateOrgSettings(org, settings); err != nil {
		return fmt.Errorf("failed to update org policy: %w", err)
	}
	fmt.Println("- New Repos Policy: Disabled.")

	// 2. Desativar em Repositórios EXISTENTES
	// Ordem Crítica: Updates -> Alerts -> Graph (para evitar erros de dependência)

	// A. Security Updates
	if err := e.SetOrgSecurityFeature(org, "dependabot_security_updates", "disable_all"); err != nil {
		return fmt.Errorf("failed to disable Security Updates: %w", err)
	}
	fmt.Println("- Existing Repos: Security Updates Disabled.")

	// B. Dependabot Alerts
	if err := e.SetOrgSecurityFeature(org, "dependabot_alerts", "disable_all"); err != nil {
		return fmt.Errorf("failed to disable Dependabot Alerts: %w", err)
	}
	fmt.Println("- Existing Repos: Dependabot Alerts Disabled.")

	// C. Dependency Graph
	if err := e.SetOrgSecurityFeature(org, "dependency_graph", "disable_all"); err != nil {
		return fmt.Errorf("failed to disable Dependency Graph: %w", err)
	}
	fmt.Println("- Existing Repos: Dependency Graph Disabled.")

	return nil
}
