package natural_language

import (
	"os"
	"github.com/coveo/go-coveo/search"
	"strings"
	"regexp"
)

var searchToken = os.Getenv("SEARCHTOKEN")

func ExtractWordsFromQuery(query search.Query) (WordCounts, error) {
	titleWords, err := ExtractWordCountsFromTitlesInQuery(query)
	if err != nil {
		return nil, err
	}

	concepts, err := ExtractWordCountsFromConceptsInQuery(query)
	if err != nil {
		return nil, err
	}

	return titleWords.Extend(concepts), nil
}

func ExtractWordCountsFromTitlesInQuery(q search.Query) (WordCounts, error){
	var stopwords Stopwords

	err := stopwords.LoadFromFile("stopwords/en.txt")
	if err != nil {
		return nil, err
	}


	searchClient, _ := search.NewClient(search.Config{
		Token: searchToken,
		UserAgent:"",
		Endpoint:"https://platformdev.cloud.coveo.com/rest/search/",
	});

	response, status := searchClient.Query(q)
	if status != nil {
		return nil, status
	}

	var text string
	for _, result := range response.Results {
		text = strings.Join([]string{text, result.Title}, " ")
	}

	text = CleanText(text)

	reg, _ := regexp.Compile("([\\w\\p{L}\\p{Nd}']+)")
	words := reg.FindAllString(text, -1)
	filteredWords := stopwords.RemoveFrom(words)
	results := CountWordOccurence(filteredWords)
	return results, nil
}

func ExtractWordCountsFromConceptsInQuery(q search.Query) (WordCounts, error) {
	conceptsList := WordCounts{}
	searchClient, _ := search.NewClient(search.Config{
		Token: searchToken,
		UserAgent:"",
		Endpoint:"https://platformdev.cloud.coveo.com/rest/search/",
	});
	response, status := searchClient.Query(q)
	if status != nil {
		return nil, status
	}

	for _, groupBy := range response.GroupByResults {
		for _, concepts := range groupBy.Values {
			conceptsList = append(conceptsList, WordCount{CleanText(concepts.Value), concepts.NumberOfResults})
		}
	}
	return conceptsList, nil
}
