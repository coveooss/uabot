package natural_language

import (
	"sort"
	"strings"
	"regexp"
)

func CountWordOccurence(words []string) map[string]int {
	word_occurence := make(map[string]int)
	for _, word := range words {
		if _, ok := word_occurence[word]; ok {
			word_occurence[word] += 1
		} else {
			word_occurence[word] = 1
		}
	}
	return word_occurence
}

func CleanText(text string) string {
	digitRegex := regexp.MustCompile("[0-9]")
	removedDigit := digitRegex.ReplaceAllString(text, "")
	toLower := strings.ToLower(removedDigit)
	return toLower
}

func RankByWordCount(wordFrequencies map[string]int) PairList{
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key string
	Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }
func (p PairList) Add(pair Pair) (PairList) {
	addedPairlist := PairList{}
	if p.ContainsKey(pair.Key) {
		for _, in_pair := range p {
			if in_pair.Key == pair.Key {
				addedPairlist = append(addedPairlist, Pair{pair.Key, in_pair.Value + pair.Value})
			}
			addedPairlist = append(addedPairlist, in_pair)
		}
	} else {
		addedPairlist = append(p, Pair{pair.Key, pair.Value})
	}
	return addedPairlist
}
func (p PairList) ContainsKey(key string) (bool) {
	for _, pair := range p {
		if pair.Key == key {
			return true
		}
	}
	return false
}

func MergePairLists(firstPairList PairList, secondPairList PairList) (PairList) {
	mergedPairList := PairList{}
	for _, pair := range firstPairList {
		mergedPairList = append(mergedPairList, pair)
	}
	for _,pair := range secondPairList {
		mergedPairList = mergedPairList.Add(pair)
	}
	return mergedPairList
}
