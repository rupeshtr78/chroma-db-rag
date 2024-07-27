package ollamamodel

import (
	"chroma-db/internal/constants"

	ollamaEmbedder "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
	"github.com/amikos-tech/chroma-go/types"
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

func GetOllamaEmbedder(ollamaUrl string, model string) (embeddings.Embedder, error) {

	ollamaLLM, err := GetOllamaModel(ollamaUrl, model)
	if err != nil {
		return nil, err
	}

	ollamaEmbedder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		return nil, err
	}

	return ollamaEmbedder, nil
}

func GetOllamaEmbeddingFn(ollamaUrl string, model string) (types.EmbeddingFunction, error) {

	embeddingFn, err := ollamaEmbedder.NewOllamaEmbeddingFunction(
		ollamaEmbedder.WithBaseURL(ollamaUrl),
		ollamaEmbedder.WithModel(constants.OllamaEmbdedModel))
	if err != nil {
		return nil, err
	}

	return embeddingFn, nil

}
