package cmd

import (
	"github.com/snapp-incubator/jira-element-proxy/internal/config"
	"github.com/snapp-incubator/jira-element-proxy/internal/webhook-proxy/cmd/api"

	"github.com/spf13/cobra"
)

// NewRootCommand creates a new webhook-proxy root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "webhook-proxy",
	}

	cfg := config.New()

	api.Register(root, cfg)

	return root
}
