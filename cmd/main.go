package main

import (
	"chroma-db/internal/prompts"
	"chroma-db/pkg/logger"
	"context"
)

var log = logger.Log

func main() {
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	queryString := "what is mirostat_eta"
	// vectorResults, err := db.RunVectorDb(ctx, queryString)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// _ = vectorResults
	// chat.ChatOllama(ctx)
	// gitquery.GitCodeQuery()

	vectorResults := `mirostat_tau Controls the balance between coherence and diversity of the output.
	// A lower value will result in more focused and coherent text. (Default: 5.0)`
	s, err := prompts.GetTemplate(queryString, vectorResults)
	if err != nil {
		log.Error().Msgf("Failed to get template: %v", err)

	}

	log.Info().Msgf("Final Prompt: %s", s)

	// chat.ChatOllama(ctx, s)
}
