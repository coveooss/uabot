# BESTTECH

These scenarios are principally for showing customer analytics cloud platform. 
# Coaxial Cable to HDMI
### Events:
```json
{ "type" : "Search", "arguments" : { "queryText" : "install hdmi-QAM", "goodQuery" : false } },
{ "type" : "SearchAndClick", "arguments" : { "queryText" : "convert coaxial cable to hdmi", "docClickTitle" : "How to Convert Coaxial Cable to HDMI", "goodQuery" : false } }
```

### Explication: 
Reveal learns when a client click a document in a second search. 
### Expected Results: 
Reveal should learns and show the doc with the title «convert coaxial cable to hdmi»
### Actual Results: 
Without reveal you should not see your solution because the products «hdmi-QAM» is not a listed products. 

# HX720-478
### Events:
```json
{ "type" : "Search", "arguments" : { "queryText" : "hx720-478", "goodQuery" : false } },
{ "type" : "SearchAndClick", "arguments" : { "queryText" : "hx720-478 stand", "docClickTitle" : "HX720 Internet TV Stand is unstable", "goodQuery" : false } }
```
### Explication:
Client has trouble installing the stand for is new HX720 TV. Start a new query with the complete serial number «hx720-478» but, the model is only hx720. After finding nothing the client asks for «HX720-478 stand» and then find a resolved case that show him step by step. 
### Expected Results: 
KB with the title : HX720 Internet TV Stand is unstable.
### Actual Results:
No results

# Net flix not connect
### Events:
```json
{ "type" : "Search", "arguments" : { "queryText" : "Net flix not connected", "goodQuery" : false } },
{ "type" : "SearchAndClick", "arguments" : { "queryText" : "netflix not connect", "docClickTitle" : "Netflix setup on some televisions", "goodQuery" : false } }
```
### Explication:
Client has trouble installing netflix and search for net flix not connect. 
### Expected Results: 
Reveal learns that net flix and netflix are the same and that the doc title «Netflix setup on some televisions» appear.
### Actual Results:
Results but not relevant. 

## Search: What is 4k? To show quick view with youtube  



