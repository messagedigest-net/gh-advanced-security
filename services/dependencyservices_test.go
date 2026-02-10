package services

import (
    "bytes"
    "regexp"
    "strings"
    "testing"

    "github.com/messagedigest-net/gh-advanced-security/model"
)

func TestPrintDependabotTable_RendersRow(t *testing.T) {
    d := &DependencyServices{}
    d.alerts = []model.DependabotAlert{
        {
            Number: 777,
            State:  "open",
            Dependency: model.Dependency{
                Package: model.Package{Ecosystem: "npm", Name: "left-pad"},
            },
            SecurityAdvisory: model.SecurityAdvisory{
                GHSAId:  "GHSA-xxxx",
                CVEId:   "CVE-2026-0001",
                Severity: "high",
            },
            SecurityVulnerability: model.SecurityVulnerability{
                VulnerableVersionRange: "<=1.3.0",
            },
        },
    }

    buf := &bytes.Buffer{}
    SetTablePrinterWriter(buf, false, 120)
    defer ClearTablePrinter()

    if err := d.printTable(); err != nil {
        t.Fatalf("printTable returned error: %v", err)
    }

    out := buf.String()
    plain := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(out, "")

    if !strings.Contains(plain, "left-pad") {
        t.Fatalf("expected package name in output, got: %q", plain)
    }
    if !strings.Contains(plain, "CVE-2026-0001") && !strings.Contains(plain, "GHSA-xxxx") {
        t.Fatalf("expected CVE/GHSA id in output, got: %q", plain)
    }
    if !strings.Contains(strings.ToLower(plain), "high") {
        t.Fatalf("expected severity in output, got: %q", plain)
    }
}
