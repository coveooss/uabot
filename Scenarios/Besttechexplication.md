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



