package model

type UpdateCodeScanningDefaultSetup struct {
	State      string   `json:"state,omitempty"`
	QuerySuite string   `json:"query_suite,omitempty"`
	Languages  []string `json:"languages,omitempty"`
}
