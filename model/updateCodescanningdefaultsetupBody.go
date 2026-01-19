package model

type UpdateCodeScanningDefaultSetup struct {
	State      string   `json:,omitempty`
	QuerySuite string   `json:"query_suite",omitempty`
	Languages  []string `json:,omitempty`
}
