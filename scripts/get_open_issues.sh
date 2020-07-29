#!/usr/bin/env bash

lpass ls > /dev/null
if [ $(echo $?) -ne 0 ]
then
  echo "You are not logged in to the LastPass CLI. Please log in and try again."
  exit 1
fi

token=$(lpass show -j "Github - Routing CI Bot" | jq -r ".[0].note" | awk '{print $NF}')
GITHUB_ACCESS_TOKEN="${token}" go run main.go
