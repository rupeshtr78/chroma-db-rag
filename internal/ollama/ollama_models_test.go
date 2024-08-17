package ollamamodel

import (
	"chroma-db/internal/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOllamaModel(t *testing.T) {
	// Test with valid URL and model
	ollamaUrl := constants.OllamaUrl
	model := constants.OllamaChatModel
	ollamaModel, err := GetOllamaModel(ollamaUrl, model)
	assert.NoError(t, err)
	assert.NotNil(t, ollamaModel)

}

func TestGetOllamaEmbeddingFn(t *testing.T) {
	// Test with valid URL and model
	ollamaUrl := constants.OllamaUrl
	model := constants.OllamaEmbdedModel
	embeddingFn, err := GetOllamaEmbeddingFn(ollamaUrl, model)
	assert.NoError(t, err)
	assert.NotNil(t, embeddingFn)

}
