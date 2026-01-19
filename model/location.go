package model

type Location struct {
	Path        string
	StartLine   int `json:"start_line"`
	EndLine     int `json:"end_line"`
	StartColumn int `json:"start_column"`
	EndColumn   int `json:"end_column"`
}
