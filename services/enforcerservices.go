package services

import (
	"fmt"
	"io"

	"github.com/messagedigest-net/gh-advanced-security/model"
)

var enforcerSvcs *EnforcerServices

type EnforcerServices struct{}

func GetEnforcerServices() *EnforcerServices {
	if enforcerSvcs == nil {
		enforcerSvcs = &EnforcerServices{}
	}
	return enforcerSvcs
}

// SetOrgSecurityFeature ativa/desativa funcionalidades em massa para a Organização
// Docs: POST /orgs/{org}/{security_product}/{enablement}
// security_product: dependency_graph, dependabot_alerts, dependabot_security_updates, secret_scanning, etc.
// enablement: enable_all, disable_all
func (e *EnforcerServices) SetOrgSecurityFeature(org, product, state string) error {
	// Ex: orgs/my-org/dependabot_alerts/enable_all
	path := fmt.Sprintf("orgs/%s/%s/%s", org, product, state)

	// POST request sem corpo (body nil)
	resp, err := client.Request("POST", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// Tratamento específico para 422 (Unprocessable Entity) que pode ocorrer se não houver permissão/plano
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("API Error %d: %s | Body: %s", resp.StatusCode, resp.Status, string(body))
}

// UpdateOrgSettings atualiza as políticas padrão para NOVOS repositórios
func (e *EnforcerServices) UpdateOrgSettings(org string, settings model.OrgUpdateRequest) error {
	path := fmt.Sprintf("orgs/%s", org)
	// Reutiliza a função 'patch' que já temos em core.go
	return patch(path, settings)
}

// Helper simples para criar ponteiros bool
func boolPtr(b bool) *bool {
	return &b
}
