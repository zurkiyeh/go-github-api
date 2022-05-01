package command

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zurkiyeh/go-github-cli/transport"
)

// Fetch pull-requests config
type configFetchPullRequest struct {
	ConfigFile string
	Debug      bool
}

func newCommandFetchPRs() *cobra.Command {
	fmt.Println("Test")
	var configFile configFetchPullRequest
	cmd := &cobra.Command{
		Use:  "fetch-pr",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Print from within the command\n")
			config, err := initConfig(configFile.ConfigFile)
			if err != nil {
				return err
			}
			logger, err := initLogger(configFile.Debug)
			if err != nil {
				return err
			}
			logger.Info("Your email address:", config.Credentials.EmailAddress)
			logger.Info("Your email address:", config.Credentials.EmailPassword)

			c := transport.NewClient(logger)
			// queryParams.Set("repo", "charmbracelet/wish")
			// "?q=repo:+is:pull-request+created:>2022-04-28"
			params := url.Values{
				"repo": {"charmbracelet/wish"},
			}
			req, _ := c.NewRequest("GET", "search/issues", params, nil)
			c.Logger.Info("Request: ", req.URL)
			err = c.Do(context.Background(), req)

			return err
		},
	}
	return setupflags(cmd, &configFile)
}

// set up flags for fetch-pr command
func setupflags(cmd *cobra.Command, c *configFetchPullRequest) *cobra.Command {
	cmd.Flags().StringVarP(&c.ConfigFile, "config", "c", c.ConfigFile, "Set configfile alternate location. Default is .config.yaml in this dir.")

	cmd.Flags().BoolVarP(&c.Debug, "debug", "d", c.Debug, "Set log level to DEBUG.")
	return cmd
}

func initLogger(setDebug bool) (*logrus.Logger, error) {
	logger := logrus.New()

	logger.SetLevel(logrus.InfoLevel)
	if setDebug {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		DisableSorting:         true,
	})

	logger.SetOutput(os.Stdout)
	return logger, nil
}
