package api

import (
	"errors"
	"fmt"
	"github.com/snapp-incubator/jira-element-proxy/internal/webhook-proxy/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/snapp-incubator/jira-element-proxy/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	app := echo.New()

	if err := app.Start(fmt.Sprintf(":%d", cfg.API.Port)); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("echo initiation failed: %s", err)
	}

	proxyHandler := handler.Proxy{ElementURL: cfg.Element.URL}

	logrus.Println("API has been started :D")

	app.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	app.GET("/element", proxyHandler.ProxyToElement)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// Register API command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "api",
			Short: "Run API to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
