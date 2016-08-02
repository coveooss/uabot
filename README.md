# uabot [![Build Status](https://travis-ci.org/coveo/uabot.svg?branch=master)](https://travis-ci.org/coveo/uabot) [![Go Report Card](https://goreportcard.com/badge/github.com/coveo/uabot)](https://goreportcard.com/report/github.com/coveo/uabot) [![license](https://img.shields.io/badge/license-Apache%20License%202.0-blue.svg)](https://github.com/coveo/uabot/blob/master/LICENSE)

Bot to send "intelligent" random usage analytics to simulate visits, queries and clicks on a site.
Works with a configuration file using json format.
Check the /Scenarios folder for examples.

## [Usage documentation](http://coveo.github.io/uabot/)
## [Code documentation](http://godoc.org/github.com/coveo/uabot/scenariolib)

## Supports
- [x] Search events
- [x] Click events
- [x] Custom events
- [x] Facet events
- [x] View events
- [x] TabChange events

## Environment variables
Needs 3 environment variables to function :

Variable | Usage
------------ | -------------
SEARCHTOKEN | API key to search
UATOKEN | API key to send events to UA
SCENARIOSURL | Url to the scenario JSON file to randomize
LOCAL | `true` for local (otherwise don't use)
GO15VENDOREXPERIMENT | Use this for go version < 1.6


#### On windows
```sh
set SEARCHTOKEN=value
set UATOKEN=value
set SCENARIOSURL=value
set LOCAL=true #if scenariosurl is a local path
set GO15VENDOREXPERIMENT=1 #if golang version < 1.6
go run main.go
```
#### On MAC
```sh
export SEARCHTOKEN = value
export UATOKEN = value
export SCENARIOSURL = value
export LOCAL=true #if scenariosurl is a local path
export GO15VENDOREXPERIMENT=1 #if golang version < 1.6
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
