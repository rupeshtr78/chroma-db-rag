package embedders

import (
	"chroma-db/internal/constants"
	"chroma-db/pkg/logger"

	huggingface "github.com/amikos-tech/chroma-go/hf"
	ollamaEmbedder "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
	"github.com/amikos-tech/chroma-go/types"
)

type EmbeddingManager interface {
	GetEmbeddingFunction() (types.EmbeddingFunction, error)
}

func NewEmbeddingManager(embedder constants.Embedder, baseurl string, model string) EmbeddingManager {
	switch embedder {
	case constants.HuggingFace:
		return &HuggingFaceEmbedder{
			BaseUrl: baseurl,
			Model:   model,
		}
	case constants.Ollama:
		return &OllamaEmbedder{
			BaseUrl: baseurl,
			Model:   model,
		}
	default:
		return nil
	}
}

type HuggingFaceEmbedder struct {
	BaseUrl string
	Model   string
}

// GetHuggingFaceEmbedding returns a new HuggingFace Embedding Function using ami-chroma-go baseUrl and model
func (hf *HuggingFaceEmbedder) GetEmbeddingFunction() (types.EmbeddingFunction, error) {
	ef, err := huggingface.NewHuggingFaceEmbeddingInferenceFunction(hf.BaseUrl)
	if err != nil {
		logger.Log.Error().Msgf("Error getting hugging face embedding function: %v\n", err)
	}
	logger.Log.Debug().Msgf("Hugging Face Embedding Function using model: %v\n", hf.Model)
	return ef, err
}

type OllamaEmbedder struct {
	BaseUrl string
	Model   string
}

// GetOllamaEmbedding returns a new Ollama Embedding Function using ami-chroma-g
func (oe *OllamaEmbedder) GetEmbeddingFunction() (types.EmbeddingFunction, error) {
	embeddingFn, err := ollamaEmbedder.NewOllamaEmbeddingFunction(
		ollamaEmbedder.WithBaseURL(oe.BaseUrl),
		ollamaEmbedder.WithModel(oe.Model))
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().Msgf("Ollama Embedding Function using model: %v\n", oe.Model)
	return embeddingFn, nil
}
