package command

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/spf13/cobra"
	"github.com/zurkiyeh/go-github-cli/html"
	"github.com/zurkiyeh/go-github-cli/transport"
)

// Fetch pull-requests config
type configFetchPullRequest struct {
	ConfigFile string
	Repo       string
	To         string
	From       string
	Debug      bool
}

var (
	auth smtp.Auth
)

func newCommandFetchPRs() *cobra.Command {
	var configFile configFetchPullRequest
	cmd := &cobra.Command{
		Use:          "fetch-pr",
		Short:        "Fetch pull-requests from your favorite repo!",
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
			// Init client to start request
			c := transport.NewClient(logger,
				&config.PersonalToken)
			if err != nil {
				c.Logger.Errorf("an error occured while parsing user input: %s", err)
				return err
			}
			// Build query based on input of user
			query, err := buildQuery("?q=is:pull-request", configFile, c.Logger)
			if err != nil {
				c.Logger.Error("error parsing user input")
				return err
			}
			c.Logger.Debug("Query is: ", query)

			// Start a new request and execute
			req, _ := c.NewRequest("GET", "search/issues", query, config.PersonalToken)
			c.Logger.Debug("Request: ", req.URL)
			respJSON, err := c.Do(context.Background(), req)
			c.Logger.Debug("Request has returned ", len(respJSON.Items), " items")

			for i, pr := range respJSON.Items {
				c.Logger.Debug(" Pull-request #", i, " has title: ", pr.Title)
			}

			// Send to email
			fromEmailAddr := config.Credentials.EmailAddress
			fromEmailToken := config.Credentials.EmailPassword
			toEmailAddr := config.Credentials.EmailRecepient
			emailSmtpHost := config.Credentials.EmailSmtpHost
			emailSmtpPort := config.Credentials.EmailSmtpPort

			auth = smtp.PlainAuth("", fromEmailAddr, fromEmailToken, emailSmtpHost)

			r := html.NewRequest([]string{config.Credentials.EmailRecepient},
				"Your Search Results",
				"Hello from go-github-cli!")
			if err := r.ParseTemplate("./html/template.html", respJSON); err == nil {
				ok, _ := r.SendEmail(auth, toEmailAddr, emailSmtpHost, emailSmtpPort)
				c.Logger.Debug("Sent email:", ok)
				c.Logger.Info(fmt.Sprintf("search completed and results have been sent to: %s", config.Credentials.EmailRecepient))
			} else {
				c.Logger.Debug(err)
			}
			return err
		},
	}
	return setupflags(cmd, &configFile)
}
