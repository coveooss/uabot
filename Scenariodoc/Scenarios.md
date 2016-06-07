## Building a scenario

This is the documentation of the parameters of a scenario.

Parameters in **bold** are mandatory

Parameter | Type | Usage
------------ | ------------- | ----------------
name | string | Name of the scenario, used mainly for debugging purposes
weight | number | An arbitrary weight gived to the scenario
events | []Events | The events happening in this scenario

*The chance that a single specific scenario will be randomized is given by weight/totalWeights*

### Supported events

Type | Description
------------ | -------------
Search | Represents one query sent to the index. Typically the submit of the search bar, search as you type, etc.
Click | Represents a click on a document that was returned by a query. Can represent either a document open or a quickview.
Facet | Represents an event sent when the user chooses a value in a facet.
Tab | Represents when the user changes the tabs on top of the result list in a search page.
SearchAndClick | Use when you want to click on a specific document after a specific search. Ties a search and a click event together.
Custom | A custom event sent to the analytics, contains custom data.
SetOrigin | An event to tell the bot to change the origin of the events (use this when the user moved between search pages for example)
Page View | An event when a user visits a page.

*Consult each event doc for information on how to build a specific event.*

### Example

```json
{
    "name"   : "A name",
    "weight" : 1,
    "events" : [ ]
}
```
