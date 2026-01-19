package model

type DeleteAnalysis struct {
	NextAnalysisUrl  string `json:"next_analysis_url"`
	ConfirmDeleteUrl string `json:"confirm_delete_url"`
}
