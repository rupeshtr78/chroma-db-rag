package embedders

// import (
// 	"chroma-db/internal/constants"

// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockEmbeddingFunction is a mock implementation of the EmbeddingFunction interface
// type MockEmbeddingFunction struct {
// 	mock.Mock
// }

// func (m *MockEmbeddingFunction) EmbedDocuments(texts []string) ([][]float32, error) {
// 	args := m.Called(texts)
// 	return args.Get(0).([][]float32), args.Error(1)
// }

// func (m *MockEmbeddingFunction) EmbedQuery(text string) ([]float32, error) {
// 	args := m.Called(text)
// 	return args.Get(0).([]float32), args.Error(1)
// }

// func TestNewEmbeddingManager(t *testing.T) {
// 	apiKey := os.Getenv(constants.OpenAIApiKey)
// 	tests := []struct {
// 		name     string
// 		embedder constants.Embedder
// 		baseurl  string
// 		model    string
// 		expected EmbeddingManager
// 	}{
// 		{
// 			name:     "HuggingFace Embedder",
// 			embedder: constants.HuggingFace,
// 			baseurl:  "http://hf.base.url",
// 			model:    "model1",
// 			expected: &HuggingFaceEmbedder{BaseUrl: "http://hf.base.url", Model: "model1"},
// 		},
// 		{
// 			name:     "Ollama Embedder",
// 			embedder: constants.Ollama,
// 			baseurl:  "http://ollama.base.url",
// 			model:    "model2",
// 			expected: &OllamaEmbedder{BaseUrl: "http://ollama.base.url", Model: "model2"},
// 		},
// 		{
// 			name:     "OpenAI Embedder",
// 			embedder: constants.OpenAI,
// 			baseurl:  "",
// 			model:    "model3",
// 			expected: &OpenAiEmbedder{ApiKey: apiKey, Model: "model3"},
// 		},
// 		{
// 			name:     "Unknown Embedder",
// 			embedder: 100,
// 			baseurl:  "",
// 			model:    "",
// 			expected: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result := NewEmbeddingManager(tt.embedder, tt.baseurl, tt.model)
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }

// func TestHuggingFaceEmbedder_GetEmbeddingFunction(t *testing.T) {
// 	hf := &HuggingFaceEmbedder{
// 		BaseUrl: "http://hf.base.url",
// 		Model:   "model1",
// 	}

// 	ef, err := hf.GetEmbeddingFunction()
// 	assert.NoError(t, err)
// 	assert.NotNil(t, ef)
// }

// func TestOllamaEmbedder_GetEmbeddingFunction(t *testing.T) {
// 	oe := &OllamaEmbedder{
// 		BaseUrl: "http://ollama.base.url",
// 		Model:   "model2",
// 	}

// 	ef, err := oe.GetEmbeddingFunction()
// 	assert.NoError(t, err)
// 	assert.NotNil(t, ef)
// }

// func TestOpenAiEmbedder_GetEmbeddingFunction_Error(t *testing.T) {
// 	os.Unsetenv(constants.OpenAIApiKey)

// 	o := &OpenAiEmbedder{
// 		ApiKey: "",
// 		Model:  "model3",
// 	}

// 	ef, err := o.GetEmbeddingFunction()
// 	assert.Error(t, err)
// 	assert.Nil(t, ef)
// }

// // MockLogger is a mock implementation of the logger interface
// type MockLogger struct {
// 	mock.Mock
// }

// func (m *MockLogger) Debug() *log.Logger {
// 	args := m.Called()
// 	return args.Get(0).(*log.Logger)
// }

// func (m *MockLogger) Error() *log.Logger {
// 	args := m.Called()
// 	return args.Get(0).(*log.Logger)
// }

// func TestHuggingFaceEmbedder_GetEmbeddingFunction_Error(t *testing.T) {
// 	hf := &HuggingFaceEmbedder{
// 		BaseUrl: "invalid-url",
// 		Model:   "model1",
// 	}

// 	_, err := hf.GetEmbeddingFunction()
// 	assert.Error(t, err)
// 	// assert.Nil(t, ef)
// }

// func TestOllamaEmbedder_GetEmbeddingFunction_Error(t *testing.T) {
// 	oe := &OllamaEmbedder{
// 		BaseUrl: "invalid-url",
// 		Model:   "model2",
// 	}

// 	ef, err := oe.GetEmbeddingFunction()
// 	assert.Error(t, err)
// 	assert.Nil(t, ef)
// }
