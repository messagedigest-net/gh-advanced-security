package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate security reports",
	Long:  `Generate CSV reports of security alerts across an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		services.ChooseSubCommand(cmd.Commands(), args, "What kind of report do you want?")
	},
}

var codeScanningReportCmd = &cobra.Command{
	Use:   "code-scanning",
	Short: "Export Code Scanning alerts to CSV",
	Run: func(cmd *cobra.Command, args []string) {
		generateReport(cmd, args, "code-scanning")
	},
}

var secretScanningReportCmd = &cobra.Command{
	Use:   "secret-scanning",
	Short: "Export Secret Scanning alerts to CSV",
	Run: func(cmd *cobra.Command, args []string) {
		generateReport(cmd, args, "secret-scanning")
	},
}

// NEW: Dependabot Report Command
var dependabotReportCmd = &cobra.Command{
	Use:   "dependabot",
	Short: "Export Dependabot alerts to CSV",
	Run: func(cmd *cobra.Command, args []string) {
		generateReport(cmd, args, "dependabot")
	},
}

// Shared logic for generating reports
func generateReport(cmd *cobra.Command, args []string, reportType string) {
	target, _ := services.GetTarget(cmd, args, "Which organization?")

	repoSvc := services.GetRepositoryServices()
	fmt.Printf("Fetching repositories for %s...\n", target)
	repos, err := repoSvc.FetchAllForOrg(target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Analyzing %d repositories. This may take a while...\n", len(repos))

	// CSV File Setup
	filename := fmt.Sprintf("%s-%s-report.csv", target, reportType)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write Headers based on type
	switch reportType {
	case "code-scanning":
		writer.Write([]string{"Repository", "Tool", "Rule", "Severity", "State", "Created At", "URL"})
	case "secret-scanning":
		writer.Write([]string{"Repository", "Secret Type", "Secret", "State", "Resolution", "Created At", "URL"})
	case "dependabot":
		writer.Write([]string{"Repository", "Package", "Severity", "State", "CVE/GHSA", "Vulnerable Version", "Created At", "URL"})
	}

	// Worker Pool setup
	type Row []string
	results := make(chan Row)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5)

	for _, repo := range repos {
		wg.Add(1)
		go func(repoName string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			// Logic Dispatcher
			if reportType == "code-scanning" {
				svc := services.GetAlertServices()
				alerts, err := svc.FetchAllCodeScanning(target, repoName)
				if err == nil {
					for _, a := range alerts {
						results <- []string{
							repoName, a.Tool.Name, a.Rule.Id, a.Rule.Severity, a.State, a.CreatedAt, a.HtmlUrl,
						}
					}
				}
			} else if reportType == "secret-scanning" {
				svc := services.GetAlertServices()
				alerts, err := svc.FetchAllSecretScanning(target, repoName)
				if err == nil {
					for _, a := range alerts {
						results <- []string{
							repoName, a.SecretType, a.Secret, a.State, a.Resolution, a.CreatedAt, a.HtmlUrl,
						}
					}
				}
			} else if reportType == "dependabot" {
				// NEW: Dependabot Logic
				svc := services.GetDependencyServices()
				alerts, err := svc.FetchAllDependabotAlerts(target, repoName)
				if err == nil {
					for _, a := range alerts {
						// Fallback logic for Identifier (CVE vs GHSA)
						id := a.SecurityAdvisory.CVEId
						if id == "" {
							id = a.SecurityAdvisory.GHSAId
						}

						results <- []string{
							repoName,
							a.Dependency.Package.Name,
							a.SecurityAdvisory.Severity,
							a.State,
							id,
							a.SecurityVulnerability.VulnerableVersionRange,
							a.CreatedAt,
							a.HtmlUrl,
						}
					}
				}
			}
			time.Sleep(50 * time.Millisecond)
		}(repo.Name)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for row := range results {
		writer.Write(row)
		count++
		if count%100 == 0 {
			fmt.Printf("\rProcessed %d alerts...", count)
		}
	}

	fmt.Printf("\nDone! Report saved to %s\n", filename)
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.AddCommand(codeScanningReportCmd)
	reportCmd.AddCommand(secretScanningReportCmd)
	reportCmd.AddCommand(dependabotReportCmd)
}
