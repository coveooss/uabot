package explorerlib

import (
	"time"
	"math/rand"
	"math"
	"strings"
)

var (
	s1 = rand.NewSource(time.Now().UnixNano())
	random = rand.New(s1)
	stopwords *Stopwords
)

type WordCount struct {
	Word  string
	Count int
}

type WordCounts struct {
	Words      []WordCount
	TotalCount int
}

func (wordCountList WordCounts) Len() int {
	return len(wordCountList.Words)
}
func (wordCountList WordCounts) Less(i, j int) bool {
	return wordCountList.Words[i].Count < wordCountList.Words[j].Count
}
func (wordCountList WordCounts) Swap(i, j int) {
	wordCountList.Words[i], wordCountList.Words[j] = wordCountList.Words[j], wordCountList.Words[i]
}
func (wordCountList WordCounts) Add(pair WordCount) (WordCounts) {
	// Do not Add a word that is in the stopword list
	loadStopWords()
	if stopwords.Contains(pair.Word) {
		return wordCountList
	}

	addedPairlist := WordCounts{}
	if wordCountList.ContainsKey(pair.Word) {
		for _, in_pair := range wordCountList.Words {
			if in_pair.Word == pair.Word {
				addedPairlist.Words = append(addedPairlist.Words, WordCount{pair.Word, in_pair.Count + pair.Count})
			} else {
				addedPairlist.Words = append(addedPairlist.Words, in_pair)
			}
		}
	} else {
		addedPairlist.Words = append(wordCountList.Words, WordCount{pair.Word, pair.Count})
	}
	addedPairlist.TotalCount += pair.Count
	return addedPairlist
}

func loadStopWords() {
	if stopwords == nil {
		stopwords = &Stopwords{}
		err := stopwords.LoadRecursivelyFromDirectory("stopwords")
		if err != nil {
			panic(err)
		}
	}
}

func (wordCounts WordCounts) ContainsKey(key string) (bool) {
	for _, pair := range wordCounts.Words {
		if pair.Word == key {
			return true
		}
	}
	return false
}

func (firstPairList WordCounts) Extend(secondPairList WordCounts) (WordCounts) {
	mergedPairList := WordCounts{}
	for _, pair := range firstPairList.Words {
		mergedPairList = mergedPairList.Add(pair)
	}
	for _, pair := range secondPairList.Words {
		mergedPairList = mergedPairList.Add(pair)
	}
	return mergedPairList
}

func (wordCounts WordCounts) PickRandomWord() string {
	if size := len(wordCounts.Words); size != 0 {
		return wordCounts.Words[random.Intn(size)].Word
	}
	return ""
}

func (wordCounts WordCounts) PickExpNWords(n int) string {
	numberOfWords := randomNumberWithExpMinMax(1, math.MaxInt64, float64(n))
	words := make([]string, 0)
	for i := 0; i < numberOfWords; i++ {
		words = append(words, wordCounts.PickRandomWord())
	}
	return strings.Join(words, " ")
}

func randomNumberWithExpMinMax(min int, max int, lambda float64) int {
	var exponentialrandomint = math.MaxInt64
	for min > exponentialrandomint || exponentialrandomint >= max {
		exponentialrandomint = int(random.ExpFloat64() * lambda + 0.5)
	}
	return exponentialrandomint
}
