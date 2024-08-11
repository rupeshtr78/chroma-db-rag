package utils

import (
	"fmt"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/rs/zerolog/log"
)

// stripStopWords removes the stop words from the text and returns a slice of words
func stripStopWords(text string) []string {
	langCode := "en"

	// remove stopwords
	cleanContent := stopwords.CleanString(text, langCode, true)
	fmt.Println(cleanContent)

	// convert to slice of words
	result := make([]string, 0)
	// split the text into words and trim the spaces
	for _, word := range strings.Split(cleanContent, " ") {
		trimmedWord := strings.TrimSpace(word)
		// remove extra spaces
		if len(trimmedWord) > 0 {
			result = append(result, trimmedWord)
		}
	}

	log.Debug().Msgf("Vector Query Strings: %v", result)

	return result
}
