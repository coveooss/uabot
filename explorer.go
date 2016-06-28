package main

import (
	"log"
	"os"
	"github.com/adambbolduc/uabot/natural-language"
	"github.com/coveo/go-coveo/search"
	"strconv"
	"github.com/erocheleau/uabot/scenariolib"
	"io/ioutil"
	"encoding/json"
)

var (
	Info *log.Logger
	searchToken = os.Getenv("SEARCHTOKEN")
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	Info = log.New(os.Stdout, "INFO | ", log.Ldate | log.Ltime)
	index, status:= natural_language.NewIndex("https://platformdev.cloud.coveo.com/rest/search/", searchToken)
	check(status)
	fields := []string{"@syssource"}

	wordsByFieldValueByLanguage := map[string][]natural_language.WordsByFieldValue{}
	// for each language
	languages, status := index.FetchLanguages()
	check(status)
	for _, language := range languages {
		// discover Words
		wordsByFieldValueByLanguage[language] = []natural_language.WordsByFieldValue{}
		// for every fields provided
		for _, field := range fields {
			values, status := index.FetchFieldValues(field)
			check(status)
			// for all values of the field
			for _, value := range values.Values {
				wordCounts := natural_language.WordCounts{}
				totalCount, status := index.FindTotalCountFromQuery(search.Query{
					AQ:"@syslanguage=\"" + language + "\" " + field + "=\"" + value.Value + "\"",
				})
				check(status)
				queryNumber := totalCount / 100 + 1
				randomWord := ""

				Info.Println(strconv.Itoa(queryNumber), language, field, value.Value)
				for i := 0; i < queryNumber; i++ {
					// build A query from the word counts in the appropriate language with a filter on the field value
					queryExpression := randomWord +
					" @syslanguage=\"" + language + "\" " +
					field + "=\"" + value.Value + "\" "
					response, status := index.FetchResponse(queryExpression, 100)
					check(status)
					// extract words from the response
					newWordCounts := natural_language.ExtractWordsFromResponse(*response)
					// update word counts
					wordCounts = wordCounts.Extend(newWordCounts)
					// pick a random word (Probability by popularity, or constant)
					randomWord = wordCounts.PickRandomWord()
				}
				wordsByFieldValueByLanguage[language] = append(wordsByFieldValueByLanguage[language], natural_language.WordsByFieldValue{
					FieldName:field,
					FieldValue:value.Value,
					Words:wordCounts,
				})
			}
		}
	}
	// collapse results from all fields
	wordCountsByLanguage := map[string]natural_language.WordCounts{}
	for language, wordCountsInLanguage := range wordsByFieldValueByLanguage {
		wordCounts := natural_language.WordCounts{}
		for _, wordCountsByFields := range wordCountsInLanguage {
			wordCounts = wordCounts.Extend(wordCountsByFields.Words)
		}
		natural_language.RankByWordCount(wordCounts)
		wordCountsByLanguage[language] = wordCounts
	}

	queriesInLanguage := make(map[string][]string)
	for language, wordCounts := range wordCountsByLanguage {
		words := []string{}
		for i, word := range wordCounts.Words{
			if i < 20 {
				words = append(words, word.Word)
			} else {
				break
			}
		}
		queriesInLanguage[language] = words
	}
	config := scenariolib.Config{}
	config.GoodQueriesInLang = queriesInLanguage
	bytes, err := json.Marshal(config)
	check(err)
	writeerr := ioutil.WriteFile("out.json", bytes, 0644)
	check(writeerr)
}
