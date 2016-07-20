# Autobot 

Application wrapper to the exisiting [uabot](https://github.com/coveo/uabot) to send "intelligent" usage analytics to simulate user interaction with an organization. Autobot grant an endpoint service to POST an autobot query.

## Installation

1. Install [Docker-engine](https://docs.docker.com/engine/installation/)
2. Use the scripts to build and run the docker image

## Usage

Once installed, you can use the following API to use Autobot :

To post a task to the robot
```
POST : [HOST]:8080/start
HEADER : {Content-Type : application/json}
BODY : {
[REQUIRED] "searchEndpoint" : YOUR-SEARCH-ENDPOINT, 
[REQUIRED] "searchToken" : YOUR-SEARCH-TOKEN, 
[REQUIRED] "analyticsEndpoint" : YOUR-ANALYTICS-ENDPOINT, 
[REQUIRED] "analyticsToken" : YOUR-ANALYTICS-TOKEN, 
[REQUIRED] "timeToLive" : LIFETIME-OF-THE-AUTOBOT, 
[REQUIRED] "originLevels" : {ORIGIN-LEVEL1 : [LIST-OF-ORIGIN-LEVEL-2]}, 
[OPTIONAL] "avgNumberWordsPerQuery" : AVERAGE-NUMBER-OF-WORDS-PER-QUERY (default=1), 
[OPTIONAL] "fetchQueryNumber" : NUMBER-OF-RESULT-IN-SEARCH-RESPONSE (default=1000), 
[OPTIONAL] "explorationRatio" : INDEX-EXPLORATION-RATIO (default=0.01), 
[OPTIONAL] "numberOfQueryPerLanguage" : MAX-NUMBER-OF-QUERY-PER-LANGUAGE (default=10), 
[OPTIONAL] "fields" : FIELDS-TO-EXPLORE-EQUALLY (default=["@syssource"]), 
}
```

To stop a task prematurely
```
POST : [HOST]:8080/stop/{workerid}
workerid    :   The id provided by the server
```

To get information about running tasks
```
GET : [HOST]:8080/info
```
