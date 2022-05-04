package command

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Default values in case flags are not passed
var (
	default_repo = "charmbracelet/wish"
)

const (
	time_layout = "2006-01-02"
)

// Initializing config file
func initConfig(cfgFile string) (*Config, error) {
	viperInstance := viper.New()
	if cfgFile != "" {
		viperInstance.SetConfigFile(cfgFile)
	} else {
		// defaults to .config.yaml in root dir
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		hd, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		viperInstance.AddConfigPath(wd)
		viperInstance.AddConfigPath(hd)
		viperInstance.SetConfigName(".config")
		viperInstance.SetConfigType("yaml")
	}

	viperInstance.AutomaticEnv()

	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	if err := viperInstance.UnmarshalExact(&config); err != nil {
		return nil, err
	}
	return &config, nil

}

// Initalize logger information based on the value of -d flag passed at CLI
func initLogger(setDebug bool) (*logrus.Logger, error) {
	logger := logrus.New()

	// Info level by default
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

// Builds a query to attach to outgoing request
func buildQuery(initQuery string, cfg configFetchPullRequest, logger *logrus.Logger) (string, error) {
	var to, from string
	var repo string
	var err error
	query := initQuery

	// Check that -f has not been overridden from the CLI. If so, try to parse the new passed value
	if cfg.From != "" {
		time, err := time.Parse(time_layout, cfg.From)
		if err != nil {
			return "", err
		}
		from, _ = extractDate(time)
		logger.Info("Overriding \"From\" time: ", from)
	} else {
		// if not set, default to a week from now (Default behavior)
		from, _ = extractDate(time.Now().AddDate(0, 0, -7))
		if err != nil {
			return "", err
		}
		logger.Info("Default \"From\" time: ", from)
		from, _ = extractDate(time.Now())
	}

	// Check that -t has not been overridden from the CLI. If so, try to parse the new passed value
	if cfg.To != "" {
		time, err := time.Parse(time_layout, cfg.To)
		if err != nil {
			return "", err
		}
		to, _ = extractDate(time)
		if err != nil {
			return "", err
		}
		logger.Info("Overriding to time: ", to)
	} else {
		// if not set, default to now (Default behavior)
		if err != nil {
			return "", err
		}
		to, err = extractDate(time.Now())
		if err != nil {
			return "", err
		}
		logger.Info("Default \"to\" time: ", to)
	}

	// Add query string based on flags passed to adher to github specifications
	if (cfg.To != "" && cfg.From != "") || (cfg.From == "" && cfg.To == "") {
		query += fmt.Sprintf("+created:%s..%s", from, to)
	} else if cfg.From != "" {
		query += fmt.Sprintf("+created:>%s", from)
	} else {
		query += fmt.Sprintf("+created:<%s", to)
	}

	// Check if repo has been overridden
	if cfg.Repo == "" {
		repo = default_repo
		logger.Info("Default repo : ", default_repo)
	} else {
		repo = cfg.Repo
		logger.Info("Overriding repo :", repo)
	}
	query += fmt.Sprintf("+repo:%s", repo)
	return query, nil
}

// Parse date to adhre to YYYY/MM/DD format
func extractDate(date time.Time) (string, error) {
	dateStr := "%02d-%02d-%02d"
	return fmt.Sprintf(dateStr, date.Year(), date.Month(), date.Day()), nil
}
