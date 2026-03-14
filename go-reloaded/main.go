package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Brain of the project
func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run . input.txt output.txt")
		return
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	//  Reading file function
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}
	text := string(data)

	text = processNumbers(text)
	text = processCase(text)
	text = fixPunctuation(text)
	text = fixQuotes(text)
	text = fixArticles(text)

	// Write file function
	err = os.WriteFile(outputFile, []byte(text), 0644)
	if err != nil {
		fmt.Println("Error writing files", err)
		return
	}
}

func processNumbers(text string) string {
	words := strings.Fields(text) // "(hex) i am a man 42 (hex) old" ==> ["i","am","a","man","42","old","fellow"]

	for k := 0; k < len(words); k++ { // k = 5; k < 7; k++

		if words[k] == "(hex)" && k > 0 { // words[5] ====> "(hex)" && (k > 0 not the first word)
			num, err := strconv.ParseInt(words[k-1], 16, 64) // "42", 16 , 64 // 62

			if err == nil {
				words[k-1] = strconv.FormatInt(num, 10) // "62"
			}

			words = append(words[:k], words[k+1:]...) // words----- [start:stop]== [0:len(words)]      words[0:4] // [i -- a]   --- words[5+1:]
			k--
			continue
		}

		if words[k] == "(bin)" && k > 0 { // words[5] ====> "(hex)"
			num, err := strconv.ParseInt(words[k-1], 2, 64)
			if err == nil {
				words[k-1] = strconv.FormatInt(num, 10)
			}
			words = append(words[:k], words[k+1:]...)
			k--
			continue
		}
	}
	return strings.Join(words, " ")
}

func processCase(text string) string {
	words := strings.Fields(text)

	for k := 0; k < len(words); k++ {
		word := words[k]

		// Merge split modifiers like "(cap," + "6)"
		if (strings.HasPrefix(word, "(up,") ||
			strings.HasPrefix(word, "(low,") ||
			strings.HasPrefix(word, "(cap,")) && k+1 < len(words) {

			word = word + words[k+1]
			words = append(words[:k+1], words[k+2:]...)
		}

		//  Detect modifier (with or without number)
		if strings.HasPrefix(word, "(up") ||
			strings.HasPrefix(word, "(low") ||
			strings.HasPrefix(word, "(cap") {

			action := ""
			count := 1

			if strings.HasPrefix(word, "(up") {
				action = "up"
			} else if strings.HasPrefix(word, "(low") {
				action = "low"
			} else {
				action = "cap"
			}

			// ✅ Extract number if exists
			if strings.Contains(word, ",") {
				clean := strings.Trim(word, "()")
				parts := strings.Split(clean, ",")
				if len(parts) > 1 {
					num, err := strconv.Atoi(strings.TrimSpace(parts[1]))
					if err == nil {
						count = num
					}
				}
			}

			// ✅ Apply transformation backward
			for j := 1; j <= count && k-j >= 0; j++ {
				switch action {
				case "up":
					words[k-j] = strings.ToUpper(words[k-j])
				case "low":
					words[k-j] = strings.ToLower(words[k-j])
				case "cap":
					content := strings.ToLower(words[k-j])
					if len(content) > 0 {
						words[k-j] = strings.ToUpper(string(content[0])) + content[1:]
					}
				}
			}

			// Remove modifier
			words = append(words[:k], words[k+1:]...)
			k--
		}
	}

	return strings.Join(words, " ")
}
func fixPunctuation(text string) string {
	reEllipsis := regexp.MustCompile(`([.,!?;:])\s+([.,!?;:])`)
	for reEllipsis.MatchString(text) {
		text = reEllipsis.ReplaceAllString(text, "$1$2")
	}
	re1 := regexp.MustCompile(`\s+([.,!?;:])`)
	text = re1.ReplaceAllString(text, "$1")

	re2 := regexp.MustCompile(`([.,!?;:]+)([^\s.,!?;:])`)
	text = re2.ReplaceAllString(text, "$1 $2")

	return text
}

func fixQuotes(text string) string {
	re := regexp.MustCompile(`'\s*([^']*?)\s*'`)
	return re.ReplaceAllString(text, "'$1'")
}

func fixArticles(text string) string {
	words := strings.Fields(text)

	vowels := "aeiouhAEIOUH"

	for k := 0; k < len(words)-1; k++ {

		if strings.ToLower(words[k]) == "a" {
			next := words[k+1]
			if len(next) > 0 && strings.ContainsRune(vowels, rune(next[0])) {
				if words[k] == "A" {
					words[k] = "An"
				} else {
					words[k] = "an"
				}
			}
		}
	}

	return strings.Join(words, " ")
}
