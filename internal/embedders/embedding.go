package embedders

import (
	"chroma-db/internal/constants"
	"chroma-db/pkg/logger"
	"os"

	huggingface "github.com/amikos-tech/chroma-go/hf"
	"github.com/amikos-tech/chroma-go/openai"
	ollamaEmbedder "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
	"github.com/amikos-tech/chroma-go/types"
)

var log = logger.Log

type EmbeddingManager interface {
	GetEmbeddingFunction() (types.EmbeddingFunction, error)
}

// NewEmbeddingManager returns a new EmbeddingManager based on the embedder
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
	case constants.OpenAI:
		return &OpenAiEmbedder{
			ApiKey: os.Getenv(constants.OpenAIApiKey),
			Model:  model,
		}
	default:
		return nil
	}
}

func CreateEmbeddingFunction(embedder constants.Embedder, embeddingModel string) (types.EmbeddingFunction, error) {
	// Get Embedding either HuggingFace or Ollama
	em := NewEmbeddingManager(embedder, constants.HuggingFaceTeiUrl, embeddingModel)

	hfef, err := em.GetEmbeddingFunction()
	if err != nil {
		log.Debug().Msgf("Error getting hugging face embedding function: %v\n", err)
		return nil, err
	}

	return hfef, nil
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

type OpenAiEmbedder struct {
	ApiKey string
	Model  string
}

func (o *OpenAiEmbedder) GetEmbeddingFunction() (types.EmbeddingFunction, error) {
	// Create new OpenAI embedding function
	openaiEf, err := openai.NewOpenAIEmbeddingFunction(o.ApiKey,
		openai.WithModel(openai.EmbeddingModel(o.Model)))
	if err != nil {
		log.Error().Msgf("Error creating embedding function: %s \n", err)
		return nil, err
	}

	return openaiEf, err
}
