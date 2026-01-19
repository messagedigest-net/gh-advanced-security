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
func (a *AlertServices) ListCodeScanning(org, repo string, jsonOutput bool, userPageSize int) error { // 1. Determine Page Size
    pageSize := GetOptimalPageSize(userPageSize)

    // 2. Prepare URL with pagination
    path := fmt.Sprintf("repos/%s/%s/code-scanning/alerts?per_page=%d", org, repo, pageSize)

    // Reset state for new request
    a.codeAlerts = []model.Alert{}

    for {
        var pageAlerts []model.Alert

        // Fetch one page
        nextUrl, err := getPages(path, &pageAlerts)
        if err != nil {
            return err
        }

        // 3. Handle JSON vs Interactive
        if jsonOutput {
            // For JSON, we fetch ALL pages silently to dump a complete valid JSON object
            a.codeAlerts = append(a.codeAlerts, pageAlerts...)
            if nextUrl == "" {
                break
            }
            path = nextUrl
            continue
        }

        // 4. Interactive Mode: Print THIS page immediately
        a.codeAlerts = pageAlerts // Temporarily set strictly for printing
        if err := a.printCodeScanningTable(); err != nil {
            return err
        }

        // 5. Check continuation
        if nextUrl == "" {
            break
        }

        // 6. Ask User
        if !AskForNextPage() {
            break
        }

        path = nextUrl
    }

    if jsonOutput {
        return jsonLister(a.codeAlerts)
    }

    return nil
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

// In services/alertservices.go

// ListPushProtectionBypasses fetches bypass requests
func (a *AlertServices) ListPushProtectionBypasses(org, repo string, jsonOutput bool) error {
    path := fmt.Sprintf("repos/%s/%s/secret-scanning/push-protection-bypasses", org, repo)

    var bypasses []model.PushProtectionBypass
    // Use the generic getPages we refactored
    _, err := getPages(path, &bypasses)
    if err != nil {
        return err
    }

    if jsonOutput {
        return jsonLister(bypasses)
    }

    tp, err := getTablePrinter()
    if err != nil {
        return err
    }

    tp.AddHeader([]string{"ID", "Secret Type", "Status", "Requester", "Comment", "Date"})
    for _, b := range bypasses {
        tp.AddField(fmt.Sprintf("%d", b.ID))
        tp.AddField(b.SecretType)
        tp.AddField(b.Status)
        tp.AddField(b.Requester.Login)

        // Truncate comment
        comment := b.RequesterComment
        if len(comment) > 40 {
            comment = comment[:37] + "..."
        }
        tp.AddField(comment)
        tp.AddField(b.CreatedAt)
        tp.EndRow()
    }

    return tp.Render()
}
