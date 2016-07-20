## Building a scenario

This is the documentation of the parameters of a scenario.

Parameters in **bold** are mandatory

Parameter | Type | Usage
------------ | ------------- | ----------------
name | string | Name of the scenario, used mainly for debugging purposes
language | string | The language of the given scenario or random if not specified
weight | number | An arbitrary weight gived to the scenario
events | []Events | The events happening in this scenario

*The chance that a single specific scenario will be randomized is given by weight/totalWeights*

### Supported events [Documentation](events.md)

Type | Description
------------ | -------------
[Search](events.md#Search) | Represents one query sent to the index. Typically the submit of the search bar, search as you type, etc.
[Click](events.md#Click) | Represents a click on a document that was returned by a query. Can represent either a document open or a quickview.
[Facet](events.md#Facet) | Represents an event sent when the user chooses a value in a facet.
[Tab](events.md#Tab) | Represents when the user changes the tabs on top of the result list in a search page.
[SearchAndClick](events.md#SearchAndClick) | Use when you want to click on a specific document after a specific search. Ties a search and a click event together.
[Custom](events.md#Custom) | A custom event sent to the analytics, contains custom data.
[SetOrigin](events.md#Origin) | An event to tell the bot to change the origin of the events (use this when the user moved between search pages for example)
[Page View](events.md#View) | An event when a user visits a page.
[FakeSearch](events.md#FakeSearch) | Fake a search by getting a valid searchID but replacing the response to the fakeResponse

[*Consult each event doc for information on how to build a specific event.*](events.md)

### Example

```json
{
    "name"   : "A name",
    "weight" : 1,
    "events" : [ ]
}
```
