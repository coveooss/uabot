# Events
Different types of events can be generated using this bot. This is the documentation on how to build them.

## List of events
1. [Search event](#Search)
2. [Click event](#Click)
3. [SearchAndClick event](#SearchAndClick)
4. [Custom event](#Custom)
5. [TabChange event](#Tab)
6. [FacetChange event](#Facet)
7. [SetOrigin event](#Origin)
8. [PageView event](#Page)

### 0. Generic event

Parameter | Type | Usage
------------ | ------------- | ----------------
type | string | The type of the event
arguments | Object | The arguments of the event, they are different for each type of events

### <a name="Search"></a> 1. Search event
Represents one query sent to the index. Typically the submit of the search bar, search as you type, etc.

`"type" : "Search"`

Arguments | Type | Usage
------------ | ------------- | ----------------
queryText | string | The query to send. Leave "blank" for a random query
goodQuery | boolean | If the random query should be a good or a bad query
ignoreEvent | boolean | Do not send the event to analytics (optional, default is false)
matchLanguage | boolean | If the query expression will be in the visit language.
customData | object | Custom data to be sent alongside the event.
caseSearch | boolean | If the query comes from a Case Creation interface
inputTitle | string | (Only used if caseSearch is true) Name of the input field on the case form that triggered the search.

#### Example

```json
{
    "type" : "Search",
    "arguments" : {
        "queryText" : "",
        "goodQuery" : true,
        "logEvent" : true,
        "caseSearch" : true,
        "inputTitle" : "Product",
        "customData" : {
            "hasclicks": true
        }
    }
}
```

###<a name="Click"></a> 2. Click Event

Represents a click on a document that was returned by a query. Can represent either a document open or a quickview.

`"type" : "Click"`

Arguments | Type | Usage
------------ | ------------- | ----------------
docNo | number | The rank of the document to click (0 base, put -1 for random)
offset | number | An offset used in random document clicking
probability | number | The probability that the user will click (between 0 and 1)
quickview | boolean | If the click is a quickview or not (default, false)
customData | object | Custom data to be sent alongside the event.
fakeClick | boolean | Click on a document in falseResponse.
fakeResponse | search.Response | A fake response from the search

#### Example

```json
{
    "type" : "Click",
    "arguments" : {
        "docNo" :-1,
        "offset" : 0,
        "probability" : 0.45,
        "quickview" : true,
        "customData" : {
            "hasclicks": true
        }
    }
}
```

###<a name="SearchAndClick"></a> 3. SearchAndClick event

Use when you want to click on a specific document after a specific search. Ties a search and a click event together.

`"type" : "SearchAndClick"`

Arguments | Type | Usage
------------ | ------------- | ----------------
**queryText** | string | The query to send. Not recommended to use with random query.
**docClickTitle** | string | The title of the document you want to click on.  (checks if string is _contained_ in title.) You can use either **docClickTitle** or the **matchField**/**matchRegex** pair.
**matchField** | string | The name of the field you want to use to find the document you want to click on. (*matchRegex* is required when using this argument)
**matchRegex** | string | The pattern you want to match to find the document you want to click on. (*matchField* is required when using this argument)
**probability** | number | Between 0 and 1, the probability the user will click
quickview | boolean | If the click is a quickview instead of a document open (default false)
caseSearch | boolean | If the event is on a Case Creation interface (default false)
inputTitle | string | If it's a case creation event, which input triggered the search
customData | object | Any custom data to send with the event

#### Example
```json
{
    "type" : "SearchAndClick",
    "arguments" : {
        "queryText" : "specific query",
        "caseSearch": true,
        "inputTitle": "Subject",
        "probability" : 0.85,
        "docClickTitle" : "specific title"
    }
}
```
```json
{
    "type" : "SearchAndClick",
    "arguments" : {
        "queryText" : "specific query",
        "caseSearch": true,
        "inputTitle": "Subject",
        "probability" : 0.85,
        "matchField" : "title",
        "matchRegexp" : "^Rocky(\\s+[IV]+)*$"
    }
}
```

###<a name="Custom"></a> 4. Custom event

A custom event sent to the analytics, contains custom data.

`"type" : "Custom"`

Arguments | Type | Usage
------------ | ------------- | ----------------
**actionType** | string | The event type of this custom event
**actionCause** | string | The cause of this event, also the event value
customData | object | Any custom data to send with the event

#### Example
```json
{
    "type" : "Custom",
    "arguments" : {
        "actionCause" : "submitButton",
        "actionType" : "caseCreation",
        "customData" : {
            "hasclicks": false,
            "product" : "XBR6 TV"
        }
    }
}
```

###<a name="Tab"></a> 5. TabChange event

Represents when the user changes the tabs on top of the result list in a search page.

`"type" : "TabChange"`

Arguments | Type | Usage
------------ | ------------- | ----------------
**tabName** | string | The name of the tab that the user switched to. This will also change originLevel2.
**tabCQ** | string | The constant query applied by this tab to the queries

#### Example
```json
{
    "type" : "TabChange",
    "arguments" : {
        "tabName" : "YOUTUBE",
        "tabCQ" : "@filetype==\"youtubevideo\""
    }
}
```

###<a name="Facet"></a> 6. FacetChange event

Represents an event sent when the user chooses a value in a facet.

`"type" : "FacetChange"`

Arguments | Type | Usage
------------ | ------------- | ----------------
**facetTitle** | string | The title of the facet that was selected
**facetValue** | string | The value that was selected in the facet
**facetField** | string | The field bound to the facet

#### Example
```json
{
    "type" : "FacetChange",
    "arguments" : {
        "facetTitle": "Type",
        "facetValue": "Message",
        "facetField": "@objecttype"
    }
}
```

###<a name="Origin"></a> 7. SetOrigin event

An event to tell the bot to change the origin of the events (use this when the user moved between search pages for example)

`"type" : "SetOrigin"`

Arguments | Type | Usage
------------ | ------------- | ----------------
originLevel1 | string | The new originLevel1
originLevel2 | string | The new originLevel2
originLevel3 | string | The new originLevel3

#### Example
```json
{
    "type" : "SetOrigin",
    "arguments" : {
        "originLevel1" : "Example1",
        "originLevel2" : "Example2",
        "originLevel3" : "Example3"
    }
}
```

###<a name="Page"></a> 8. PageView event

An event when a user visits a page. (This event is similar to a Click event, do a Search, then instead of a click, send a view event).

`"type" : "View"`

Arguments | Type | Usage
------------ | ------------- | ----------------
offset | number | And offset to send a view event after a search.
**probability** | number | Between 0 and 1, the probability the user will click
**docNo** | number | The rank of the document to click (0 base, put -1 for random)
**pageViewField** | string | The field to use on the result as contentIdKey and contentIdValue
contentType | string | The type of the content that was viewed

#### Example
```json
{
    "type" : "View",
    "arguments" : {
        "pageuri" : "https://example.com",
        "pagereferrer" : "https://www.google.com",
        "pagetitle" : "Page Example"
    }
}

```
