package services

import (
    "fmt"

    "github.com/messagedigest-net/gh-advanced-security/model"
)

var alertSvcs *AlertServices

type AlertServices struct {
    codeAlerts   []model.Alert
    secretAlerts []model.SecretScanningAlert
    next         string
}

func GetAlertServices() *AlertServices {
    if alertSvcs == nil {
        alertSvcs = &AlertServices{}
    }
    return alertSvcs
}

// ListCodeScanning fetches and displays Code Scanning alerts
func (a *AlertServices) ListCodeScanning(org, repo string, jsonOutput bool) error {
    // Reset state for new request
    a.codeAlerts = []model.Alert{}
    a.next = ""

    path := fmt.Sprintf("repos/%s/%s/code-scanning/alerts", org, repo)

    // Pagination Loop
    for {
        var pageAlerts []model.Alert

        // Use the new generic getPages[T]
        nextUrl, err := getPages(path, &pageAlerts)
        if err != nil {
            return err
        }

        a.codeAlerts = append(a.codeAlerts, pageAlerts...)

        if nextUrl == "" {
            break
        }
        path = nextUrl
    }

    if jsonOutput {
        return jsonLister(a.codeAlerts)
    }

    return a.printCodeScanningTable()
}

// ListSecretScanning fetches and displays Secret Scanning alerts
func (a *AlertServices) ListSecretScanning(org, repo string, jsonOutput bool) error {
    // Reset state
    a.secretAlerts = []model.SecretScanningAlert{}
    a.next = ""

    path := fmt.Sprintf("repos/%s/%s/secret-scanning/alerts", org, repo)

    // Pagination Loop
    for {
        var pageAlerts []model.SecretScanningAlert

        nextUrl, err := getPages(path, &pageAlerts)
        if err != nil {
            return err
        }

        a.secretAlerts = append(a.secretAlerts, pageAlerts...)

        if nextUrl == "" {
            break
        }
        path = nextUrl
    }

    if jsonOutput {
        return jsonLister(a.secretAlerts)
    }

    return a.printSecretScanningTable()
}

// Helper to print Code Scanning table
func (a *AlertServices) printCodeScanningTable() error {
    tp, err := getTablePrinter()
    if err != nil {
        return err
    }

    tp.AddHeader([]string{"ID", "State", "Tool", "Rule ID", "Description", "Created At"})

    for _, alert := range a.codeAlerts {
        tp.AddField(fmt.Sprintf("%d", alert.Numer)) // Note: 'Numer' matches your existing model
        tp.AddField(alert.State)
        tp.AddField(alert.Tool.Name)
        tp.AddField(alert.Rule.Id)

        // Truncate description if too long
        desc := alert.Rule.Description
        if len(desc) > 50 {
            desc = desc[:47] + "..."
        }
        tp.AddField(desc)
        tp.AddField(alert.CreatedAt)
        tp.EndRow()
    }

    return tp.Render()
}

// Helper to print Secret Scanning table
func (a *AlertServices) printSecretScanningTable() error {
    tp, err := getTablePrinter()
    if err != nil {
        return err
    }

    tp.AddHeader([]string{"ID", "State", "Secret Type", "Resolution", "Push Protection", "Created At"})

    for _, alert := range a.secretAlerts {
        tp.AddField(fmt.Sprintf("%d", alert.Number))
        tp.AddField(alert.State)
        tp.AddField(alert.SecretTypeDisplayName)

        resolution := alert.Resolution
        if resolution == "" {
            resolution = "-"
        }
        tp.AddField(resolution)

        tp.AddField(enabledOrDisabled(alert.PushProtectionBypassed))
        tp.AddField(alert.CreatedAt)
        tp.EndRow()
    }

    return tp.Render()
}
