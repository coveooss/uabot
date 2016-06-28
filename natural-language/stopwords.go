package natural_language

import (
	"io/ioutil"
	"strings"
)

type Stopwords struct {
	StopwordsList []string
}

func (stopwords *Stopwords) LoadFromFile(path string) (error) {
	file, err := ioutil.ReadFile(path)
	for _, word := range strings.Split(string(file), "\r\n") {
		stopwords.StopwordsList = append(stopwords.StopwordsList, word)
	}
	return err
}

func (stopwords *Stopwords) LoadRecursivelyFromDirectory(path string) (error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range dir {
		fileName := path + "/" + file.Name()
		if file.IsDir() {
			err := stopwords.LoadRecursivelyFromDirectory(fileName)
			if err != nil {
				return err
			}
		} else {
			stopwords.LoadFromFile(fileName)
		}
	}
	return nil
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
		if (stopword == word) {
			return true
		}
	}
	return false
}
