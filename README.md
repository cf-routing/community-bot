# Community Bot
A script that helps the CF Networking team stay updated on Github Community issues and PRs.

## Getting Started
To get started, make sure you clone this repo into your `GOPATH`:
```
cd $GOPATH
mkdir -p src/github.com/cf-routing
cd src/github.com/cf-routing
git clone git@github.com:cf-routing/community-bot.git
cd community-bot
go get github.com/golang/dep/cmd/dep
dep ensure
```

## Usage
The bot prints out all the repos relevant to the routing team, along with the
issues and PRs associated with each repo. The issues/PRs are sorted by Least
Recently Updated.

Run the following script to get all open issues:
```
scripts/get_open_issues.sh
```

You can also provide `SINCE=N` where N is the number of days back you would like
to see issues from. If not provided, the program will return _all_ open issues and
PRs.

Alternatively, run it manually:

1. Get the github access token from LastPass:
   ```
   github_access_token=$(lpass show -j "Github - Routing CI Bot" | jq -r ".[0].note" | awk '{print $NF}')
   ```

2. Run main:
   ```
   GITHUB_ACCESS_TOKEN=${github_access_token} go run main.go
   ```
