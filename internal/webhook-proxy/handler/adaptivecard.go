package handler

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
}

type TextBlock struct {
	Type    string        `json:"type"`
	Text    string        `json:"text"`
	Wrap    bool          `json:"wrap,omitempty"`
	Weight  string        `json:"weight,omitempty"`
	Size    string        `json:"size,omitempty"`
	Inlines []interface{} `json:"inlines,omitempty"`
}

type FactSet struct {
	Type  string `json:"type"`
	Facts []Fact `json:"facts"`
}

type Fact struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type MentionText struct {
	Type      string        `json:"type"`
	Text      string        `json:"text"`
	Mentioned MentionedUser `json:"mentioned"`
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
