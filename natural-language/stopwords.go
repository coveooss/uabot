package natural_language

import (
	"io/ioutil"
	"strings"
)

type Stopwords struct {
	StopwordsList []string
}

func (stopwords *Stopwords) LoadFromFile(path string) (error){
	file, err := ioutil.ReadFile(path)
	stopwords.StopwordsList = strings.Split(string(file), "\r\n")
	return err
}

func (stopwords *Stopwords) RemoveFrom(words []string) []string {
	filteredWords := []string{}
	for _, word := range words {
		if (!stopwords.Contains(word)) {
			filteredWords = append(filteredWords, word)
		}
	}
	return filteredWords
}

func (stopwords *Stopwords) Contains(word string) bool {
	for _, stopword := range stopwords.StopwordsList {
		if(stopword == word){
			return true
		}
	}
	return false
}
