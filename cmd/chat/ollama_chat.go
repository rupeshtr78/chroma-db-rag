package chat

import (
	"chroma-db/internal/constants"
	ollamamodel "chroma-db/internal/ollama"
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
)

func ChatOllama(ctx context.Context, prompt string) {

	l, err3 := ollamamodel.GetOllamaModel(constants.OllamaUrl, constants.OllamaChatModel)
	if err3 != nil {
		log.Default().Println(err3)
	}

	// prompt := "Why is Sky Blue?"

	_, err4 := l.Call(ctx, prompt,
		llms.WithMaxTokens(1024),
		llms.WithStreamingFunc(handleStreamingFunc),
		llms.WithSeed(42),
		llms.WithTemperature(0.5), // 0.5 0.9
		llms.WithTopP(0.9),

		// llms.WithTopK(40),
	)
	if err4 != nil {
		log.Default().Println(err4)
	}
	// log.Default().Println(s)
}

func handleStreamingFunc(ctx context.Context, chunk []byte) error {
	if len(chunk) == 0 {
		return nil
	}
	fmt.Print(string(chunk))
	return nil
}
