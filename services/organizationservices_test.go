package services

import (
    "bytes"
    "regexp"
    "strings"
    "testing"

    "github.com/messagedigest-net/gh-advanced-security/model"
)

func TestPrintOrgTable_RendersRow(t *testing.T) {
    o := &OrganizationServices{}
    o.organizations = []model.Organization{
        {
            Login:  "example-org",
            URL:    "https://github.com/example-org",
            DependencyGraphEnabledForNewRepositories:              true,
            DependabotAlertsEnabledForNewRepositories:             false,
            DependabotSecurityUpdatesEnabledForNewRepositories:    true,
            AdvancedSecurityEnabledForNewRepositories:             true,
            SecretScanningEnabledForNewRepositories:               false,
            SecretScanningPushProtectionEnabledForNewRepositories: true,
        },
    }

    buf := &bytes.Buffer{}
    SetTablePrinterWriter(buf, false, 120)
    defer ClearTablePrinter()

    if err := o.printOrgTable(); err != nil {
        t.Fatalf("printOrgTable returned error: %v", err)
    }

    out := buf.String()
    plain := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(out, "")

    if !strings.Contains(plain, "example-org") {
        t.Fatalf("expected org login in output, got: %q", plain)
    }
    if !strings.Contains(plain, "https://github.com/example-org") {
        t.Fatalf("expected org URL in output, got: %q", plain)
    }
    if !strings.Contains(plain, "Enabled") || !strings.Contains(plain, "Disabled") {
        t.Fatalf("expected enabled/disabled markers in output, got: %q", plain)
    }
}
