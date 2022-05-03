package command

import (
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		Use:  "fetch-pr",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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
	cmd.Flags().StringVarP(&c.Repo, "repo", "r", c.ConfigFile, "Specify repo to be searched. Format: \"Org/repo_name\". Default will be charmbracelet/wish")
	cmd.Flags().StringVarP(&c.To, "to", "t", c.ConfigFile, "End time for search period")
	cmd.Flags().StringVarP(&c.From, "from", "f", c.ConfigFile, "start time for search period")

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

func buildQuery(cfg configFetchPullRequest) (string, error) {
	var to, from string
	var repo string
	var err error
	query := "?q=is:pull-request"

	if cfg.From != "" {
		time, err := time.Parse(time_layout, cfg.From)
		if err != nil {
			return "", err
		}
		from, _ = extractDate(time)
		fmt.Printf("Overriding \"From\" time: %s\n", from)
	} else {
		// set default values for both to and from times
		from, _ = extractDate(time.Now().AddDate(0, 0, -7))
		if err != nil {
			return "", err
		}
		fmt.Printf("Default \"From\" time: %s\n", from)
		to, _ = extractDate(time.Now())

	}

	if cfg.To != "" {
		time, err := time.Parse(time_layout, cfg.To)
		if err != nil {
			return "", err
		}
		to, _ = extractDate(time)
		fmt.Printf("Overriding to time: %s\n", to)
	} else {
		// set default values for both to and from times
		if err != nil {
			return "", err
		}
		fmt.Printf("Default to time: %s\n", to)
		to, _ = extractDate(time.Now())
	}

	if (cfg.To != "" && cfg.From != "") || (cfg.From == "" && cfg.To == "") {
		query += fmt.Sprintf("+created:%s..%s", from, to)
	} else if cfg.From != "" {
		query += fmt.Sprintf("+created:>%s", from)
	} else {
		query += fmt.Sprintf("+created:<%s", to)
	}

	if cfg.Repo == "" {
		repo = default_repo
		fmt.Printf("Default repo : %s\n", repo)
	} else {
		repo = cfg.Repo
		fmt.Printf("Overriding repo : %s\n", repo)
	}
	query += fmt.Sprintf("+repo:%s", repo)
	return query, nil
}

func extractDate(date time.Time) (string, error) {
	dateStr := "%02d-%02d-%02d"
	return fmt.Sprintf(dateStr, date.Year(), date.Month(), date.Day()), nil
}
