package main

import (
	"chroma-db/cmd/chat"
	"chroma-db/cmd/db"
	"chroma-db/internal/prompts"
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	queryString := "what is mirostat_eta"
	vectorResults, err := db.RunVectorDb(ctx, queryString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_ = vectorResults
	// chat.ChatOllama(ctx)
	// gitquery.GitCodeQuery()

	// content := `mirostat_tau Controls the balance between coherence and diversity of the output.
	// // A lower value will result in more focused and coherent text. (Default: 5.0)`
	s, err := prompts.GetTemplate(queryString, vectorResults)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(s)

	chat.ChatOllama(ctx, s)
}
