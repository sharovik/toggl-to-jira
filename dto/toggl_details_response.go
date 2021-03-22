package dto

import "time"

type TogglDetailsResponse struct {
	TotalGrand      int `json:"total_grand"`
	TotalBillable   int `json:"total_billable"`
	TotalCurrencies []struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	} `json:"total_currencies"`
	TotalCount int        `json:"total_count"`
	PerPage    int        `json:"per_page"`
	Data       []DataItem `json:"data"`
}

type DataItem struct {
	ID              int           `json:"id"`
	Pid             int           `json:"pid"`
	Tid             interface{}   `json:"tid"`
	UID             int           `json:"uid"`
	Description     string        `json:"description"`
	Start           time.Time     `json:"start"`
	End             time.Time     `json:"end"`
	Updated         time.Time     `json:"updated"`
	Dur             int64         `json:"dur"`
	User            string        `json:"user"`
	UseStop         bool          `json:"use_stop"`
	Client          string        `json:"client"`
	Project         string        `json:"project"`
	ProjectColor    string        `json:"project_color"`
	ProjectHexColor string        `json:"project_hex_color"`
	Task            interface{}   `json:"task"`
	Billable        float64       `json:"billable"`
	IsBillable      bool          `json:"is_billable"`
	Cur             string        `json:"cur"`
	Tags            []interface{} `json:"tags"`
}
