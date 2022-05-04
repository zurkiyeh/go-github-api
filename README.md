# go-github-cli
golang-github-api is a CLI utility that fetches pull-requests from your favorite reposiroty within a specified time period. This will develop over time to become a CLI utility to search for any type of issue on github. Finally, the results will be sent to a specific email address.

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
## Other requirements:
- Have docker installed
- In order to receive the email, you will have to create an App password. Instructions on how to do so for gmail can be [found here](https://support.google.com/accounts/answer/185833?hl=en).
- Github allows unauthenticated users up to 60 calls/hr. If you wish, you can authenticated which allows for more calls. This CLI supports Basic authenticated (authentication with header) only for now. To generate a token, refer to [github official docs here](https://github.com/settings/tokens/new). To read more about Github Basic Authentication, refer to [this section of the docs](https://docs.github.com/en/rest/overview/other-authentication-methods#basic-authentication).

# Usage
```
Fetch pull-requests from your favorite repo!

Usage:
  go-github-api fetch-pr [flags]

Flags:
  -c, --config string   Set configfile alternate location. Default is .config.yaml in root dir.
  -d, --debug           Set log level to DEBUG.
  -f, --from string     start time for search period
  -h, --help            help for fetch-pr
  -r, --repo string     Specify repo to be searched. Format: "Org/repo_name". Default will be charmbracelet/wish
  -t, --to string       End time for search period
  ```

## Notes about supported flags
- Flags -t/-f support YYYY/MM/DD date format. They do not support time in Hours/Minutess/Seconds for now. This was a design choice to improve user experience. Otherwise users will have to specify long time formats which can be tedious.
- The -t (\"to\") flag defaults to today's date if not specified and the -f (\"from\") defaults to a week from now. If you explicitly specify only one of the flags, the other won't be included in the query. In other words, default behavior applies only when neither of the flags are specified.
- Pass the -d flags to read more debug output
- You can override the default repo with your own with -r flag. Please make sure to pass Org-name/repo-name. An example would be passing this repo as such:  
> -r zurkiyeh/go-github-api

## Flag Example
This example will only show flag options. To learn more about how you can run the utility, refer to Run section below.  

To fetch all pull-requests from this repo in the last week:
> -r zurkiyeh/go-github-api  

To fetch all pull-requests from charmbracelet/wish between May 2017 and today:
> -f 2017-05-01
  
# Run go-github-cli
You can check out the repo and run locally as follows:  

In the root dir:  
> go run cmd/cli/main.go fetch-pr [FLAGS]  
## Running with Docker
First, in the root directory, build the docker image:
> docker build -t go-github-search-api:latest-local .

To run:
> docker run -it go-github-search-api:latest-local fetch-pr [FLAGS]

# Testing
The command package was not tested as the tests are provided by the cobra package

To test client_test.go, you must have the env var GOOD_TOKEN defined 
> export GOOD_ENV=PLACE_YOUR_TOKEN_HERE

To run the tests. Navigate to root directory, then:
> go test ./... -v
## Testing with Docker
In case you don not have the go set up on your machine and do not want to go through configuring it. A dev image is also available.   
  
Build the test (Dockerfile.dev) image:
> docker build -f Dockerfile.dev -t go-github-search-api:testing-local .

Then run the image with -it to map terminal to container. The default entrypoint here is specified to /bin/bash:
> docker run -it go-github-search-api:testing-local

Finally, execute the same command as you would in a regular environement:
> go test ./... -v
# Future Improvements
- Write output to Writer object which allows to forward to varies output sources
- Use context for http calls