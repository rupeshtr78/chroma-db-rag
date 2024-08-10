package ollamamodel

import (
	"chroma-db/internal/constants"

	ollamaEmbedder "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/llms/ollama"
)

// GetOllamaModel returns a new Ollama LLM model using langchain-go
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

// GetOllamaEmbeddingFn returns a new Ollama Embedding Function using ami-chroma-go
// TODO abstract out to support hf embedding function delete after refactor
func GetOllamaEmbeddingFn(ollamaUrl string, model string) (types.EmbeddingFunction, error) {

	embeddingFn, err := ollamaEmbedder.NewOllamaEmbeddingFunction(
		ollamaEmbedder.WithBaseURL(ollamaUrl),
		ollamaEmbedder.WithModel(constants.OllamaEmbdedModel))
	if err != nil {
		return nil, err
	}

	return embeddingFn, nil

}
