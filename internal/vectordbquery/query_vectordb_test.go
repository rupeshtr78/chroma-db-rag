package vectordbquery_test

import (
	"context"
	"errors"
	"testing"

	"chroma-db/internal/vectordbquery"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCollection is a mock implementation of chromago.Collection
type MockCollection struct {
	mock.Mock
	chromago.Collection
}

func (m *MockCollection) EmbeddingFunction() ([]*types.Embedding, error) {
	args := m.Called()
	return args.Get(0).([]*types.Embedding), args.Error(1)
}

func (m *MockCollection) QueryWithOptions(ctx context.Context, options ...types.CollectionQueryOption) (*chromago.QueryResults, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(*chromago.QueryResults), args.Error(1)
}

// MockEmbeddingFunction is a mock implementation of chromago.EmbeddingFunction
type MockEmbeddingFunction struct {
	mock.Mock
	types.EmbeddingFunction
}

func (m *MockEmbeddingFunction) EmbedDocuments(ctx context.Context, documents []string) ([]*types.Embedding, error) {
	args := m.Called(ctx, documents)
	return args.Get(0).([]*types.Embedding), args.Error(1)
}

func TestQueryVectorDbWithOptions_Success(t *testing.T) {
	ctx := context.Background()
	mockCollection := new(MockCollection)
	mockEmbeddingFunction := new(MockEmbeddingFunction)

	embedding := []*types.Embedding{{ArrayOfFloat32: &[]float32{0.1, 0.2, 0.3}}}
	queryResults := &chromago.QueryResults{
		Documents: [][]string{{"doc1", "doc2"}},
		Distances: [][]float32{{0.1, 0.2}},
		Metadatas: [][]map[string]interface{}{{{"key1": "value1"}, {"key2": "value2"}}},
	}

	mockEmbeddingFunction.On("EmbedDocuments", ctx, mock.Anything).Return(embedding, nil)
	mockCollection.On("EmbeddingFunction").Return(mockEmbeddingFunction)
	mockCollection.On("QueryWithOptions", ctx, mock.Anything).Return(queryResults, nil)

	result, err := vectordbquery.QueryVectorDbWithOptions(ctx, &mockCollection.Collection, []string{"query1"})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, queryResults, result)
}

func TestQueryVectorDbWithOptions_EmbedError(t *testing.T) {
	ctx := context.Background()
	mockCollection := new(MockCollection)
	mockEmbeddingFunction := new(MockEmbeddingFunction)

	mockEmbeddingFunction.On("EmbedDocuments", ctx, mock.Anything).Return(nil, errors.New("embedding error"))
	mockCollection.On("EmbeddingFunction").Return(mockEmbeddingFunction)

	result, err := vectordbquery.QueryVectorDbWithOptions(ctx, &mockCollection.Collection, []string{"query1"})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "embedding error", err.Error())
}

func TestQueryVectorDbWithOptions_QueryError(t *testing.T) {
	ctx := context.Background()
	mockCollection := new(MockCollection)
	mockEmbeddingFunction := new(MockEmbeddingFunction)

	embedding := []*types.Embedding{{ArrayOfFloat32: &[]float32{0.1, 0.2, 0.3}}}

	mockEmbeddingFunction.On("EmbedDocuments", ctx, mock.Anything).Return(embedding, nil)
	mockCollection.On("EmbeddingFunction").Return(mockEmbeddingFunction)
	mockCollection.On("QueryWithOptions", ctx, mock.Anything).Return(nil, errors.New("query error"))

	result, err := vectordbquery.QueryVectorDbWithOptions(ctx, &mockCollection.Collection, []string{"query1"})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "query error", err.Error())
}
