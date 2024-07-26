package chat

import (
	"chroma-db/internal/constants"
	ollamamodel "chroma-db/internal/ollama"
	"chroma-db/pkg/logger"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

var log = logger.Log

func ChatOllama(ctx context.Context, prompt string) {

	l, err := ollamamodel.GetOllamaModel(constants.OllamaUrl, constants.OllamaChatModel)
	if err != nil {
		log.Err(err).Msg("Failed to get Ollama model")
	}

	_, err = l.Call(ctx, prompt,
		llms.WithMaxTokens(1024),
		llms.WithStreamingFunc(handleStreamingFunc),
		llms.WithSeed(42),
		llms.WithTemperature(0.5), // 0.5 0.9
		llms.WithTopP(0.9),

		// llms.WithTopK(40),
	)
	if err != nil {
		log.Err(err).Msg("Failed to call Ollama model")
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
