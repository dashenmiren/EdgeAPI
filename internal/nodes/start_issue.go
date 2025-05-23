// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package nodes

type StartIssue struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

func NewStartIssue(code string, message string, suggestion string) *StartIssue {
	return &StartIssue{
		Code:       code,
		Message:    message,
		Suggestion: suggestion,
	}
}
