package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/snapp-incubator/jira-element-proxy/internal/webhook-proxy/request"
	"io"
	"net/http"
)

const (
	DisplayName = "Service Desk"
)

type (
	Proxy struct {
		ElementURL string
	}

	ElementBody struct {
		Text        string `json:"text"`
		DisplayName string `json:"displayName"`
	}
)

func (p *Proxy) ProxyToElement(c echo.Context) error {
	req := &request.Jira{}

	err := c.Bind(req)
	if err != nil {
		logrus.Errorf("failed to read request body: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if p.proxyRequest(generateElementText(req), p.ElementURL) {
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func (p *Proxy) proxyRequest(txt string, url string) bool {
	body, err := json.Marshal(ElementBody{
		Text:        txt,
		DisplayName: DisplayName,
	})
	if err != nil {
		logrus.Errorf("marshal request body error: %s", err)
		return false
	}

	proxyReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logrus.Errorf("create proxy request error: %s", err)
		return false
	}

	proxyReq.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		logrus.Errorf("proxy request to element error: %s", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("element response body read error: %s", err)
		return false
	}

	logrus.Errorf("element response body read error: %s", responseBody)
	return false
}

func generateElementText(req *request.Jira) string {
	return fmt.Sprintf(
		"Type: %s\nSummary: %s\nIssuer: %s\nURL: %s",
		req.Issue.Fields.CustomField11401.RequestType.Name, req.Issue.Fields.Summary, req.User.Name,
		req.Issue.Fields.CustomField11401.Links.Web)
}
