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
func (a *AlertServices) ListCodeScanning(org, repo string, jsonOutput bool, userPageSize int, fetchAll bool) error {
    pageSize := GetOptimalPageSize(userPageSize)
    path := fmt.Sprintf("repos/%s/%s/code-scanning/alerts?per_page=%d", org, repo, pageSize)

    a.codeAlerts = []model.Alert{}

    for {
        var pageAlerts []model.Alert
        nextUrl, err := getPages(path, &pageAlerts)
        if err != nil {
            return err
        }

        if jsonOutput {
            a.codeAlerts = append(a.codeAlerts, pageAlerts...)
            if nextUrl == "" {
                break
            }
            path = nextUrl
            continue
        }

        // Interactive Render
        a.codeAlerts = pageAlerts
        if err := a.printCodeScanningTable(); err != nil {
            return err
        }

        if nextUrl == "" {
            break
        }

        // LOGIC: Pause only if NOT fetching all
        if !fetchAll {
            if !AskForNextPage() {
                break
            }
        }
        path = nextUrl
    }

    if jsonOutput {
        return jsonLister(a.codeAlerts)
    }
    return nil
}

// ListSecretScanning fetches and displays Secret Scanning alerts
func (a *AlertServices) ListSecretScanning(org, repo string, jsonOutput bool, userPageSize int, fetchAll bool) error {
    pageSize := GetOptimalPageSize(userPageSize)
    path := fmt.Sprintf("repos/%s/%s/secret-scanning/alerts?per_page=%d", org, repo, pageSize)

    a.secretAlerts = []model.SecretScanningAlert{}

    for {
        var pageAlerts []model.SecretScanningAlert
        nextUrl, err := getPages(path, &pageAlerts)
        if err != nil {
            return err
        }

        if jsonOutput {
            a.secretAlerts = append(a.secretAlerts, pageAlerts...)
            if nextUrl == "" {
                break
            }
            path = nextUrl
            continue
        }

        a.secretAlerts = pageAlerts
        if err := a.printSecretScanningTable(); err != nil {
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
        return jsonLister(a.secretAlerts)
    }
    return nil
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
func (a *AlertServices) ListPushProtectionBypasses(org, repo string, jsonOutput bool, userPageSize int, fetchAll bool) error {
    pageSize := GetOptimalPageSize(userPageSize)
    path := fmt.Sprintf("repos/%s/%s/secret-scanning/push-protection-bypasses?per_page=%d", org, repo, pageSize)

    var allBypasses []model.PushProtectionBypass // We need a local accumulator or struct field

    for {
        var pageBypasses []model.PushProtectionBypass
        nextUrl, err := getPages(path, &pageBypasses)
        if err != nil {
            return err
        }

        if jsonOutput {
            allBypasses = append(allBypasses, pageBypasses...)
            if nextUrl == "" {
                break
            }
            path = nextUrl
            continue
        }

        // For bypasses, we don't have a struct field to store them temporarily in the service
        // (unless you added one), so we can pass the slice directly to a helper or render inline.
        // Assuming you add 'printBypassTable(bypasses)' helper:
        if err := a.printBypassTable(pageBypasses); err != nil {
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
        return jsonLister(allBypasses)
    }
    return nil
}

// Helper for Bypasses
func (a *AlertServices) printBypassTable(bypasses []model.PushProtectionBypass) error {
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
