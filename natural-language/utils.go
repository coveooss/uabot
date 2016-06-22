package natural_language

import (
	"sort"
	"strings"
	"regexp"
)

func CountWordOccurence(words []string) WordCounts{
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
		wordCounts = append(wordCounts, WordCount{Key: key, Value: value})
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

type WordCount struct {
	Key string
	Value int
}

type WordCounts []WordCount

func (wordCountList WordCounts) Len() int { return len(wordCountList) }
func (wordCountList WordCounts) Less(i, j int) bool { return wordCountList[i].Value < wordCountList[j].Value }
func (wordCountList WordCounts) Swap(i, j int){ wordCountList[i], wordCountList[j] = wordCountList[j], wordCountList[i] }
func (wordCountList WordCounts) Add(pair WordCount) (WordCounts) {
	addedPairlist := WordCounts{}
	if wordCountList.ContainsKey(pair.Key) {
		for _, in_pair := range wordCountList {
			if in_pair.Key == pair.Key {
				addedPairlist = append(addedPairlist, WordCount{pair.Key, in_pair.Value + pair.Value})
			} else {
				addedPairlist = append(addedPairlist, in_pair)
			}
		}
	} else {
		addedPairlist = append(wordCountList, WordCount{pair.Key, pair.Value})
	}
	return addedPairlist
}
func (wordCounts WordCounts) ContainsKey(key string) (bool) {
	for _, pair := range wordCounts{
		if pair.Key == key {
			return true
		}
	}
	return false
}

func (firstPairList WordCounts) Extend(secondPairList WordCounts) (WordCounts) {
	mergedPairList := WordCounts{}
	for _, pair := range firstPairList {
		mergedPairList = append(mergedPairList, pair)
	}
	for _,pair := range secondPairList {
		mergedPairList = mergedPairList.Add(pair)
	}
	return mergedPairList
}
