package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Credentials struct {
		EmailAddress  string
		EmailPassword string
	}
	EncryptionKey string
}

var (
	defaultOutput = os.Stdout

	// timeFormats = []string{
	// 	"2006-01-02",
	// 	"2006-01-02T15:04",
	// 	"2006-01-02T15:04:05",
	// }
)

// Execute executes the root command.
func Execute() {
	rootCmd := newCommandRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(defaultOutput, err)
	}
}

func newCommandRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "go-github-api",
		Short:   "Fetch pull-requests from your favorite repo",
		Version: "1.0.0",
	}
	cmd.AddCommand(
		newCommandFetchPRs(),
	)
	return cmd
}

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
