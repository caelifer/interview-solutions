/*
Goat Latin is a made-up language based off of English, sort of like Pig Latin.
The rules of Goat Latin are as follows:

1. If a word begins with a consonant (i.e. not a vowel), remove the first
letter and append it to the end, then add 'ma'.
For example, the word 'goat' becomes 'oatgma'.

2. If a word begins with a vowel, append 'ma' to the end of the word.
For example, the word 'I' becomes 'Ima'.

3. Add one letter "a" to the end of each word per its word index in the
sentence, starting with 1. That is, the first word gets "a" added to the
end, the second word gets "aa" added to the end, the third word in the
sentence gets "aaa" added to the end, and so on.

Write a function that, given a string of words making up one sentence, returns
that sentence in Goat Latin. For example:

        string_to_goat_latin('I speak Goat Latin')

should print: 'Imaa peaksmaaa oatGmaaaa atinLmaaaaa'
*/

package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Global splitter
var splitter = regexp.MustCompile("[ \t\v]")

func main() {
	string_to_goat_latin("I speak Goat Latin")
}

func string_to_goat_latin(s string) {
	words := splitter.Split(s, -1)
	trans := make([]string, 0, len(words))
	// Translate word by word
	for i, word := range words {
		trans = append(trans, translateWord(word, i+1))
	}
	// Print translation
	fmt.Println(strings.Join(trans, " "))
}

// Translate word
func translateWord(word string, n int) string {
	// Split word into slice of runes
	runes := []rune(word)
	// Pre-allocate maximum needed buffer on the stack
	buf := make([]rune, len(runes)+1+2+n)
	// Copy word's runes into our local buffer
	copy(buf, runes)
	// Check if the first letter is a consonant
	if !isVowel(buf[0]) {
		// Rule #1
		buf = append(buf[1:], buf[0])
	}
	// Rule #2 + end of rule #1
	buf = append(buf, 'm', 'a')
	// Rule #3
	buf = append(buf, repeatRune('a', n)...)
	// Return translated word as a string
	return string(squize(buf))
}

// Check if rune is an English vowel
func isVowel(r rune) bool {
	switch r {
	// a e i o u A E I O U
	case 'a', 'A', 'e', 'E', 'i', 'I', 'o', 'O', 'u', 'U':
		return true
	}
	return false
}

// Generate slice of N runes
func repeatRune(r rune, n int) []rune {
	res := make([]rune, n)
	for i := range res {
		res[i] = r
	}
	return res
}

func squize(runes []rune) []byte {
	res := make([]byte, 0, len(runes))
	for _, v := range []byte(string(runes)) {
		if v > 0 {
			res = append(res, v)
		}
	}
	return res
}
