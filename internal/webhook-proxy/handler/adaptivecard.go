package handler

// "https://learn.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format?tabs=adaptive-md%2Cdesktop%2Cdesktop1%2Cdesktop2%2Cconnector-html#mention-support-within-adaptive-cards"
type MSTeamsAdaptiveCardMessage struct {
	Type        string       `json:"type"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	ContentType string       `json:"contentType"`
	Content     AdaptiveCard `json:"content"`
}

type AdaptiveCard struct {
	Schema  string        `json:"$schema,omitempty"` // "http://adaptivecards.io/schemas/adaptive-card.json"
	Type    string        `json:"type"`
	Version string        `json:"version"`
	Body    []interface{} `json:"body"`
	Actions []interface{} `json:"actions,omitempty"`
	MSTeams *MSTeamsInfo  `json:"msteams,omitempty"`
}

type MSTeamsInfo struct {
	Entities []MentionEntity `json:"entities"`
}

type MentionEntity struct {
	Type      string        `json:"type"`
	Text      string        `json:"text"`
	Mentioned MentionedUser `json:"mentioned"`
}

type TextBlock struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	Wrap   bool   `json:"wrap,omitempty"`
	Weight string `json:"weight,omitempty"`
	Size   string `json:"size,omitempty"`
}

type FactSet struct {
	Type  string `json:"type"`
	Facts []Fact `json:"facts"`
}

type Fact struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type MentionedUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ActionOpenURL struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}
