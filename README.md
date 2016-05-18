[![Build Status](https://travis-ci.org/erocheleau/uabot.svg?branch=master)](https://travis-ci.org/erocheleau/uabot)

# uabot
Bot to generate data to UA

Needs 3 environment variables to function :

SEARCHTOKEN  => API key to search

UATOKEN      => API key to send events to UA

SCENARIOSURL => Url to the scenario JSON file to randomize

GO15VENDOREXPERIMENT=1

(Only if running golang < 1.6) To use the vendor folder for the lib for analytics/search api

## On windows
```sh
set SEARCHTOKEN=value
set UATOKEN=value
set SCENARIOSURL=value
# To use a local scenario json file
set LOCAL=true
# Only if running golang < v1.6
set GO15VENDOREXPERIMENT=1
go run main.go
```
## On MAC
```sh
export SEARCHTOKEN = value
export UATOKEN = value
export SCENARIOSURL = value
# To use a local scenario json file
export LOCAL=true
# Only if running golang < v1.6
export GO15VENDOREXPERIMENT=1
go build
go run main.go
```

## To trigger a Docker rebuild, push with `latest`tag
```sh
1. Commit your changes
# We need to use the -f option here because tag latest already exists
2. git tag -f -a latest -m "Rebuild reason here"
# Push changes to branch + push changes to tag, you will need to supply credentials twice.
3. git push && git push -f --tags
```

[![forthebadge](http://forthebadge.com/images/badges/made-with-crayons.svg)](http://forthebadge.com)
