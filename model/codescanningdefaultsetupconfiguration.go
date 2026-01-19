package model

type CodeScanningDefaultSetupConfiguration struct {
	State      string
	Languages  []string
	QuerySuite string `json:"query_suite"`
	UpdatedAt  string `json:"updated_at"`
	Schedule   string
}
