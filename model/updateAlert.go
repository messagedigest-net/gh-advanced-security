package model

type UpdateAlert struct {
	State            string
	DismissedReason  string `json:"dismissed_reason",omitempty`
	DismissedComment string `json:"dismissed_comment",omitempty`
}
