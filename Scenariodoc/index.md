# Documentation for the usage analytics bot configuration file format

This is the documentation index file for all the possible options of the data generation bot.
This file should be kept up to date with the changes, otherwise contact the author.

## [Documentation of how to build a scenario](Scenarios.md)

## General parameters

Parameters that are necessary are in **bold**, otherwise there is a default value in defaults\defauls.go
Parameters that are optional without a default value are in *italics*

Parameter | Type | Usage | Default
------------ | ------------- | ---------------- | -----------------
*orgName* | string | The name of the cloud org | (none)
searchendpoint | string | Endpoint where to direct the search queries | https://cloudplatform.coveo.com/rest/search/
analyticsendpoint | string | Endpoint where to direct the usage analytics events | https://usageanalytics.coveo.com/rest/v15/analytics/
**randomGoodQueries** | []string | The dataset of random queries (good ones) | ""
**randomBadQueries** | []string | The dataset of random queries (bad ones) | ""
[**scenarios**](Scenarios.md) | []Scenarios | The dataset of scenarios to execute | (none) See [documentation](Scenarios.md)
timeBetweenVisits | number | The time to wait between each visits (between 0 and X seconds) | 120 seconds
timeBetweenActions | number | The time to wait between each actions (between 0 and X seconds) | 3 seconds
*pipeline* | string | The name of the pipeline the queries will use | (none)
*defaultOriginLevel1* | string | The name of the originLevel1 param by default | (none)
partialMatch | boolean | Enable partial match on the queries | false
partialMatchKeywords | number | Number of words after which to enable partial match | (none)
partialMatchTreshold | string | Number of words considered in a partial match (string because you can send "50%") | (none)
allowAnonymousVisits | boolean | If you allow some of the visits to be anonymous | false
anonymousTreshold | number | Number between 0 and 1 of the % of anonymous visits | 0
globalfilter | string | A filter to be applied to all queries | ""
languages | []string | A list of random languages for the visits | (none)

### Change default datasets parameters

All the parameters in this section have a default dataset defined in the .\defaults\defaults.go file. But you can override them by setting some yourself in the config file.

*The actual emails of the user will be composed of a mix of firstname+"."+lastname+emailsuffix*

*In order to influence random values in the datasets, just copy the same value multiple time*

Parameter | Type | Usage
------------ | ------------- | ----------------
Emails | []string | Email suffixes for the users in the form of `@suffix.com`
FirstNames | []string | A list of possible first names
LastNames | []string | A list of possible last names
RandomIPs | []string | A list of possible IP
UserAgents | []string | A list of UserAgent strings
Languages | []string | A list of possible languages (TBD)
MobileUserAgents | []string | A list of UserAgents strings that are on mobile

### Example

```json
{
  "searchendpoint"          : "https://cloudplatform.coveo.com/rest/search/",
  "analyticsendpoint"       : "https://usageanalytics.coveo.com/rest/v15/analytics/",
  "partialMatch"            : true,
  "partialMatchKeywords"    : 4,
  "partialMatchTreshold"    : "50%",
  "defaultOriginLevel1"     : "origin",
  "timeBetweenVisits"       : 120,
  "timeBetweenActions"      : 3,
  "allowAnonymousVisits"	: true,
  "anonymousTreshold"       : 0.5,
  "orgName"     			: "orgname",
  "pipeline"                : "pipeline",
  "allowEntitlements"       : true,
  "randomCustomData"        : [ { "apiname": "nameofacustomdimension", "values" : [ "value 1", "value 2", "value 3" ] } ],
  "languages"               : ["en", "fr"],
  "globalfilter"            : "@uri",
  "randomGoodQueries"       : [ "First query", "Second query", "Third query", "etc..." ],
  "randomBadQueries"        : [ "First bad query", "Second bad query", "You can even use query syntax @source=Sharepoint", "etc..." ],
  "scenarios"               : []
}
```
