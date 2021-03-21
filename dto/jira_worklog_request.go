package dto

type JiraWorklogRequest struct {
	Timespent string  `json:"timeSpent"`
	Comment   Comment `json:"comment"`
	Started   string  `json:"started"`
}

type Comment struct {
	Type    string        `json:"type"`
	Version int           `json:"version"`
	Content []string `json:"content"`
}
