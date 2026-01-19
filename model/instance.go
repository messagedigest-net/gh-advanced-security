package model

type Instance struct {
	Ref             string
	AnalysisKey     string `json:"analysis_key"`
	Category        string
	Environment     string
	State           string
	CommitSha       string `json:"commit_sha"`
	Message         Message
	Location        Location
	Classifications []Classification
}
