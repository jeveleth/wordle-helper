package utils

// check if the word is valid
// if valid, add to list of valid words
// if not valid, skip to next word
// return list of valid words

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/dlclark/regexp2"
)

var clean []string

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (a PairList) Len() int           { return len(a) }
func (a PairList) Less(i, j int) bool { return a[i].Value < a[j].Value }
func (a PairList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type Config struct {
	Include     string
	Exclude     string
	Pattern     string
	AntiPattern string
}

// GetCommonLettersCount returns a map of letters and their frequency
func GetCommonLettersCount(words []string) map[string]int {
	str := strings.Join(words, " ")
	smashString := strings.ReplaceAll(str, " ", "")
	m := make(map[string]int)
	var b strings.Builder

	for _, char := range smashString {
		if _, ok := m[string(char)]; ok {
			m[string(char)]++
		} else {
			m[string(char)] = 1
			b.WriteString(string(char))
		}
	}
	return m
}

// RanksByFreq returns a slice of pairs of words and their frequency
func RanksByFreq(letterFreq map[string]int) PairList {
	pl := make(PairList, len(letterFreq))
	i := 0
	// first sort by value to get frequency
	for k, v := range letterFreq {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func SortRankedPairsAlaphabetically(pl PairList) PairList {
	// loop through pl to check if there are multiple letters with same frequency
	// if so, add to slice of pairs
	newPairList := make(PairList, 0)

	var sameVals []Pair
	for i, v := range pl {
		if i == 0 {
			sameVals = append(sameVals, v)
			continue
		} else {
			// if the current value is the same as the last value in sameVals
			if v.Value == sameVals[len(sameVals)-1].Value {
				sameVals = append(sameVals, v)
				// sort the slice of pairs by value
				sort.Slice(sameVals, func(i, j int) bool {
					return sameVals[i].Key < sameVals[j].Key
				})
			} else {
				// if the current value is not the same as the previous value
				// remove the previous value from the slice
				newPairList = append(newPairList, sameVals...)
				sameVals = nil
				sameVals = append(sameVals, v)
			}
		}
	}
	if sameVals != nil {
		newPairList = append(newPairList, sameVals...)
		sameVals = nil
	}
	return newPairList
}

// RankCommonLetters returns a slice of letters and their frequency
func RankCommonLetters(words []string) PairList {
	letterFreq := GetCommonLettersCount(words)
	pl := RanksByFreq(letterFreq)
	return SortRankedPairsAlaphabetically(pl)
}

func excludeOnly(b strings.Builder, wrds []string, c Config) string {
	excluded := SetExcludedLetters(c.Exclude, wrds)
	b.WriteString(fmt.Sprintf(
		`%v words found that exclude '%v'.`,
		len(excluded),
		c.Exclude),
	)
	b.WriteString("The words are:\n")
	for _, v := range excluded {
		b.WriteString(fmt.Sprintf("%v\n", v))
	}
	b.WriteString("\nThe most common letters are:\n")
	pl := RankCommonLetters(excluded)
	for _, v := range pl {
		b.WriteString(fmt.Sprintf("\n%v: %v times", v.Key, v.Value))
	}
	return b.String()
}

func includeOnly(b strings.Builder, wrds []string, c Config) string {
	included := SetIncludeLetters(c.Include, wrds)
	b.WriteString(fmt.Sprintf(
		`%v words found that include '%v'. `,
		len(included),
		c.Include))
	b.WriteString("The words are:\n")
	for _, v := range included {
		b.WriteString(fmt.Sprintf("%v\n", v))
	}
	b.WriteString("\nThe most common letters are:\n")
	pl := RankCommonLetters(included)
	for _, v := range pl {
		b.WriteString(fmt.Sprintf(
			"\n%v: %v times",
			v.Key,
			v.Value))
	}
	return b.String()
}

func processEmptyPattern(b strings.Builder, c Config, parsedWords []string) string {
	b.WriteString(
		fmt.Sprintf(`%v words found that exclude '%v' and include '%v'. `, len(parsedWords), c.Exclude, c.Include),
	)
	b.WriteString("The words are:\n")
	for _, wd := range parsedWords {
		b.WriteString(fmt.Sprintf("%v\n", wd))
	}
	b.WriteString("\nThe most common letters are:\n")
	pl := RankCommonLetters(parsedWords)
	for _, v := range pl {
		b.WriteString(fmt.Sprintf("\n%v: %v times", v.Key, v.Value))
	}
	return b.String()
}

func processNonEmptyPattern(b strings.Builder, c Config, parsedWords []string) string {
	str := ParsePattern(c.Pattern)
	finalWords := SetComboPattern(str, parsedWords)
	b.WriteString(fmt.Sprintf(`%v words found that exclude '%v', include '%v', and match the pattern '%v'. `,
		len(finalWords),
		c.Exclude,
		c.Include,
		c.Pattern,
	))
	b.WriteString("The words are:\n")
	for _, wd := range finalWords {
		b.WriteString(fmt.Sprintf("%v\n", wd))
	}
	b.WriteString("\nThe most common letters are:\n")
	pl := RankCommonLetters(finalWords)
	for _, v := range pl {
		b.WriteString(fmt.Sprintf("\n%v: %v times", v.Key, v.Value))
	}
	return b.String()
}

func includeAndExclude(b strings.Builder, wrds []string, c Config) string {
	var parsedWords []string
	clean = SetExcludedLetters(c.Exclude, SetIncludeLetters(c.Include, wrds))
	parsedWords = clean
	if c.AntiPattern != "" {
		str := ParseAntiPattern(c.AntiPattern)
		parsedWords = SetAntiPattern(str, clean)
	}
	if c.Pattern == "" {
		return processEmptyPattern(b, c, parsedWords)
	} else {
		return processNonEmptyPattern(b, c, parsedWords)
	}
}

// FiveLetterWords returns a slice of words that are five letters long
func FiveLetterWords(wrds []string) []string {
	var fiveLetterWords []string
	for _, wd := range wrds {
		if len(wd) == 5 {
			fiveLetterWords = append(fiveLetterWords, wd)
		}
	}
	return fiveLetterWords
}

// ReadWords reads words from a file
func ReadWords() ([]string, error) {
	wrds := make([]string, 0)
	f, err := os.Open("words.txt")

	if err != nil {
		return wrds, err
	}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		wrds = append(wrds, scanner.Text())
	}
	// don't close too soon
	defer f.Close()
	return wrds, nil
}

// SetIncludeLetters returns a slice of words that include the letters in the string
func SetIncludeLetters(letters string, wrds []string) []string {
	var b strings.Builder
	for _, char := range letters {
		b.WriteString(fmt.Sprintf(`(?=.*(%[1]v))`, strings.ToLower(string(char))))
	}
	// We use regexp2 because the stdlib regexp doesn't support positive lookahead
	re2 := regexp2.MustCompile(fmt.Sprintf(`%v.*`, b.String()), 0)
	var wrds2 []string
	for _, v := range wrds {
		if isMatch, _ := re2.MatchString(v); isMatch {
			wrds2 = append(wrds2, v)
		}
	}
	return wrds2
}

// ParseAntiPattern parses the anti-pattern string
func ParseAntiPattern(antiPattern string) string {
	re := regexp.MustCompile("-")
	var b strings.Builder
	for _, char := range antiPattern {
		if re.MatchString(string(char)) {
			b.WriteString(`([a-z])`)
		} else {
			b.WriteString(fmt.Sprintf(`([^%[1]v])`, strings.ToLower(string(char))))
		}
	}
	return b.String()
}

// ParsePattern parses the pattern string
func ParsePattern(pattern string) string {
	re := regexp.MustCompile("-")
	var b strings.Builder
	for _, char := range pattern {
		if re.MatchString(string(char)) {
			b.WriteString(`([a-z])`)
		} else {
			b.WriteString(fmt.Sprintf(`[%[1]v]`, strings.ToLower(string(char))))
		}
	}
	return b.String()
}

// SetAntiPattern returns a slice of words that match the anti-pattern
func SetAntiPattern(pattern string, wrds []string) []string {
	re := regexp.MustCompile(pattern)
	var newWords []string
	for _, wrd := range wrds {
		if re.MatchString(wrd) {
			newWords = append(newWords, wrd)
		}
	}
	return newWords
}

// SetExcludedLetters returns a slice of words that exclude the letters in the string
func SetExcludedLetters(exclude string, wrds []string) []string {
	var cleanedWrds []string
	re2 := regexp2.MustCompile(fmt.Sprintf(`(?!\S*[%v])\b\S+`, exclude), 0)
	for _, wrd := range wrds {
		if isMatch, _ := re2.MatchString(wrd); isMatch {
			cleanedWrds = append(cleanedWrds, wrd)
		}
	}
	return cleanedWrds
}

// SetComboPattern returns a slice of words that match the pattern
func SetComboPattern(pattern string, words []string) []string {
	re := regexp.MustCompile(pattern)
	var newWords []string
	for _, wrd := range words {
		if re.MatchString(wrd) {
			newWords = append(newWords, wrd)
		}
	}
	return newWords
}

func Report(wrds []string, c Config) string {
	var b strings.Builder
	switch {
	case c.Exclude != "" && c.Include == "":
		return excludeOnly(b, wrds, c)
	case c.Exclude == "" && c.Include != "":
		return includeOnly(b, wrds, c)
	case c.Exclude != "" && c.Include != "":
		return includeAndExclude(b, wrds, c)
	default:
		return fmt.Sprintf("%v words found", len(wrds))
	}
}
