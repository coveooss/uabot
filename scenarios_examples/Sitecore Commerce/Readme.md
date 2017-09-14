The scenarios above are for the Habitat Retail demo solution: https://github.com/Sitecore/Sitecore.Demo.Retail

The files are split in 4 Sitecore xDB personas native to the demo:

* Gamer
* Audiophile
* Home appliances buyer
* Kitchen appliances buyer

Which will each use several devices/browsers:

* Windows 10 with chrome/IE
* MAC with chrome
* Android with Chrome
* IOS with Safari

The current scenario is pointing at a live coveo demo site : retail.coveodemo.com which is tied to a Coveo cloud platform. Most of the nodes are exactly what is used for that demo so that you can see a live working example instead of a placeholder.

Which means that in order to use these scenarios, you must change the following nodes:

* orgName => Set your own Coveo Cloud organisation ID (not the display name, but the ID found in the url).
* globalFilter => The filters are replicating those of the demo search page. Simply create a search page and inspect the aq and cq parameters found in the request sent when performing a query.
* origin (all level) for all scenarios => The originLevel1, originLevel2 and originLevel3 are tied to the current demo site, simply change them for the DOM ID of the search page, the DOM ID of the main tab of the search interface (default if none) and the url of the page redirecting to the search page. If you have multiple tabs and you want to randomize them, you will need to create different scenarios with different origins. For the origin 3, some scenarios will target the home page, other the search page itself and some specific pages.
* referrer for all scenarios => Should be the same as the origin 3.

Be aware that these scenarios might not work if the base code of the UABot changes. If anything goes wrong please post on answers.coveo.com.
