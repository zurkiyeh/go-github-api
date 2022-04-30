package command

import (
	"fmt"

	"github.com/spf13/cobra"
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
			fmt.Printf("Your email address: %s \n", config.Credentials.EmailAddress)
			fmt.Printf("Your email password: %s\n", config.Credentials.EmailPassword)
			return nil
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
