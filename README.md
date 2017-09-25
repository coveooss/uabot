# uabot [![Build Status](https://travis-ci.org/coveo/uabot.svg?branch=master)](https://travis-ci.org/coveo/uabot) [![Go Report Card](https://goreportcard.com/badge/github.com/coveo/uabot)](https://goreportcard.com/report/github.com/coveo/uabot) [![license](https://img.shields.io/badge/license-Apache%20License%202.0-blue.svg)](https://github.com/coveo/uabot/blob/master/LICENSE)

Bot to send "intelligent" random usage analytics to simulate visits, queries and clicks on a site.
Works with a configuration file using json format.
Check the [/scenarios_examples](https://github.com/coveo/uabot/tree/master/scenarios_examples) folder for examples.

## How to use:

1. [Download executable](https://github.com/coveo/uabot/releases/latest).
2. Set Environment variables (refer to the section below).
3. Build your scenarios ([How to build scenarios](http://coveo.github.io/uabot/scenario.html)).
4. Execute the bot.

[Examples of scenarios](https://github.com/coveo/uabot/tree/master/scenarios_examples)

<hr/>

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

## Developper section
<hr/>

### To trigger a Docker rebuild, push with `latest` tag
```sh
1. Commit your changes
# We need to use the -f option here because tag latest already exists
2. git tag -f -a latest -m "Rebuild reason here"
# Push changes to branch + push changes to tag, you will need to supply credentials twice.
3. git push && git push -f --tags
```

### To release a newer version of the bot, with Travis automated builds
```sh
# Create a new tag with the version number to use.
1. git tag -a [v0.9.9] -m "Release comment here"
# Push tag.
2. git push --tags origin master
# It takes a little bit of time for Travis to generate the artefacts
3. Visit https://github.com/coveo/uabot/releases to view the releases to ensure your new release is visible.
```

[![forthebadge](http://forthebadge.com/images/badges/made-with-crayons.svg)](http://forthebadge.com)
