# uabot [![Build Status](https://travis-ci.org/coveo/uabot.svg?branch=master)](https://travis-ci.org/coveo/uabot) [![Go Report Card](https://goreportcard.com/badge/github.com/coveo/uabot)](https://goreportcard.com/report/github.com/coveo/uabot) [![license](https://img.shields.io/badge/license-Apache%20License%202.0-blue.svg)](https://github.com/coveo/uabot/blob/master/LICENSE)

Bot to send "intelligent" random usage analytics to simulate visits, queries and clicks on a site.
Works with a configuration file using json format.
Check the [/scenarios_examples](https://github.com/coveo/uabot/tree/master/scenarios_examples) folder for examples.

## How to use:

1. [Download executable](https://github.com/coveo/uabot/releases/latest).
2. Set Environment variables (refer to the section below).
3. Build your scenarios ([How to build scenarios](#configuration-files)).
4. Execute the bot.

### Tracing

You can use the argument `-trace` to get more logs when debugging your scenarios.

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
LOCAL | `true` for local (default false)


### On windows

```sh
set SEARCHTOKEN=value
set UATOKEN=value
set SCENARIOSURL=value
set LOCAL=true #if scenariosurl is a local path
```

### On MAC

```sh
export SEARCHTOKEN = value
export UATOKEN = value
export SCENARIOSURL = value
export LOCAL=true #if scenariosurl is a local path
```

## Configuration Files

### Documentation

- [How to build scenarios](http://coveo.github.io/uabot/scenario.html)
- [Examples of scenarios](https://github.com/coveo/uabot/tree/master/scenarios_examples)

### Visual Studio Code Snippets

The [uabot.code-snippets](/uabot.code-snippets) file is a snippets library for the Visual Studio Code editor. It contains snippets for configurations, scenarios, all supported events, and custom events for e-commerce websites.

#### Configuration Snippets

Prefix | Description
-- | --
uabotconfig | Coveo UABot Complete Configuration
uabotconfig | Coveo UABot Simple Configuration
uabotrandomcustomdata | Coveo UABot Random Custom Data
uabotscenario | Coveo UABot Complete Scenario
uabotscenario | Coveo UABot Simple Scenario

#### Event Snippets

Prefix | Description
-- | --
uaboteventsearch | Coveo UABot Complete Search Event
uaboteventcustomdata | Coveo UABot Event Custom Data
uaboteventsearch | Coveo UABot GoodQuery Search Event
uaboteventsearch | Coveo UABot BadQuery Search Event
uaboteventsearch | Coveo UABot Custom Query Search Event
uaboteventclick | Coveo UABot Complete Click Event
uaboteventclick | Coveo UABot Precise Click Event
uaboteventclick | Coveo UABot Random Click Event
uaboteventsearchandclick | Coveo UABot Complete SearchAndClick Event
uaboteventsearchandclick | Coveo UABot Simple SearchAndClick Event
uaboteventtabchange | Coveo UABot TabChange Event
uaboteventfacetchange | Coveo UABot Complete FacetChange Event
uaboteventfacetchange | Coveo UABot Simple FacetChange Event
uaboteventsetorigin | Coveo UABot SetOrigin Event
uaboteventview | Coveo UABot Complete View Event
uaboteventview | Coveo UABot Precise View Event
uaboteventview | Coveo UABot Random View Event
uaboteventsetreferrer | Coveo UABot SetReferrer Event
uaboteventcustom | Coveo UABot Custom Event

#### Custom E-Commerce Event Snippets

Prefix | Description
-- | --
uaboteventcustomcommercedetailview | Coveo UABot Custom Commerce detailView Event
uaboteventcustomcommerceaddtocart | Coveo UABot Custom Commerce addToCart Event
uaboteventcustomcommerceremovefromcart | Coveo UABot Custom Commerce removeFromCart Event
uaboteventcustomcommerceaddpurchase | Coveo UABot Custom Commerce addPurchase Event
uaboteventcustomcommerceremovepurchase | Coveo UABot Custom Commerce removePurchase Event
uaboteventcustomcommerceaddrating | Coveo UABot Custom Commerce addRating Event
uaboteventcustomcommerceremoverating | Coveo UABot Custom Commerce removeRating Event
uaboteventcustomcommerceaddbookmark | Coveo UABot Custom Commerce addBookmark Event
uaboteventcustomcommerceremovebookmark | Coveo UABot Custom Commerce removeBookmark Event
uaboteventcustomcommerceaddcompare | Coveo UABot Custom Commerce addCompare Event
uaboteventcustomcommerceremovecompare | Coveo UABot Custom Commerce removeCompare Event

#### How to Install

1. Download the [uabot.code-snippets](/uabot.code-snippets) file.
2. Place the file in your Visual Studio Code snippets folder
    1. Windows: `%APPDATA%\Code\User\snippets`
    2. Mac: `$HOME/Library/Application Support/Code/User/snippets`
    3. Linux: `$HOME/.config/Code/User/snippets`
3. Open Visual Studio Code.

#### How to Use

In any plain text or JSON file:

1. Type a snippet prefix. You can also type a partial prefix.
2. Use the up/down arrow keys to choose the desired snippet.
3. Hit the `[Tab]` or `[Return]` key to trigger the snippet completion.
4. Fill in the blanks by replacing the help text. Use the `[Tab]` key to switch between parameters.
5. Repeat and enjoy!

## Developper section

<hr/>

### To build an executable

```sh
# Install the dependencies
1. go get
# Build an executable -o sets the output name
2. go build -o myexecutable
3. Run the executable
```

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
