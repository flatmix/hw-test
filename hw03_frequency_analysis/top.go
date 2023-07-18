package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordItem struct {
	Count int
	Word  string
}

func Top10(text string) []string {
	words := strings.Fields(text)
	wordsMap := make(map[string]int, 0)

	for _, wordKey := range words {
		word := strings.ToLower(strings.Trim(wordKey, `,-!.'"`))
		if len(word) == 0 {
			continue
		}
		wordsMap[word]++
	}

	wordsArr := make([]WordItem, len(wordsMap))

	i := 0
	for word, count := range wordsMap {
		wordsArr[i] = WordItem{
			Count: count,
			Word:  word,
		}
		i++
	}

	sort.Slice(wordsArr, func(i, j int) bool {
		if wordsArr[i].Count != wordsArr[j].Count {
			return wordsArr[i].Count > wordsArr[j].Count
		}
		return wordsArr[i].Word < wordsArr[j].Word
	})

	result := []string{}

	for iw, structWord := range wordsArr {
		if iw < 10 {
			result = append(result, structWord.Word)
		} else {
			break
		}
	}

	return result
}
