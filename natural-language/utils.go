package natural_language

import (
	"sort"
	"strings"
	"regexp"
)

func CountWordOccurence(words []string) WordCounts {
	word_occurence := make(map[string]int)
	for _, word := range words {
		if _, ok := word_occurence[word]; ok {
			word_occurence[word] += 1
		} else {
			word_occurence[word] = 1
		}
	}
	wordCounts := WordCounts{}
	for key, value := range word_occurence {
		wordCounts.Words = append(wordCounts.Words, WordCount{Word: key, Count: value})
	}
	return wordCounts
}

func CleanText(text string) string {
	digitRegex := regexp.MustCompile("[0-9]")
	removedDigit := digitRegex.ReplaceAllString(text, "")
	toLower := strings.ToLower(removedDigit)
	return toLower
}

func RankByWordCount(wordCounts WordCounts) WordCounts {
	sort.Sort(sort.Reverse(wordCounts))
	return wordCounts
}

type WordsByFieldValue struct {
	FieldName  string
	FieldValue string
	Words      WordCounts
}


