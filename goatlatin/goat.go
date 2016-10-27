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

  would return: 'Imaa peaksmaaa oatGmaaaa atinLmaaaaa'
*/
package main

import (
	"fmt"
	"regexp"
)

// Global splitter
var splitter = regexp.MustCompile("[ \t\v]")

func main() {
	string_to_goat_latin("I speak Goat Latin")
}

func string_to_goat_latin(s string) {
	cnt := 0 // word counter
	for _, word := range splitter.Split(s, -1) {
		// Increment our word counter
		cnt++
		// Split word into []rune
		w := []rune(word)

		// Check first letter
		if !isVowel(w[0]) {
			w = append(w[1:], w[0])
		}

		// Rule #2 + end of rule #1
		w = append(w, 'm', 'a')

		// Rule #3
		w = append(w, repeatRune('a', cnt)...)

		fmt.Printf("%s ", string(w))
	}

	// Add new line if we process at least one word
	if cnt > 0 {
		fmt.Print("\n")
	}
}

// Check if rune is an English vowel
func isVowel(r rune) bool {
	switch r {
	// a e i o u A E I O U
	case 'a', 'A', 'e', 'E', 'i', 'I', 'o', 'O', 'u', 'U':
		return true
	default:
		return false
	}
}

// Generate string with c runes
func repeatRune(r rune, mult int) []rune {
	res := make([]rune, mult)
	for i := range res {
		res[i] = r
	}
	return res
}
