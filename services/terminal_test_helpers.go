//go:build test
// +build test

package services

import (
    "io"

    "github.com/cli/go-gh/v2/pkg/tableprinter"
)

// SetTablePrinterWriter configures the test hook to return a fresh printer
// that writes to `out`. This file is only compiled when tests are run with
// the `-tags test` build tag.
func SetTablePrinterWriter(out io.Writer, isTTY bool, width int) {
    testTablePrinter = func() (tableprinter.TablePrinter, bool) {
        return tableprinter.New(out, isTTY, width), true
    }
}

// ClearTablePrinter clears the test hook.
func ClearTablePrinter() {
    testTablePrinter = nil
}
