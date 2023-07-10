package handler

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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

	elementBody struct {
		Text        string `json:"text"`
		DisplayName string `json:"displayName"`
	}
)

func (p *Proxy) ProxyToElement(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Errorf("failed to read request body: %s", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	logrus.Infof("jira request body: %s", string(body))

	if p.proxyRequest(body, p.ElementURL) {
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func (p *Proxy) proxyRequest(body []byte, url string) bool {
	body, err := json.Marshal(elementBody{
		Text:        string(body),
		DisplayName: DisplayName,
	})
	if err != nil {
		logrus.Errorf("proxy request to element error: %s", err)
		return false
	}

	proxyReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logrus.Errorf("proxy request to element error: %s", err)
		return false
	}

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
	logrus.Errorf("proxy request to element error: %s", responseBody)
	return false
}
