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
	DisplayName      = "Service Desk"
	PLATFORM_SUBTEAM = "platform"
	NETWORK_SUBTEAM  = "network"
	RUNTIME_SUBTEAM  = "runtime"
)

type (
	Proxy struct {
		ElementConf config.Element
	}

	ElementBody struct {
		Text        string `json:"text"`
		DisplayName string `json:"displayName"`
	}
)

func (p *Proxy) ProxyToElementHandler(isComment bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		subteam := c.Param("team")
		req := &request.JiraRequest{}
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading body")
		}

		c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
		fmt.Printf("Request Body: %s\n", string(body))

		err = c.Bind(req)
		if err != nil {
			logrus.Errorf("failed to read request body: %s", err.Error())
			return c.NoContent(http.StatusBadRequest)
		}
		generatedElementText := generateElementText(req, isComment)
		logrus.Printf("team name: %s\n generatedElementText: %s", subteam, generatedElementText)
		switch subteam {
		case PLATFORM_SUBTEAM:
			logrus.Printf("using platform url %s", p.ElementConf.PlatformURL)
			if p.proxyRequest(generatedElementText, p.ElementConf.PlatformURL) {
				return c.NoContent(http.StatusOK)
			}
		case NETWORK_SUBTEAM:
			logrus.Printf("using network url %s", p.ElementConf.NetworkURL)
			if p.proxyRequest(generatedElementText, p.ElementConf.NetworkURL) {
				return c.NoContent(http.StatusOK)
			}
		case RUNTIME_SUBTEAM:
			logrus.Printf("using runtime url %s", p.ElementConf.RuntimeURL)
			if p.proxyRequest(generatedElementText, p.ElementConf.RuntimeURL) {
				return c.NoContent(http.StatusOK)
			}
		default:
			if p.proxyRequest(generatedElementText, p.ElementConf.URL) {
				return c.NoContent(http.StatusOK)
			}
		}

		return c.NoContent(http.StatusInternalServerError)
	}
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
	defer func() {
		_ = resp.Body.Close()
	}()

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

func generateElementText(req *request.JiraRequest, isComment bool) string {
	if isComment {
		return fmt.Sprintf(
			"📰\nNew Comment Added\nType: %s\nSummary: %s\nIssuer: %s\nURL: %s\nAssignee: %s\n",
			req.Fields.CustomField10003.RequestType.Name,
			req.Fields.Summary,
			req.Fields.Creator.DisplayName,
			req.Fields.CustomField10003.Links.Web,
			req.Fields.Assignee.DisplayName,
		)
	} else {
		return fmt.Sprintf(
			"🎯\nType: %s\nSummary: %s\nIssuer: %s\nURL: %s\nAssignee: %s",
			req.Fields.CustomField10003.RequestType.Name,
			req.Fields.Summary,
			req.Fields.Creator.DisplayName,
			req.Fields.CustomField10003.Links.Web,
			req.Fields.Assignee.DisplayName,
		)
	}
}
