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

// FetchAllForOrg silently retrieves ALL repositories for automation (Bulk Enforcers).
// It handles pagination automatically without user interaction.
func (r *RepositoryServices) FetchAllForOrg(org string) ([]model.Repository, error) {
	var allRepos []model.Repository

	// Start with a large page size for efficiency in background tasks
	path := fmt.Sprintf("orgs/%s/repos?per_page=100", org)

	for {
		var pageRepos []model.Repository
		nextUrl, err := getPages(path, &pageRepos)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, pageRepos...)

		if nextUrl == "" {
			break
		}
		path = nextUrl
	}

	return allRepos, nil
}

// ListFor is the INTERACTIVE version (Fetch -> Print -> Prompt)
// Update signature to accept fetchAll
func (r *RepositoryServices) ListFor(org string, user bool, jsonOutput bool, userPageSize int, fetchAll bool) error {
	pageSize := GetOptimalPageSize(userPageSize)

	kind := "orgs"
	if user {
		kind = "users"
	}
	path := fmt.Sprintf("%s/%s/repos?per_page=%d", kind, org, pageSize)

	r.repositories = []model.Repository{}

	for {
		var pageRepos []model.Repository
		nextUrl, err := getPages(path, &pageRepos)
		if err != nil {
			return err
		}

		// JSON or "Fetch All" should behave similarly regarding accumulation/looping
		if jsonOutput {
			r.repositories = append(r.repositories, pageRepos...)
			if nextUrl == "" {
				break
			}
			path = nextUrl
			continue
		}

		// Interactive Rendering
		r.repositories = pageRepos
		if err := r.printRepoTable(); err != nil {
			return err
		}

		if nextUrl == "" {
			break
		}

		// LOGIC FIX: Only ask if NOT fetching all
		if !fetchAll {
			if !AskForNextPage() {
				break
			}
		}

		path = nextUrl
	}

	if jsonOutput {
		return jsonLister(r.repositories)
	}
	return nil
}

// printRepoTable renders the repository list to the terminal
func (r *RepositoryServices) printRepoTable() error {
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
		"NonProviderPatterns",
		"ValidityChecks",
		"PushProtection",
		"DependabotSecurityUpdates",
	})

	for _, repo := range r.repositories {
		// Nil check for SecurityAndAnalysis to prevent panic on public repos
		asStatus := "N/A"
		ssStatus := "N/A"
		npStatus := "N/A"
		ppStatus := "N/A"
		vcStatus := "N/A"
		duStatus := "N/A"

		if repo.SecurityAndAnalysis.AdvancedSecurity.Status != "" {
			asStatus = repo.SecurityAndAnalysis.AdvancedSecurity.Status
		}
		if repo.SecurityAndAnalysis.SecretScanning.Status != "" {
			ssStatus = repo.SecurityAndAnalysis.SecretScanning.Status
		}
		if repo.SecurityAndAnalysis.SecretScanningNonProviderPatterns.Status != "" {
			npStatus = repo.SecurityAndAnalysis.SecretScanningNonProviderPatterns.Status
		}
		if repo.SecurityAndAnalysis.SecretScanningValidityChecks.Status != "" {
			vcStatus = repo.SecurityAndAnalysis.SecretScanningValidityChecks.Status
		}
		if repo.SecurityAndAnalysis.SecretScanningPushProtection.Status != "" {
			ppStatus = repo.SecurityAndAnalysis.SecretScanningPushProtection.Status
		}
		if repo.SecurityAndAnalysis.DependabotSecurityUpdates.Status != "" {
			duStatus = repo.SecurityAndAnalysis.DependabotSecurityUpdates.Status
		}

		tablePrinter.AddField(repo.FullName)
		tablePrinter.AddField(repo.Owner.Login)
		tablePrinter.AddField(enabledOrDisabled(repo.Private))
		tablePrinter.AddField(repo.HtmlUrl)
		tablePrinter.AddField(repo.Description)
		tablePrinter.AddField(repo.Homepage)
		tablePrinter.AddField(repo.Language)
		tablePrinter.AddField(strings.Join(repo.Topics, ","))
		tablePrinter.AddField(asStatus)
		tablePrinter.AddField(ssStatus)
		tablePrinter.AddField(npStatus)
		tablePrinter.AddField(vcStatus)
		tablePrinter.AddField(ppStatus)
		tablePrinter.AddField(duStatus)
		tablePrinter.EndRow()
	}

	return tablePrinter.Render()
}

// Show renders detailed info for a single repository (from your previous Repomix)
func (r *RepositoryServices) Show(name string, jsonOutput bool) error {
	repo, err := r.Get(name)
	if err != nil {
		return err
	}

	if jsonOutput {
		return jsonLister(repo)
	}

	tablePrinter, err := getTablePrinter()
	if err != nil {
		return err
	}

	tablePrinter.AddField("Repository")
	tablePrinter.AddField(repo.FullName)
	tablePrinter.EndRow()

	tablePrinter.AddField("Visibility")
	tablePrinter.AddField(repo.Visibility)
	tablePrinter.EndRow()

	tablePrinter.AddField("URL")
	tablePrinter.AddField(repo.HtmlUrl)
	tablePrinter.EndRow()

	tablePrinter.AddField("Security Settings:")
	tablePrinter.EndRow()

	// Safe handling for nil SecurityAndAnalysis
	sas := repo.SecurityAndAnalysis

	// Helper for safe status access
	safeStatus := func(status string) string {
		if status == "" {
			return "disabled/not available"
		}
		return status
	}

	tablePrinter.AddField("\tAdvanced Security")
	tablePrinter.AddField(safeStatus(sas.AdvancedSecurity.Status))
	tablePrinter.EndRow()

	tablePrinter.AddField("\tSecret Scanning")
	tablePrinter.AddField(safeStatus(sas.SecretScanning.Status))
	tablePrinter.EndRow()

	tablePrinter.AddField("\tNon Provider Paterns")
	tablePrinter.AddField(safeStatus(sas.SecretScanningNonProviderPatterns.Status))
	tablePrinter.EndRow()

	tablePrinter.AddField("\tValidity Checks")
	tablePrinter.AddField(safeStatus(sas.SecretScanningValidityChecks.Status))
	tablePrinter.EndRow()

	tablePrinter.AddField("\tPush Protection")
	tablePrinter.AddField(safeStatus(sas.SecretScanningPushProtection.Status))
	tablePrinter.EndRow()

	tablePrinter.AddField("\tDependabot Security Updates")
	tablePrinter.AddField(safeStatus(sas.DependabotSecurityUpdates.Status))
	tablePrinter.EndRow()

	return tablePrinter.Render()
}
