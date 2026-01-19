package model

type Analysis struct {
	Instance
	Error        string
	CreatedAt    string `json:"created_at"`
	ResultsCount int    `json:"results_count"`
	ID           int
	URL          string
	SarifID      string `json:"sarif_id"`
	Tool         Tool
	Deletable    bool
	Warning      string
}
