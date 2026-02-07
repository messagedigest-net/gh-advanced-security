package services

import (
	"fmt"
	"sync"

	"github.com/messagedigest-net/gh-advanced-security/model"
)

var orgSvcs *OrganizationServices

type OrganizationServices struct {
	organizations []model.Organization
}

func GetOrganizationServices() *OrganizationServices {
	if orgSvcs == nil {
		orgSvcs = &OrganizationServices{}
	}

	return orgSvcs
}

func (o *OrganizationServices) Get(name string) (*model.Organization, error) {
	org := &model.Organization{}
	path := fmt.Sprintf("orgs/%s", name)
	err := client.Get(path, org)
	return org, err
}

func (o *OrganizationServices) GetAll() error {
	err := client.Get("user/orgs", &o.organizations)
	return err
}

// Update signature
func (o *OrganizationServices) List(jsonOutput bool, userPageSize int, fetchAll bool) error {
	pageSize := GetOptimalPageSize(userPageSize)
	path := fmt.Sprintf("user/orgs?per_page=%d", pageSize)

	o.organizations = []model.Organization{}

	for {
		var pageOrgs []model.Organization
		nextUrl, err := getPages(path, &pageOrgs)
		if err != nil {
			return err
		}

		enrichedOrgs := o.enrichOrgsInParallel(pageOrgs)

		if jsonOutput {
			o.organizations = append(o.organizations, enrichedOrgs...)
			if nextUrl == "" {
				break
			}
			path = nextUrl
			continue
		}

		o.organizations = enrichedOrgs
		if err := o.printOrgTable(); err != nil {
			return err
		}

		if nextUrl == "" {
			break
		}

		// LOGIC FIX:
		if !fetchAll {
			if !AskForNextPage() {
				break
			}
		}

		path = nextUrl
	}

	if jsonOutput {
		return jsonLister(o.organizations)
	}

	fmt.Println("* setting for new repositories")
	return nil
}

// enrichOrgsInParallel pega uma lista de orgs "magras" e busca os detalhes completos de cada uma
func (o *OrganizationServices) enrichOrgsInParallel(simpleOrgs []model.Organization) []model.Organization {
	var wg sync.WaitGroup
	var mu sync.Mutex

	results := make([]model.Organization, len(simpleOrgs))

	// Limitador de concorrência para não estourar rate limit
	semaphore := make(chan struct{}, 5)

	for i, simpleOrg := range simpleOrgs {
		wg.Add(1)
		go func(index int, name string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Adquire token
			defer func() { <-semaphore }() // Libera token

			// Busca detalhes completos: GET /orgs/{name}
			fullOrg, err := o.Get(name)

			mu.Lock()
			if err == nil {
				results[index] = *fullOrg
			} else {
				// Se falhar (ex: falta de permissão), mantém o objeto original (básico)
				// ou loga o erro se tiver logger
				results[index] = simpleOrgs[index]
			}
			mu.Unlock()
		}(i, simpleOrg.Login)
	}

	wg.Wait()
	return results
}

// Helper method to keep List clean (move your existing table logic here)
func (o *OrganizationServices) printOrgTable() error {
	tablePrinter, err := getTablePrinter()
	if err != nil {
		return err
	}

	tablePrinter.AddHeader([]string{
		"Org",
		"URL",
		"Adv Security*",
		"Secret Scanning*",
		"Push Protection*",
		"Dep. Graph*",
		"Dep. Alerts*",
		"Dep. Updates*",
	})

	for _, i := range o.organizations {
		tablePrinter.AddField(i.Login)
		tablePrinter.AddField(i.URL)
		tablePrinter.AddField(enabledOrDisabled(i.AdvancedSecurityEnabledForNewRepositories))
		tablePrinter.AddField(enabledOrDisabled(i.SecretScanningEnabledForNewRepositories))
		tablePrinter.AddField(enabledOrDisabled(i.SecretScanningPushProtectionEnabledForNewRepositories))
		tablePrinter.AddField(enabledOrDisabled(i.DependencyGraphEnabledForNewRepositories))
		tablePrinter.AddField(enabledOrDisabled(i.DependabotAlertsEnabledForNewRepositories))
		tablePrinter.AddField(enabledOrDisabled(i.DependabotSecurityUpdatesEnabledForNewRepositories))
		tablePrinter.EndRow()
	}
	return tablePrinter.Render()
}

func (o *OrganizationServices) Show(name string, json bool) error {
	org, err := o.Get(name)
	if err != nil {
		return err
	}
	if json {
		return jsonLister(org)
	}

	tablePrinter, err := getTablePrinter()
	if err != nil {
		return err
	}

	tablePrinter.AddField("Organization")
	tablePrinter.AddField(org.Login)
	tablePrinter.EndRow()
	tablePrinter.AddField("API URL")
	tablePrinter.AddField(org.URL)
	tablePrinter.EndRow()
	tablePrinter.AddField("Settings for new repositories:")
	tablePrinter.EndRow()
	tablePrinter.AddField("\tDependency Graph")
	tablePrinter.AddField(enabledOrDisabled(org.DependencyGraphEnabledForNewRepositories))
	tablePrinter.EndRow()
	tablePrinter.AddField("\tDependabot Alerts")
	tablePrinter.AddField(enabledOrDisabled(org.DependabotAlertsEnabledForNewRepositories))
	tablePrinter.EndRow()
	tablePrinter.AddField("\tDependabot Security Updates")
	tablePrinter.AddField(enabledOrDisabled(org.DependabotSecurityUpdatesEnabledForNewRepositories))
	tablePrinter.EndRow()
	tablePrinter.AddField("\tEnable Advanced Security")
	tablePrinter.AddField(enabledOrDisabled(org.AdvancedSecurityEnabledForNewRepositories))
	tablePrinter.EndRow()
	tablePrinter.AddField("\tSecret Scanning")
	tablePrinter.AddField(enabledOrDisabled(org.SecretScanningEnabledForNewRepositories))
	tablePrinter.EndRow()
	tablePrinter.AddField("\tSecret Scanning Push Protection")
	tablePrinter.AddField(enabledOrDisabled(org.SecretScanningPushProtectionEnabledForNewRepositories))
	tablePrinter.EndRow()
	tablePrinter.AddField("\tSecret Scanning Push Protection Custom Link")
	tablePrinter.AddField(org.SecretScanningPushProtectionCustomLink)
	tablePrinter.EndRow()
	tablePrinter.AddField("\tSecret Scanning Push Protection Custom Link Enabled")
	tablePrinter.AddField(enabledOrDisabled(org.SecretScanningPushProtectionCustomLinkEnabled))
	tablePrinter.EndRow()

	return tablePrinter.Render()
}
