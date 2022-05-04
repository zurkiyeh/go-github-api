package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Struct to parse the config file
type Config struct {
	Credentials struct {
		EmailAddress   string
		EmailPassword  string
		EmailRecepient string
		EmailSmtpHost  string
		EmailSmtpPort  string
	}
	PersonalToken string
}

var (
	defaultOutput = os.Stdout
)

// Execute executes the root command.
func Execute() {
	rootCmd := newCommandRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(defaultOutput, err)
	}
}

// Init the root command info
func newCommandRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "go-github-api",
		Short:        "CLI utility to interact with Github search API",
		Version:      "1.0.0",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		newCommandFetchPRs(),
	)
	return cmd
}

// set up flags for fetch-pr command
func setupflags(cmd *cobra.Command, c *configFetchPullRequest) *cobra.Command {
	cmd.Flags().StringVarP(&c.ConfigFile, "config", "c", c.ConfigFile, "Set configfile alternate location. Default is .config.yaml in this dir.")
	cmd.Flags().StringVarP(&c.Repo, "repo", "r", c.ConfigFile, "Specify repo to be searched. Format: \"Org/repo_name\". Default will be charmbracelet/wish")
	cmd.Flags().StringVarP(&c.To, "to", "t", c.ConfigFile, "End time for search period. Format YYYY/MM/DD")
	cmd.Flags().StringVarP(&c.From, "from", "f", c.ConfigFile, "start time for search period. Format YYYY/MM/DD")

	cmd.Flags().BoolVarP(&c.Debug, "debug", "d", c.Debug, "Set log level to DEBUG.")
	return cmd
}
