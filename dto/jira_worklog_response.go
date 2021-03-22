package dto

type JiraWorkLogResponse struct {
	Self   string `json:"self"`
	Author struct {
		Self         string `json:"self"`
		Accountid    string `json:"accountId"`
		Emailaddress string `json:"emailAddress"`
		Avatarurls   struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		Displayname string `json:"displayName"`
		Active      bool   `json:"active"`
		Timezone    string `json:"timeZone"`
		Accounttype string `json:"accountType"`
	} `json:"author"`
	Updateauthor struct {
		Self         string `json:"self"`
		Accountid    string `json:"accountId"`
		Emailaddress string `json:"emailAddress"`
		Avatarurls   struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		Displayname string `json:"displayName"`
		Active      bool   `json:"active"`
		Timezone    string `json:"timeZone"`
		Accounttype string `json:"accountType"`
	} `json:"updateAuthor"`
	Created          string   `json:"created"`
	Updated          string   `json:"updated"`
	Started          string   `json:"started"`
	Timespent        string   `json:"timeSpent"`
	Timespentseconds int      `json:"timeSpentSeconds"`
	ID               string   `json:"id"`
	Issueid          string   `json:"issueId"`
	Errormessages    []string `json:"errorMessages"`
	Errors           struct {
		Timelogged string `json:"timeLogged"`
	} `json:"errors"`
}
