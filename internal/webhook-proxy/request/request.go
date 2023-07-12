package request

type Jira struct {
	Timestamp          int64  `json:"timestamp"`
	WebhookEvent       string `json:"webhookEvent"`
	IssueEventTypeName string `json:"issue_event_type_name"`
	User               struct {
		Self         string `json:"self"`
		Name         string `json:"name"`
		Key          string `json:"key"`
		EmailAddress string `json:"emailAddress"`
		AvatarUrls   struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		DisplayName string `json:"displayName"`
		Active      bool   `json:"active"`
		TimeZone    string `json:"timeZone"`
	} `json:"user"`
	Issue struct {
		ID     string `json:"id"`
		Self   string `json:"self"`
		Key    string `json:"key"`
		Fields struct {
			CustomField11401 struct {
				Links struct {
					Web string `json:"web"`
				} `json:"_links"`
				RequestType struct {
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"requestType"`
			} `json:"customfield_11401"`
			Summary   string `json:"summary"`
			IssueType struct {
				Description string `json:"description"`
				Name        string `json:"name"`
			} `json:"issuetype"`
		} `json:"fields"`
	} `json:"issue"`
}
