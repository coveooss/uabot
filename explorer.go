package main

import (
	"log"
	"os"
	"github.com/erocheleau/uabot/scenariolib"
	"encoding/json"
	"io/ioutil"
	"github.com/coveo/go-coveo/search"
	"strings"
	"regexp"
)

var (
	Info *log.Logger
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	Info = log.New(os.Stdout, "INFO | ", log.Ldate|log.Ltime)

	var stopwords natural_language.Stopwords

	err := stopwords.LoadFromFile("stopwords/en.txt")
	check(err)

	var config scenariolib.Config

	config = scenariolib.Config{}

	config.OrgName = "orgName"


	searchToken := os.Getenv("SEARCHTOKEN")

	searchClient, _:= search.NewClient(search.Config{
		Token: searchToken,
		UserAgent:"",
		Endpoint:"https://platformdev.cloud.coveo.com/rest/search/",
	});

	response, status := searchClient.Query(search.Query{
		Q: "rabbitmq",
		NumberOfResults: 1000,
		GroupByRequests: []*search.GroupByRequest{&search.GroupByRequest{
			Field: "@concepts",
			MaximumNumberOfValues:10,
			SortCriteria: "occurences",
			InjectionDepth: 2000,
		}},
	})
	check(status)


	conceptsList := natural_language.PairList{}

	for _, groupBy := range response.GroupByResults {
		for _, concepts := range groupBy.Values {
			conceptsList = append(conceptsList, natural_language.Pair{natural_language.CleanText(concepts.Value), concepts.NumberOfResults})
		}
	}

	var text string
	for _, result := range response.Results {
		text = strings.Join([]string{text, result.Title}, " ")
	}

	reg, _ := regexp.Compile("(\\w+)")

	text = natural_language.CleanText(text)

	words := reg.FindAllString(text, -1)

	filteredWords := stopwords.RemoveFrom(words)

	config.GoodQueries = filteredWords

	jsonConfig, err:= json.Marshal(config)

	check(err)

	fileerr := ioutil.WriteFile("outputScenario.json", jsonConfig, 0644)
	check(fileerr)

	result := natural_language.RankByWordCount(natural_language.CountWordOccurence(filteredWords))

	Info.Println(natural_language.MergePairLists(result, conceptsList))
}

