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

	MSTeamsMessage struct {
		Text string `json:"text"`
	}
)

func (p *Proxy) ProxyToMSTeamsHandler(isComment bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		subteam := c.Param("team")
		req := &request.JiraRequest{} // Your existing JiraRequest struct

		// Read and re-buffer body for binding
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

		// Re-use existing text generation, assuming it's suitable for Teams plain text
		generatedText := generateTeamsTextMessage(req, isComment)
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
		default: // Includes empty subteam (e.g. POST / or POST /comment)
			targetTeamsURL = p.MSTeamsConf.URL
			logrus.Printf("Using default MS Teams url: %s for team param '%s'", targetTeamsURL, subteam)
		}

		if targetTeamsURL == "" {
			logrus.Warnf("No MS Teams URL configured for team '%s' (or default). Skipping notification.", subteam)
			// Return OK because the request was processed, but no action taken for this part.
			// Or, if a URL is mandatory, return an error.
			return c.NoContent(http.StatusOK)
		}

		if p.sendToMSTeams(generatedText, targetTeamsURL) {
			return c.NoContent(http.StatusOK)
		}

		return c.NoContent(http.StatusInternalServerError)
	}
}

// sendToMSTeams is the new function to send messages to a specific MS Teams webhook URL
func (p *Proxy) sendToMSTeams(textPayload string, webhookURL string) bool {
	message := MSTeamsMessage{
		Text: textPayload,
	}

	body, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("MS Teams: marshal request body error: %s", err)
		return false
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		logrus.Errorf("MS Teams: create request error: %s for URL %s", err, webhookURL)
		return false
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("MS Teams: request error: %s for URL %s", err, webhookURL)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logrus.Infof("MS Teams: successfully sent webhook to %s", webhookURL)
		return true
	}

	responseBodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logrus.Errorf("MS Teams: response body read error: %s (after non-success status %d from %s)", readErr, resp.StatusCode, webhookURL)
		return false
	}
	logrus.Errorf("MS Teams: failed to send webhook. Status: %d, URL: %s, Response: %s", resp.StatusCode, webhookURL, string(responseBodyBytes))
	return false
}

func generateTeamsTextMessage(req *request.JiraRequest, isComment bool) string {
	creatorName := "N/A"
	if req.Fields.Creator.DisplayName != "" {
		creatorName = req.Fields.Creator.DisplayName
	} else if req.Fields.Creator.Name != "" {
		creatorName = req.Fields.Creator.Name
	}

	assigneeName := "N/A"
	if req.Fields.Assignee.DisplayName != "" {
		assigneeName = req.Fields.Assignee.DisplayName
	} else if req.Fields.Assignee.Name != "" {
		assigneeName = req.Fields.Assignee.Name
	}

	requestTypeName := "N/A"
	webLink := "N/A"
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

	var titlePrefix string
	if isComment {
		titlePrefix = "📰 **New Comment Added**"
	} else {
		titlePrefix = "🎯 **New Issue/Update**"
	}
	return fmt.Sprintf(
		"%s\n\n"+
			"**Type:** %s \n\n"+
			"**Summary:** %s \n\n"+
			"**Issuer:** %s \n\n"+
			"**URL:** [%s](%s) \n\n"+
			"**Assignee:** %s",
		titlePrefix,
		requestTypeName,
		summary,
		creatorName,
		webLink, webLink,
		assigneeName,
	)
}
