package chat

import (
	"chroma-db/internal/constants"
	ollamamodel "chroma-db/internal/ollama"
	"chroma-db/pkg/logger"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

var log = logger.Log

type ChatManager interface {
	Chat(ctx context.Context, prompt string)
}

func NewChatManager(provider constants.LLMProvider, model string, url string, apiKey string) ChatManager {
	switch provider {
	case constants.HuggingFaceChat:
		return &HuggingFaceChat{
			Model: model,
			Url:   url,
		}
	case constants.OllamaChat:
		return &OllamaChat{
			Model: model,
			Url:   url,
		}
	case constants.OpenAIChat:
		return &OpenAiChat{
			Model:  model,
			ApiKey: apiKey,
		}
	default:
		return nil
	}
}

type OllamaChat struct {
	Model string
	Url   string
}

func (o *OllamaChat) Chat(ctx context.Context, prompt string) {

	l, err := ollamamodel.GetOllamaModel(o.Url, o.Model)
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

type OpenAiChat struct {
	Model  string
	ApiKey string
}

func (o *OpenAiChat) Chat(ctx context.Context, prompt string) {

	openaiLLm, err := openai.New(
		openai.WithToken(o.ApiKey),
		openai.WithModel(o.Model),
	)
	if err != nil {
		log.Err(err).Msg("Failed to get OpenAI model")
	}

	_, err = openaiLLm.Call(ctx, prompt,
		llms.WithMaxTokens(1024),
		llms.WithStreamingFunc(handleStreamingFunc),
		llms.WithSeed(42),
		llms.WithTemperature(0.5), // 0.5 0.9
		llms.WithTopP(0.9),
	)

	if err != nil {
		log.Err(err).Msg("Failed to call OpenAI model")
	}

}

// TODO: Implement HuggingFaceChat
type HuggingFaceChat struct {
	Model string
	Url   string
}

func (h *HuggingFaceChat) Chat(ctx context.Context, prompt string) {
	fmt.Println("TODO HuggingFaceChat")
}

// handleStreamingFunc is a callback function that is called for each chunk of data returned by the model
func handleStreamingFunc(ctx context.Context, chunk []byte) error {
	if len(chunk) == 0 {
		return nil
	}
	fmt.Print(string(chunk))
	return nil
}
