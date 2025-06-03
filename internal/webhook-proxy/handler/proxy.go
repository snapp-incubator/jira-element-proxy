package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/snapp-incubator/jira-element-proxy/internal/config"
	"github.com/snapp-incubator/jira-element-proxy/internal/webhook-proxy/request"
)

const (
	PLATFORM_SUBTEAM = "platform"
	NETWORK_SUBTEAM  = "network"
	RUNTIME_SUBTEAM  = "runtime"
)

type (
	Proxy struct {
		MSTeamsConf config.MSTeamsConfig
	}
)

func (p *Proxy) ProxyToMSTeamsHandler(isComment bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		subteam := c.Param("team")
		req := &request.JiraRequest{}

		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Errorf("Failed to read request body: %s", err.Error())
			return c.String(http.StatusInternalServerError, "Error reading body")
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err = c.Bind(req)
		if err != nil {
			logrus.Errorf("failed to bind request body: %s", err.Error())
			logrus.Debugf("Problematic request body for bind: %s", string(bodyBytes))
			return c.NoContent(http.StatusBadRequest)
		}

		generatedText := generateTeamsAdaptiveCard(req, isComment)
		logrus.Printf("Team: %s, Generated Text: %s", subteam, generatedText)

		var targetTeamsURL string
		switch subteam {
		case PLATFORM_SUBTEAM:
			targetTeamsURL = p.MSTeamsConf.PlatformURL
			logrus.Printf("Using MS Teams platform url: %s", targetTeamsURL)
		case NETWORK_SUBTEAM:
			targetTeamsURL = p.MSTeamsConf.NetworkURL
			logrus.Printf("Using MS Teams network url: %s", targetTeamsURL)
		case RUNTIME_SUBTEAM:
			targetTeamsURL = p.MSTeamsConf.RuntimeURL
			logrus.Printf("Using MS Teams runtime url: %s", targetTeamsURL)
		default:
			targetTeamsURL = p.MSTeamsConf.URL
			logrus.Printf("Using default MS Teams url: %s for team param '%s'", targetTeamsURL, subteam)
		}

		if targetTeamsURL == "" {
			logrus.Warnf("No MS Teams URL configured for team '%s' (or default). Skipping notification.", subteam)

			return c.NoContent(http.StatusOK)
		}

		if p.sendToMSTeams(generatedText, targetTeamsURL) {
			return c.NoContent(http.StatusOK)
		}

		return c.NoContent(http.StatusInternalServerError)
	}
}

func (p *Proxy) sendToMSTeams(card AdaptiveCard, webhookURL string) bool {
	payload := MSTeamsAdaptiveCardMessage{
		Type: "message", // Standard type for messages with attachments
		Attachments: []Attachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content:     card,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("MS Teams (AdaptiveCard): marshal request body error: %s", err)
		return false
	}

	logrus.Debugf("MS Teams (AdaptiveCard): Sending JSON payload: %s", string(body))

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		logrus.Errorf("MS Teams (AdaptiveCard): create request error: %s for URL %s", err, webhookURL)
		return false
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("MS Teams (AdaptiveCard): request error: %s for URL %s", err, webhookURL)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logrus.Infof("MS Teams (AdaptiveCard): successfully sent webhook to %s", webhookURL)
		return true
	}

	responseBodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logrus.Errorf("MS Teams (AdaptiveCard): response body read error: %s (after non-success status %d from %s)", readErr, resp.StatusCode, webhookURL)
		return false
	}
	logrus.Errorf("MS Teams (AdaptiveCard): failed to send webhook. Status: %d, URL: %s, Response: %s", resp.StatusCode, webhookURL, string(responseBodyBytes))
	return false
}

func generateTeamsAdaptiveCard(req *request.JiraRequest, isComment bool) AdaptiveCard {
	creatorDisplayName := "N/A"
	creatorMentionID := ""

	if req.Fields.Creator.EmailAddress != "" {
		creatorMentionID = req.Fields.Creator.EmailAddress
		if req.Fields.Creator.DisplayName != "" {
			creatorDisplayName = req.Fields.Creator.DisplayName
		} else if req.Fields.Creator.Name != "" {
			creatorDisplayName = req.Fields.Creator.Name
		} else {
			creatorDisplayName = req.Fields.Creator.EmailAddress
		}
	} else if req.Fields.Creator.DisplayName != "" {
		creatorDisplayName = req.Fields.Creator.DisplayName
	} else if req.Fields.Creator.Name != "" {
		creatorDisplayName = req.Fields.Creator.Name
	}

	assigneeDisplayName := "N/A"
	assigneeMentionID := ""
	if req.Fields.Assignee.EmailAddress != "" {
		assigneeMentionID = req.Fields.Assignee.EmailAddress
		if req.Fields.Assignee.DisplayName != "" {
			assigneeDisplayName = req.Fields.Assignee.DisplayName
		} else if req.Fields.Assignee.Name != "" {
			assigneeDisplayName = req.Fields.Assignee.Name
		} else {
			assigneeDisplayName = req.Fields.Assignee.EmailAddress
		}
	} else if req.Fields.Assignee.DisplayName != "" {
		assigneeDisplayName = req.Fields.Assignee.DisplayName
	} else if req.Fields.Assignee.Name != "" {
		assigneeDisplayName = req.Fields.Assignee.Name
	}

	requestTypeName := "N/A"
	webLink := "#"
	if req.Fields.CustomField10003.RequestType.Name != "" {
		requestTypeName = req.Fields.CustomField10003.RequestType.Name
	}
	if req.Fields.CustomField10003.Links.Web != "" {
		webLink = req.Fields.CustomField10003.Links.Web
	}

	summary := "N/A"
	if req.Fields.Summary != "" {
		summary = req.Fields.Summary
	}

	var title string
	if isComment {
		title = "📰 New Comment Added"
	} else {
		title = "🎯 New Issue/Update"
	}

	cardBody := []interface{}{
		TextBlock{Type: "TextBlock", Text: title, Weight: "bolder", Size: "medium", Wrap: true},
		FactSet{Type: "FactSet", Facts: []Fact{
			{Title: "Type:", Value: requestTypeName},
			{Title: "Summary:", Value: summary},
		}},
	}

	issuerTextBlock := TextBlock{Type: "TextBlock", Wrap: true}
	if creatorMentionID != "" {
		issuerTextBlock.Text = fmt.Sprintf("**Issuer:** <at>%s</at>", creatorDisplayName)
		issuerTextBlock.Inlines = []interface{}{
			MentionText{
				Type: "mention",
				Text: fmt.Sprintf("<at>%s</at>", creatorDisplayName),
				Mentioned: MentionedUser{
					ID:   creatorMentionID,
					Name: creatorDisplayName,
				},
			},
		}
	} else {
		issuerTextBlock.Text = fmt.Sprintf("**Issuer:** %s", creatorDisplayName)
	}
	cardBody = append(cardBody, issuerTextBlock)

	assigneeTextBlock := TextBlock{Type: "TextBlock", Wrap: true}
	if assigneeMentionID != "" {
		assigneeTextBlock.Text = fmt.Sprintf("**Assignee:** <at>%s</at>", assigneeDisplayName)
		assigneeTextBlock.Inlines = []interface{}{
			MentionText{
				Type: "mention",
				Text: fmt.Sprintf("<at>%s</at>", assigneeDisplayName),
				Mentioned: MentionedUser{
					ID:   assigneeMentionID,
					Name: assigneeDisplayName,
				},
			},
		}
	} else {
		assigneeTextBlock.Text = fmt.Sprintf("**Assignee:** %s", assigneeDisplayName)
	}
	cardBody = append(cardBody, assigneeTextBlock)

	return AdaptiveCard{
		Type:    "AdaptiveCard",
		Version: "1.5",
		Body:    cardBody,
		Actions: []interface{}{
			ActionOpenURL{Type: "Action.OpenUrl", Title: "View Issue in Jira", URL: webLink},
		},
	}
}
