# Events
Different types of events can be generated using this bot. This is the documentation on how to build them.

## List of events
1. Search event
2. Click event
3. SearchAndClick event
4. Custom event
5. TabChange event
6. FacetChange event
7. SetOrigin event
8. PageView event

### 0. Generic event

Parameter | Type | Usage
------------ | ------------- | ----------------
type | string | The type of the event
arguments | Object | The arguments of the event, they are different for each type of events

### 1. Search event
Represents one query sent to the index. Typically the submit of the search bar, search as you type, etc.

`"type" : "Search"`

Arguments | Type | Usage
------------ | ------------- | ----------------
queryText | string | The query to send. Leave "blank" for a random query
goodQuery | boolean | If the random query should be a good or a bad query
caseSearch | boolean | If the query comes from a Case Creation interface
inputTitle | string | The title of the input that triggered the search if it was a case search
customData | object | Custom data to be sent alongside the event.

#### Example

```json
{
    "type" : "Search",
    "arguments" : {
        "queryText" : "",
        "goodQuery" : true,
        "caseSearch" : true,
        "inputTitle" : "Product",
        "customData" : {
            "hasclicks": true
        }
    }
}
```

### 2. Click Event

Represents a click on a document that was returned by a query. Can represent either a document open or a quickview.

`"type" : "Click"`

Arguments | Type | Usage
------------ | ------------- | ----------------
docNo | number | The rank of the document to click (0 base, put -1 for random)
offset | number | An offset used in random document clicking
probability | number | The probability that the user will click (between 0 and 1)
quickview | boolean | If the click is a quickview or not (default, false)
customData | object | Custom data to be sent alongside the event.

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
