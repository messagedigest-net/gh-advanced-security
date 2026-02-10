package services

import (
    "regexp"
    "strings"
    "testing"
)

// stripANSI removes common ANSI color sequences so tests can assert on plain text
func stripANSI(s string) string {
    re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
    return re.ReplaceAllString(s, "")
}

func TestSeverityWithIcon_ContainsIconAndText(t *testing.T) {
    cases := map[string]string{
        "high":   "❗",
        "error":  "❗",
        "medium": "⚠",
        "warning":"⚠",
        "low":    "ℹ",
        "info":   "ℹ",
    }

    for sev, icon := range cases {
        out := SeverityWithIcon(sev)
        plain := stripANSI(out)
        if !strings.Contains(plain, icon) {
            t.Fatalf("severity %q: expected icon %q in %q", sev, icon, plain)
        }
        // ensure severity text is present (case-insensitive)
        if !strings.Contains(strings.ToLower(plain), strings.ToLower(sev)) {
            t.Fatalf("severity %q: expected text %q in %q", sev, sev, plain)
        }
    }
}

func TestStateColored_FormatsStateText(t *testing.T) {
    states := []string{"open", "closed", "dismissed", "resolved", "fixed"}
    for _, s := range states {
        out := StateColored(s)
        plain := stripANSI(out)
        // Expect the capitalized form to be present
        expected := strings.Title(strings.ToLower(s))
        if !strings.Contains(plain, expected) {
            t.Fatalf("state %q: expected %q in %q", s, expected, plain)
        }
    }
}
