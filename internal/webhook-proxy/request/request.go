package request

type JiraRequest struct {
	Self      string    `json:"self"`
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	Changelog Changelog `json:"changelog"`
	Fields    Fields    `json:"fields"`
}

type Changelog struct {
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Histories  interface{} `json:"histories"`
}

type Fields struct {
	FixVersions      []interface{}    `json:"fixVersions"`
	Priority         Priority         `json:"priority"`
	Labels           []interface{}    `json:"labels"`
	Assignee         User             `json:"assignee"`
	Status           Status           `json:"status"`
	Components       []interface{}    `json:"components"`
	Creator          User             `json:"creator"`
	Subtasks         []interface{}    `json:"subtasks"`
	Reporter         User             `json:"reporter"`
	Progress         Progress         `json:"progress"`
	Votes            Votes            `json:"votes"`
	IssueType        IssueType        `json:"issuetype"`
	Project          Project          `json:"project"`
	Watches          Watches          `json:"watches"`
	Created          int64            `json:"created"`
	Updated          int64            `json:"updated"`
	TimeTracking     TimeTracking     `json:"timetracking"`
	CustomField10405 CustomField      `json:"customfield_10405"`
	CustomField10406 []string         `json:"customfield_10406"`
	Summary          string           `json:"summary"`
	CustomField10003 CustomField10003 `json:"customfield_10003"`
}

type Priority struct {
	Self       string `json:"self"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IconURL    string `json:"iconUrl"`
	NamedValue string `json:"namedValue"`
}

type User struct {
	Self         string     `json:"self"`
	Name         string     `json:"name"`
	Key          string     `json:"key"`
	EmailAddress string     `json:"emailAddress"`
	AvatarURLs   AvatarURLs `json:"avatarUrls"`
	DisplayName  string     `json:"displayName"`
	Active       bool       `json:"active"`
	TimeZone     string     `json:"timeZone"`
}

type AvatarURLs struct {
	Four8X48  string `json:"48x48"`
	Two4X24   string `json:"24x24"`
	One6X16   string `json:"16x16"`
	Three2X32 string `json:"32x32"`
}

type Status struct {
	Self           string         `json:"self"`
	Description    string         `json:"description"`
	IconURL        string         `json:"iconUrl"`
	Name           string         `json:"name"`
	ID             int            `json:"id"`
	StatusCategory StatusCategory `json:"statusCategory"`
}

type StatusCategory struct {
	Self      string `json:"self"`
	ID        int    `json:"id"`
	Key       string `json:"key"`
	ColorName string `json:"colorName"`
	Name      string `json:"name"`
}

type Progress struct {
	Progress int `json:"progress"`
	Total    int `json:"total"`
}

type Votes struct {
	Self     string `json:"self"`
	Votes    int    `json:"votes"`
	HasVoted bool   `json:"hasVoted"`
}

type IssueType struct {
	Self        string `json:"self"`
	ID          int    `json:"id"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
	Name        string `json:"name"`
	Subtask     bool   `json:"subtask"`
	NamedValue  string `json:"namedValue"`
}

type Project struct {
	Self           string     `json:"self"`
	ID             int        `json:"id"`
	Key            string     `json:"key"`
	Name           string     `json:"name"`
	AvatarURLs     AvatarURLs `json:"avatarUrls"`
	ProjectTypeKey string     `json:"projectTypeKey"`
	Simplified     bool       `json:"simplified"`
}

type Watches struct {
	Self       string `json:"self"`
	WatchCount int    `json:"watchCount"`
	IsWatching bool   `json:"isWatching"`
}

type TimeTracking struct {
	OriginalEstimateSeconds  int `json:"originalEstimateSeconds"`
	RemainingEstimateSeconds int `json:"remainingEstimateSeconds"`
	TimeSpentSeconds         int `json:"timeSpentSeconds"`
}

type CustomField struct {
	Self     string `json:"self"`
	Value    string `json:"value"`
	ID       string `json:"id"`
	Disabled bool   `json:"disabled"`
}

type CustomField10003 struct {
	Links         Links         `json:"_links"`
	RequestType   RequestType   `json:"requestType"`
	CurrentStatus CurrentStatus `json:"currentStatus"`
}

type Links struct {
	JiraRest string `json:"jiraRest"`
	Web      string `json:"web"`
	Self     string `json:"self"`
}

type RequestType struct {
	ID            string   `json:"id"`
	Links         Links    `json:"_links"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	HelpText      string   `json:"helpText"`
	ServiceDeskID string   `json:"serviceDeskId"`
	GroupIDs      []string `json:"groupIds"`
	Icon          Icon     `json:"icon"`
}

type Icon struct {
	ID    string `json:"id"`
	Links Links  `json:"_links"`
}

type CurrentStatus struct {
	Status     string     `json:"status"`
	StatusDate StatusDate `json:"statusDate"`
}

type StatusDate struct {
	ISO8601     string `json:"iso8601"`
	Jira        string `json:"jira"`
	Friendly    string `json:"friendly"`
	EpochMillis int64  `json:"epochMillis"`
}
