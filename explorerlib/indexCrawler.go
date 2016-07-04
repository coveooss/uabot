package explorerlib

import (
	"github.com/coveo/go-coveo/search"
)

func FindWordsByLanguageInIndex(index Index, fields []string, documentsExplorationPercentage float64, fetchNumberOfResults int) (map[string]WordCounts, error) {
	wordCountsByLanguage := make(map[string]WordCounts)
	wordsByFieldValueByLanguage := map[string][]WordsByFieldValue{}
	languages, status := index.FetchLanguages()
	if status != nil {
		return nil, status
	}
	// for each language
	for _, language := range languages {
		// discover Words
		// for every fields provided
		for _, field := range fields {
			values, status := index.FetchFieldValues(field)
			if status != nil {
				return nil, status
			}
			// for all values of the field
			for _, value := range values.Values {
				wordCounts := WordCounts{}
				totalCount, status := index.FindTotalCountFromQuery(search.Query{
					AQ:"@syslanguage=\"" + language + "\" " + field + "=\"" + value.Value + "\"",
				})
				if status != nil {
					return nil, status
				}
				var queryNumber int
				if tempQueryNumber := (int(float64(totalCount) * documentsExplorationPercentage) / fetchNumberOfResults); tempQueryNumber > 0 {
					queryNumber = tempQueryNumber
				} else {
					queryNumber = 1
				}
				randomWord := ""

				for i := 0; i < queryNumber; i++ {
					// build A query from the word counts in the appropriate language with a filter on the field value
					queryExpression := randomWord +
					" @syslanguage=\"" + language + "\" " +
					field + "=\"" + value.Value + "\" "
					response, status := index.FetchResponse(queryExpression, fetchNumberOfResults)
					if status != nil {
						return nil, status
					}
					// extract words from the response
					newWordCounts := ExtractWordsFromResponse(*response)
					// update word counts
					wordCounts = wordCounts.Extend(newWordCounts)
					// pick a random word (Probability by popularity, or constant)
					randomWord = wordCounts.PickRandomWord()
				}
				taggedLanguage := LanguageToTag(language)
				wordsByFieldValueByLanguage[taggedLanguage] = append(wordsByFieldValueByLanguage[taggedLanguage], WordsByFieldValue{
					FieldName:field,
					FieldValue:value.Value,
					Words:wordCounts,
				})
			}
		}
	}
	// collapse results from all fields
	for language, wordCountsInLanguage := range wordsByFieldValueByLanguage {
		wordCounts := WordCounts{}
		for _, wordCountsByFields := range wordCountsInLanguage {
			wordCounts = wordCounts.Extend(wordCountsByFields.Words)
		}
		RankByWordCount(wordCounts)
		wordCountsByLanguage[language] = wordCounts
	}
	return wordCountsByLanguage, nil
}
