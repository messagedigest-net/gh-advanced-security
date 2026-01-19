package services

import (
	"fmt"
	"strings"

	"github.com/messagedigest-net/gh-advanced-security/model"
)

var repoSvcs *RepositoryServices

type RepositoryServices struct {
	repositories []model.Repository
	next         string
	org          string
}

func GetRepositoryServices() *RepositoryServices {
	if repoSvcs == nil {
		repoSvcs = &RepositoryServices{}
	}

	return repoSvcs
}

func (r *RepositoryServices) Get(name string) (*model.Repository, error) {
	repo := &model.Repository{}
	path := fmt.Sprintf("orgs/%s", name)
	err := client.Get(path, repo)
	return repo, err
}

func (r *RepositoryServices) getAll(org string, user bool) (err error) {
	var path string
	var repos []model.Repository
	if strings.Compare(r.org, org) == 0 && r.HasNext() {
		path = r.next
	} else {
		r.repositories = []model.Repository{}
		r.org = org
		kind := "orgs"
		if user {
			kind = "users"
		}
		path = fmt.Sprintf("%s/%s/repos", kind, org)
	}

	r.next, err = getPages(path, &repos)
	if err != nil {
		return err
	}
	fmt.Println(path, r.next)

	r.repositories = append(r.repositories, repos...)

	return nil
}

func (r *RepositoryServices) HasNext() bool {
	return len(r.next) > 0
}

func (r *RepositoryServices) ListFor(org string, user bool, json bool) (err error) {
	for {
		if err = r.getAll(org, user); err != nil {
			return err
		}
		if !r.HasNext() {
			break
		}
	}

	if json {
		return jsonLister(r.repositories)
	}

	tablePrinter, err := getTablePrinter()
	if err != nil {
		return err
	}

	tablePrinter.AddHeader([]string{
		"Full Name",
		"Owner",
		"Private",
		"HtmlUrl",
		"Description",
		"Homepage",
		"Language",
		"Topics",
		"AdvancedSecurity",
		"SecretScanning",
		"PushProtection",
		"ValidityChecks",
	})
	for _, repo := range r.repositories {
		tablePrinter.AddField(repo.FullName)
		tablePrinter.AddField(repo.Owner.Login)
		tablePrinter.AddField(enabledOrDisabled(repo.Private))
		tablePrinter.AddField(repo.HtmlUrl)
		tablePrinter.AddField(repo.Description)
		tablePrinter.AddField(repo.Homepage)
		tablePrinter.AddField(repo.Language)
		tablePrinter.AddField(strings.Join(repo.Topics, ","))
		tablePrinter.AddField(repo.SecurityAndAnalysis.AdvancedSecurity.Status)
		tablePrinter.AddField(repo.SecurityAndAnalysis.SecretScanning.Status)
		tablePrinter.AddField(repo.SecurityAndAnalysis.SecretScanningPushProtection.Status)
		tablePrinter.AddField(repo.SecurityAndAnalysis.SecretScanningValidityChecks.Status)
		tablePrinter.EndRow()
	}

	return tablePrinter.Render()
}

func (r *RepositoryServices) Show(name string, json bool) error {
	repo, err := r.Get(name)
	if err != nil {
		return err
	}
	if json {
		return jsonLister(repo)
	}

	tablePrinter, err := getTablePrinter()
	if err != nil {
		return err
	}

	//TODO: Ajustar saída
	/*
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
	*/
	return tablePrinter.Render()
}
