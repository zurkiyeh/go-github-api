package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Fetch pull-requests config
type configFetchPullRequest struct {
	ConfigFile string

	LogFile   string
	StateFile string

	Output string
	Indent bool

	Debug       bool
	JSONLogging bool
}

func newCommandFetchPRs() *cobra.Command {
	fmt.Println("Test")
	var config configFetchPullRequest
	cmd := &cobra.Command{
		Use:  "fetch-pr",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Print from within the command\n")
			return nil
		},
	}
	return setupflags(cmd, &config)
}

func setupflags(cmd *cobra.Command, c *configFetchPullRequest) *cobra.Command {

	return cmd
}
