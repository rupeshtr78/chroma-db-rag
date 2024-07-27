package main

import (
	"chroma-db/cmd/db"
	"chroma-db/internal/queryvectordb"
	"chroma-db/pkg/logger"
	"context"
)

var log = logger.Log

func main() {
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	collection, err := db.LoadDataToVectorDB(ctx, "test/Model Params.pdf")
	if err != nil {
		log.Error().Msgf("Failed to load data to vector db: %v", err)
	}

	// Query the collection
	queryTexts := []string{"what is mirostat_tau?"}
	err = queryvectordb.QueryVectorDbWithOptions(ctx, collection, queryTexts)
	if err != nil {
		log.Error().Msgf("Failed to query vector db: %v", err)
	}

	// s, err := prompts.GetTemplate(queryString, vectorResults)
	// if err != nil {
	// 	log.Error().Msgf("Failed to get template: %v", err)

	// }

	// log.Info().Msgf("Final Prompt: %s", s)

	// chat.ChatOllama(ctx, s)
}
