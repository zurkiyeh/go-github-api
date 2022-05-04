# go-github-cli
Golang github api cli utility that fetches latest pull-requests from your favorite reposiroty within a specified time period.

# Usage
```
Fetch pull-requests from your favorite repo!

Usage:
  go-github-api fetch-pr [flags]

Flags:
  -c, --config string   Set configfile alternate location. Default is .config.yaml in this dir.
  -d, --debug           Set log level to DEBUG.
  -f, --from string     start time for search period
  -h, --help            help for fetch-pr
  -r, --repo string     Specify repo to be searched. Format: "Org/repo_name". Default will be charmbracelet/wish
  -t, --to string       End time for search period
  ```

# Requirements
- cli will look for a .config.yaml file and load it
- if .config.yaml file doesnt exist or an email&password havent been specified, the output will be stdout instead
- No user input validation on time 


# Pre-requisites
## Config file
You must initialize a config file in the root dir called ".config.yaml". Location of config file can be overridden with -c flag. Take note of the name as it is added to .gitignore and should not be committed in the repo. The config file should look as follows:
```
---
Credentials:
  EmailAddress: your.email@org.com
  EmailPassword: eM@!L_PaSsw0rd # Refer to section "Generating Email Token" to see how you can generate a token in place of actual password  
  EmailRecepient: your.recepient@org.com #to field
  EmailSmtpHost : smtp.org.com #smtp info should be "smtp.gmail.com" for gmail
  EmailSmtpPort : 587 # Default is 587 
PersonalToken: GITHUB_TOKEN # Refer to section "Access Tokens" to learn how to generate the Github token
```
## other requirements:
- have docker installed
# How to use:
https://api.github.com/search/issues?q=is:pull-request+created:2006-01-02..2022-05-01+repo:gin-gonic/gin

# Set up Email to send emails:
https://support.google.com/accounts/answer/185833?hl=en


# Access Tokens
https://github.com/settings/tokens
Generate Token:
https://github.com/settings/tokens/new


only supports days
default behavior set to search for a week from now


# Run with Docker
first build in the main dir:
> docker build -t go-github-search-api:latest-local .

To run:
> docker run -it go-github-search-api:latest-local [FLAGS]

# Testing
The command package was not tested as the tests are provided by the cobra package

To test client_test.go, you must have the env var GOOD_TOKEN defined 
> export GOOD_ENV=PLACE_YOUR_TOKEN_HERE

To run the tests:
