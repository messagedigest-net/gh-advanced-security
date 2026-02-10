package services

import (
    "bytes"
    "regexp"
    "strings"
    "testing"

    "github.com/messagedigest-net/gh-advanced-security/model"
)

func TestPrintRepoTable_RendersRow(t *testing.T) {
    r := &RepositoryServices{}
    r.repositories = []model.Repository{
        {
            FullName: "example-org/example-repo",
            Owner:    model.Organization{Login: "example-org"},
            Private:  true,
            HtmlUrl:  "https://github.com/example-org/example-repo",
            Description: "Test repo",
            Homepage: "https://example.org",
            Language: "Go",
            Topics: []string{"cli", "security"},
            SecurityAndAnalysis: model.SecurityAndAnalysis{
                AdvancedSecurity: model.Status{Status: "enabled"},
                SecretScanning: model.Status{Status: "disabled"},
                SecretScanningNonProviderPatterns: model.Status{Status: "configured"},
                SecretScanningValidityChecks: model.Status{Status: "ok"},
                SecretScanningPushProtection: model.Status{Status: "bypassed"},
                DependabotSecurityUpdates: model.Status{Status: "enabled"},
            },
        },
    }

    buf := &bytes.Buffer{}
    SetTablePrinterWriter(buf, false, 160)
    defer ClearTablePrinter()

    if err := r.printRepoTable(); err != nil {
        t.Fatalf("printRepoTable returned error: %v", err)
    }

    out := buf.String()
    plain := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(out, "")

    if !strings.Contains(plain, "example-org/example-repo") {
        t.Fatalf("expected full name in output, got: %q", plain)
    }
    if !strings.Contains(plain, "example-org") {
        t.Fatalf("expected owner in output, got: %q", plain)
    }
    if !strings.Contains(plain, "Test repo") {
        t.Fatalf("expected description in output, got: %q", plain)
    }
    if !strings.Contains(plain, "Go") {
        t.Fatalf("expected language in output, got: %q", plain)
    }
}

func TestShow_RendersFields(t *testing.T) {
    repo := &model.Repository{
        FullName: "example-org/example-repo",
        Visibility: "public",
        HtmlUrl: "https://github.com/example-org/example-repo",
        SecurityAndAnalysis: model.SecurityAndAnalysis{
            AdvancedSecurity: model.Status{Status: "enabled"},
            SecretScanning: model.Status{Status: "disabled"},
            SecretScanningNonProviderPatterns: model.Status{Status: ""},
            SecretScanningValidityChecks: model.Status{Status: ""},
            SecretScanningPushProtection: model.Status{Status: ""},
            DependabotSecurityUpdates: model.Status{Status: "enabled"},
        },
    }

    buf := &bytes.Buffer{}
    SetTablePrinterWriter(buf, false, 120)
    defer ClearTablePrinter()

    if err := PrintRepoDetails(repo); err != nil {
        t.Fatalf("PrintRepoDetails returned error: %v", err)
    }

    plain := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(buf.String(), "")
    if !strings.Contains(plain, "example-org/example-repo") || !strings.Contains(plain, "public") {
        t.Fatalf("expected show output to include full name and visibility, got: %q", plain)
    }
    if !strings.Contains(plain, "Advanced Security") || !strings.Contains(plain, "enabled") {
        t.Fatalf("expected security settings in output, got: %q", plain)
    }
}
