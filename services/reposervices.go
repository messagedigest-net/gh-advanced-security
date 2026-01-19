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
	path := fmt.Sprintf("repos/%s", name)
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

// In services/reposervices.go

func (r *RepositoryServices) Show(name string, jsonOutput bool) error {
	// 1. Fetch the Repo
	repo, err := r.Get(name)
	if err != nil {
		return err
	}

	// 2. JSON "Escape Hatch"
	if jsonOutput {
		return jsonLister(repo)
	}

	// 3. Human-Readable Table "Vibe"
	tablePrinter, err := getTablePrinter()
	if err != nil {
		return err
	}

	// Basic Info
	tablePrinter.AddField("Repository")
	tablePrinter.AddField(repo.FullName)
	tablePrinter.EndRow()

	tablePrinter.AddField("Visibility")
	tablePrinter.AddField(repo.Visibility)
	tablePrinter.EndRow()

	tablePrinter.AddField("URL")
	tablePrinter.AddField(repo.HtmlUrl)
	tablePrinter.EndRow()

	// 4. Security Configuration Section
	tablePrinter.AddField("Security Settings:")
	tablePrinter.EndRow()

	// Safe handling for nil SecurityAndAnalysis
	sas := repo.SecurityAndAnalysis

	// Secret Scanning
	tablePrinter.AddField("\tSecret Scanning")
	if sas.SecretScanning.Status != "" {
		tablePrinter.AddField(sas.SecretScanning.Status)
	} else {
		tablePrinter.AddField("disabled/not available")
	}
	tablePrinter.EndRow()

	// Push Protection
	tablePrinter.AddField("\tPush Protection")
	if sas.SecretScanningPushProtection.Status != "" {
		tablePrinter.AddField(sas.SecretScanningPushProtection.Status)
	} else {
		tablePrinter.AddField("disabled")
	}
	tablePrinter.EndRow()

	// Advanced Security (GHAS) License Status
	tablePrinter.AddField("\tAdvanced Security")
	if sas.AdvancedSecurity.Status != "" {
		tablePrinter.AddField(sas.AdvancedSecurity.Status)
	} else {
		tablePrinter.AddField("disabled")
	}
	tablePrinter.EndRow()

	return tablePrinter.Render()
}
