package services

import (
	"fmt"

	"github.com/messagedigest-net/gh-advanced-security/model"
)

var depSvcs *DependencyServices

type DependencyServices struct {
	alerts []model.DependabotAlert
}

func GetDependencyServices() *DependencyServices {
	if depSvcs == nil {
		depSvcs = &DependencyServices{}
	}
	return depSvcs
}

// ListDependabotAlerts fetches alerts using your standardized pagination
func (d *DependencyServices) ListDependabotAlerts(org, repo string, jsonOutput bool, userPageSize int, fetchAll bool) error {
	pageSize := GetOptimalPageSize(userPageSize)
	path := fmt.Sprintf("repos/%s/%s/dependabot/alerts?per_page=%d", org, repo, pageSize)

	d.alerts = []model.DependabotAlert{}

	for {
		var pageAlerts []model.DependabotAlert
		nextUrl, err := getPages(path, &pageAlerts)
		if err != nil {
			return err
		}

		if jsonOutput {
			d.alerts = append(d.alerts, pageAlerts...)
			if nextUrl == "" {
				break
			}
			path = nextUrl
			continue
		}

		d.alerts = pageAlerts
		if err := d.printTable(); err != nil {
			return err
		}

		if nextUrl == "" {
			break
		}

		if !fetchAll {
			if !AskForNextPage() {
				break
			}
		}
		path = nextUrl
	}

	if jsonOutput {
		return jsonLister(d.alerts)
	}
	return nil
}

// ExportSBOM fetches the CycloneDX SBOM for the repository
func (d *DependencyServices) ExportSBOM(org, repo string) error {
	// The SBOM API returns a large JSON object (CycloneDX format)
	// We map it to interface{} to preserve the exact structure without
	// needing massive struct definitions.
	path := fmt.Sprintf("repos/%s/%s/dependency-graph/sbom", org, repo)

	var sbom map[string]interface{}

	// We can reuse getPages for a single fetch, or just direct client usage.
	// Since getPages handles the Request/Unmarshal logic nicely, let's use it
	// even though there's no pagination for SBOMs.
	_, err := getPages(path, &sbom)
	if err != nil {
		return err
	}

	// Always output SBOM as JSON (it's a data format)
	return jsonLister(sbom)
}

func (d *DependencyServices) printTable() error {
	tp, err := getTablePrinter()
	if err != nil {
		return err
	}

	tp.AddHeader([]string{"ID", "State", "Severity", "Package", "CVE/GHSA", "Version Range"})

	for _, alert := range d.alerts {
		tp.AddField(fmt.Sprintf("%d", alert.Number))
		tp.AddField(alert.State)
		tp.AddField(alert.SecurityAdvisory.Severity)
		tp.AddField(alert.Dependency.Package.Name)

		// Show CVE if available, otherwise GHSA
		id := alert.SecurityAdvisory.CVEId
		if id == "" {
			id = alert.SecurityAdvisory.GHSAId
		}
		tp.AddField(id)

		tp.AddField(alert.SecurityVulnerability.VulnerableVersionRange)
		tp.EndRow()
	}

	return tp.Render()
}

// FetchAllDependabotAlerts retrieves ALL dependabot alerts silently for reporting
func (d *DependencyServices) FetchAllDependabotAlerts(org, repo string) ([]model.DependabotAlert, error) {
	var allAlerts []model.DependabotAlert
	path := fmt.Sprintf("repos/%s/%s/dependabot/alerts?per_page=100", org, repo)

	for {
		var pageAlerts []model.DependabotAlert
		nextUrl, err := getPages(path, &pageAlerts)
		if err != nil {
			return nil, err
		}
		allAlerts = append(allAlerts, pageAlerts...)
		if nextUrl == "" {
			break
		}
		path = nextUrl
	}
	return allAlerts, nil
}
