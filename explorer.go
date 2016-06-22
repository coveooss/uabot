package main

import (
	"github.com/adambbolduc/uabot/natural-language"
	"github.com/coveo/go-coveo/search"
	"log"
	"os"
)

var (
	Info *log.Logger
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var searchToken = os.Getenv("SEARCHTOKEN")

func main() {
	Info = log.New(os.Stdout, "INFO | ", log.Ldate|log.Ltime)

	query := search.Query{
		Q:               "@syslanguage=Catalan",
		NumberOfResults: 1000,
		GroupByRequests: []*search.GroupByRequest{
			&search.GroupByRequest{
				Field: "@source",
				MaximumNumberOfValues: -1,
				SortCriteria:          "nosort",
				InjectionDepth:        2000,
			},
		},
	}
	wordCounts, err := natural_language.ExtractWordsFromQuery(query)
	check(err)
	natural_language.RankByWordCount(wordCounts)
	Info.Println(wordCounts)
	searchClient, _ := search.NewClient(search.Config{
		Token:     searchToken,
		UserAgent: "",
		Endpoint:  "https://platformdev.cloud.coveo.com/rest/search/",
	})

	sources, status := searchClient.ListFacetValues("@sysauthor", 1000)
	check(status)
	for _, source := range sources.Values {
		Info.Println(source.Value)
	}
}
