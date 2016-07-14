# Autobot 

Application wrapper to the exisiting [uabot](https://github.com/erocheleau/uabot) to send "intelligent" usage analytics to simulate user interaction with an organization. Autobot grant an endpoint service to POST an autobot query.  Autobot is split in two phases, the crawling phase and the event phase.

# Phases
### Crawling Phase

In the crawling phase, autobot will do some search on the index and will find words on titles and on different fields of the results.  These words are split into languages and they are processed to make valid queries.  The words may or may not be representative of the index, but they surely occur quite often if they appear in the query expressions.  In the end of that phase, a configuration file will be available for the uabot to proceed to the event pahse.

### Event Phase

In the event phase, Uabot takes the lead and apply the provided scenarios from the crawling phase.  Scenarios are duplicated for every languages found in the index provided, but they are modulated by their document occurences.  Let say you have 100 documents in English and 10 documents in French, then statistically speaking you can expect to have 10 times more usage analytics in english than in french. 
There are two main scenarios:

Search with a random query in the provided language
Click on one of the firsts documents with probability 40%
Search with a new random query in the same language 
Click on one of the firsts documents with probability 80%

Search with a random query in the provided language
Apply 20 times : View on one of the documents provided by the search event

## Installation

1. Install [Docker-engine](https://docs.docker.com/engine/installation/)
2. Use the scripts to build and run the docker image

## Usage

Once installed, you can use the following API to use Autobot :

To post a task to the robot

| Arugments | Usage |
| :----------:|:-------|
|SearchEndPoint| Defines the search endpoint on which you want to get your search results |
| SearchToken | Defines the token to use to identify your access on the search endpoint |
| AnalyticsEndpoint | Defines the analytics endpoint on which you want to post your events |
| AnalyticsToken | Defines the token to use to identify tour access on the analytics endpoint |
| Timeout | Defines the time life in minutes of your uabot until he stops |
| OriginLevels | Defines a map of the origin level 1 associated with a list of the origin level 2 |
| AvgNumberWordsPerQuery | Defines the average number of words to find per query, following an exponential distribution  (Default=1)|
| FertchQueryNumber | Defines the amount of results a single search will return (Default=1000) |
| ExplorationRatio | Defines the number of query to apply on each fields => ExplorationRatio * NumberOfDocumentsInIndex / FetchQueryNumber (Default=0.01)|
| NumberOfQueryPerLanguage | Defines the number of query expression the crawling phase will give (Default=10) |
| Fields | Defines the fields on which you want to explore each values equally (Default=["@syssource"] |
```
POST : [HOST]:8080/start
HEADER : {Content-Type : application/json}
BODY : {
[REQUIRED] "searchEndpoint" : YOUR-SEARCH-ENDPOINT, 
[REQUIRED] "searchToken" : YOUR-SEARCH-TOKEN, 
[REQUIRED] "analyticsEndpoint" : YOUR-ANALYTICS-ENDPOINT, 
[REQUIRED] "analyticsToken" : YOUR-ANALYTICS-TOKEN, 
[REQUIRED] "timeout" : LIFETIME-OF-THE-AUTOBOT, 
[REQUIRED] "originLevels" : {ORIGIN-LEVEL1 : [LIST-OF-ORIGIN-LEVEL-2]}, 
[OPTIONAL] "avgNumberWordsPerQuery" : AVERAGE-NUMBER-OF-WORDS-PER-QUERY (default=1), 
[OPTIONAL] "fetchQueryNumber" : NUMBER-OF-RESULT-IN-SEARCH-RESPONSE (default=1000), 
[OPTIONAL] "explorationRatio" : INDEX-EXPLORATION-RATIO (default=0.01), 
[OPTIONAL] "numberOfQueryPerLanguage" : NUMBER-OF-QUERY-PER-LANGUAGE (default=10), 
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
