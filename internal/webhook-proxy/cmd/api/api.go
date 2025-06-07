package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/snapp-incubator/jira-element-proxy/internal/webhook-proxy/handler"

	"github.com/snapp-incubator/jira-element-proxy/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	app := echo.New()

	proxyHandler := handler.Proxy{MSTeamsConf: cfg.MSTeams}

	logrus.Println("API has been started (MS Teams mode) :D")

	app.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	app.POST("/:team", proxyHandler.ProxyToMSTeamsHandler(false))
	app.POST("/comment/:team", proxyHandler.ProxyToMSTeamsHandler(true))
	app.POST("/", proxyHandler.ProxyToMSTeamsHandler(false))
	app.POST("/comment", proxyHandler.ProxyToMSTeamsHandler(true))

	if err := app.Start(fmt.Sprintf(":%d", cfg.API.Port)); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("echo initiation failed: %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "api",
			Short: "Run API to serve the requests (MS Teams mode)",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
