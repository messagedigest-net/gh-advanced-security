//go:build test
// +build test

package services

import (
    "github.com/messagedigest-net/gh-advanced-security/model"
)

// PrintRepoDetails is a test helper that renders repository details using the
// same table printing logic as `RepositoryServices.Show`, but accepts a repo
// directly to avoid network calls. Compiled only with `-tags test`.
func PrintRepoDetails(repo *model.Repository) error {
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

    sas := repo.SecurityAndAnalysis
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
