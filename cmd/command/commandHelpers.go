package command

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig(cfgFile string) (*Config, error) {
	viperInstance := viper.New()
	if cfgFile != "" {
		viperInstance.SetConfigFile(cfgFile)
	} else {
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
