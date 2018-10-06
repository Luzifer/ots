package rconfig

import "strings"

type characterClass [2]rune

func (c characterClass) Contains(r rune) bool {
	return c[0] <= r && c[1] >= r
}

type characterClasses []characterClass

func (c characterClasses) Contains(r rune) bool {
	for _, cc := range c {
		if cc.Contains(r) {
			return true
		}
	}
	return false
}

var (
	charGroupUpperLetter = characterClass{'A', 'Z'}
	charGroupLowerLetter = characterClass{'a', 'z'}
	charGroupNumber      = characterClass{'0', '9'}
	charGroupLowerNumber = characterClasses{charGroupLowerLetter, charGroupNumber}
)

func deriveEnvVarName(s string) string {
	var (
		words []string
		word  []rune
	)

	for _, l := range s {
		switch {
		case charGroupUpperLetter.Contains(l):
			if len(word) > 0 && charGroupLowerNumber.Contains(word[len(word)-1]) {
				words = append(words, string(word))
				word = []rune{}
			}
			word = append(word, l)

		case charGroupLowerLetter.Contains(l):
			if len(word) > 1 && charGroupUpperLetter.Contains(word[len(word)-1]) {
				words = append(words, string(word[0:len(word)-1]))
				word = word[len(word)-1:]
			}
			word = append(word, l)

		case charGroupNumber.Contains(l):
			word = append(word, l)

		default:
			if len(word) > 0 {
				words = append(words, string(word))
			}
			word = []rune{}
		}
	}
	words = append(words, string(word))

	return strings.ToUpper(strings.Join(words, "_"))
}
