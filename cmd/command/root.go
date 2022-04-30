package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
