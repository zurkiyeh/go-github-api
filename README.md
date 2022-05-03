# go-github-cli
Golang github api cli utility that fetches latest pull-requests within a time period

# Requirements
- cli will look for a .config.yaml file and load it
- if .config.yaml file doesnt exist or an email&password havent been specified, the output will be stdout instead
- No user input validation on time 



# How to use:
https://api.github.com/search/issues?q=is:pull-request+created:2006-01-02..2022-05-01+repo:gin-gonic/gin

# Set up Email to send emails:
https://support.google.com/accounts/answer/185833?hl=en


# Access Tokens
https://github.com/settings/tokens
Generate Token:
https://github.com/settings/tokens/new