package command

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/zurkiyeh/go-github-cli/transport"
)

var (
	default_repo = "charmbracelet/wish"
)

const (
	time_layout = "2006-01-02T15:04:05"
)

// Fetch pull-requests config
type configFetchPullRequest struct {
	ConfigFile string
	Repo       string
	To         string
	From       string
	Debug      bool
}

func newCommandFetchPRs() *cobra.Command {
	var configFile configFetchPullRequest
	cmd := &cobra.Command{
		Use:          "fetch-pr",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := initConfig(configFile.ConfigFile)
			if err != nil {
				return err
			}
			logger, err := initLogger(configFile.Debug)
			if err != nil {
				return err
			}
			// logger.Info("Your email address:", config.Credentials.EmailAddress)
			// logger.Info("Your email address:", config.Credentials.EmailPassword)

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				sigChan := getSigChan()
				for {
					select {
					case <-sigChan:
						cancel()
						return
					}
				}
			}()

			c := transport.NewClient(logger,
				&config.PersonalToken)
			if err != nil {
				c.Logger.Errorf("an error occured while parsing user input: %s", err)
				return err
			}
			query, err := buildQuery(configFile)
			if err != nil {
				c.Logger.Error("error parsing user input")
				return err
			}
			c.Logger.Debug("Query is: ", query)
			req, _ := c.NewRequest("GET", "search/issues", query, config.PersonalToken)

			c.Logger.Info("Request: ", req.URL)
			err = c.Do(ctx, req)

			return err
		},
	}
	return setupflags(cmd, &configFile)
}
