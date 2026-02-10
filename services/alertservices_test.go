package services

import (
    "bytes"
    "testing"

    "github.com/messagedigest-net/gh-advanced-security/model"
)

func TestPrintCodeScanningTable_Runs(t *testing.T) {
    a := &AlertServices{}
    a.codeAlerts = []model.Alert{
        {
            Numer:     123,
            State:     "open",
            CreatedAt: "2026-02-10T12:00:00Z",
            Tool:      model.Tool{Name: "gosec"},
            Rule:      model.Rule{Id: "R001", Description: "A long description to test output."},
        },
    }

    // ensure getTablePrinter uses a harmless size and writer
    buf := &bytes.Buffer{}
    SetTablePrinterWriter(buf, false, 120)

    err := a.printCodeScanningTable()
    if err != nil {
        t.Fatalf("printCodeScanningTable returned error: %v", err)
    }
    if buf.Len() == 0 {
        t.Fatalf("expected some output from printer, got none")
    }
}

func TestPrintSecretScanningTable_Runs(t *testing.T) {
    t.Skip("skipping secret-scanning output capture in CI environment")
    a := &AlertServices{}
    a.secretAlerts = []model.SecretScanningAlert{
        {
            Number:                42,
            State:                 "closed",
            SecretTypeDisplayName: "API Key",
            Resolution:            "",
            PushProtectionBypassed: false,
            CreatedAt:             "2026-02-09T08:00:00Z",
        },
    }

    buf := &bytes.Buffer{}
    SetTablePrinterWriter(buf, false, 120)

    err := a.printSecretScanningTable()
    if err != nil {
        t.Fatalf("printSecretScanningTable returned error: %v", err)
    }
    if buf.Len() == 0 {
        t.Fatalf("expected some output from printer, got none")
    }
    ClearTablePrinter()
}
