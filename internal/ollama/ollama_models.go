package ollamamodel

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

func GetOllamaModel(ollamaUrl string, model string) (*ollama.LLM, error) {

	ollamaModel, err := ollama.New(
		ollama.WithServerURL(ollamaUrl),
		ollama.WithModel(model),
	)
	if err != nil {
		return nil, err
	}
	return ollamaModel, nil

}

func GetOllamaEmbedding(ollamaUrl string, model string) (embeddings.Embedder, error) {

	ollamaLLM, err := GetOllamaModel(ollamaUrl, model)
	if err != nil {
		return nil, err
	}

	ollamaEmbeder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		return nil, err
	}

	return ollamaEmbeder, nil
}
