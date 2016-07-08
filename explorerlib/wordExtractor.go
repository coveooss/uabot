package explorerlib

import (
	"github.com/coveo/go-coveo/search"
	"regexp"
	"strings"
)

func ExtractWordsFromResponse(response search.Response) WordCounts {
	titleWords := ExtractWordCountsFromTitlesInResponse(response)
	concepts := ExtractWordCountsFromConceptsInResponse(response)

	return titleWords.Extend(concepts)
}

func ExtractWordCountsFromTitlesInResponse(response search.Response) WordCounts {
	var text string
	for _, result := range response.Results {
		text = strings.Join([]string{text, result.Title}, " ")
	}
	text = CleanText(text)
	reg, _ := regexp.Compile("([\\w\\p{L}\\p{Nd}']+)")
	words := reg.FindAllString(text, -1)
	results := CountWordOccurence(words)
	return results
}

func ExtractWordCountsFromConceptsInResponse(response search.Response) WordCounts {
	conceptsList := WordCounts{}
	for _, groupBy := range response.GroupByResults {
		for _, concepts := range groupBy.Values {
			conceptsList = conceptsList.Add(WordCount{CleanText(concepts.Value), concepts.NumberOfResults})
		}
	}
	return conceptsList
}
