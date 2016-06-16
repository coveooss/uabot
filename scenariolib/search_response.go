
type Result struct {
  URI string `json:URI`
  Title string `json:Title`
  ClickUri string `json:ClickUri`
  sysurihash string `json:sysurihash`
  syscollection string `json:syscollection`
  syssource string `json:syssource`
}

type SearchResponse struct {
  SearchUID string `json:SearchUID`
  TotalCount int `json:TotalCount`
  Results []Result `json:Results`
}
