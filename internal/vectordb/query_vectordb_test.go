package vectordb

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

type MockEmbeddingFunc struct {
	mock.Mock
}

func (m *MockEmbeddingFunc) EmbedDocuments(ctx context.Context, docs []string) ([]*types.Embedding, error) {
	args := m.Called(ctx, docs)
	return args.Get(0).([]*types.Embedding), args.Error(1)
}

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) QueryWithOptions(ctx context.Context, options ...types.CollectionQueryOption) (*chromago.QueryResults, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(*chromago.QueryResults), args.Error(1)
}

func (m *MockCollection) EmbeddingFunction() EmbeddingFunc {
	args := m.Called()
	return args.Get(0).(EmbeddingFunc)
}

func TestEmbedQuery(t *testing.T) {
	ctx := context.Background()
	mockEmbeddingFunc := new(MockEmbeddingFunc)

	query := []string{"example query"}
	embeddings := []*types.Embedding{
		{ArrayOfFloat32: &[]float32{0.1, 0.2, 0.3}},
	}

	mockEmbeddingFunc.On("EmbedDocuments", ctx, query).Return(embeddings, nil)

	result, err := EmbedQuery(ctx, mockEmbeddingFunc, query)
	assert.NoError(t, err)
	assert.Equal(t, embeddings, result)

	mockEmbeddingFunc.AssertExpectations(t)
}

func TestQueryVectorDbWithOptions(t *testing.T) {
	ctx := context.Background()
	mockCollection := new(MockCollection)
	mockEmbeddingFunc := new(MockEmbeddingFunc)

	queryTexts := []string{"query1", "query2"}
	embeddings := []*types.Embedding{
		{ArrayOfFloat32: &[]float32{0.1, 0.2, 0.3}},
	}

	mockEmbeddingFunc.On("EmbedDocuments", ctx, queryTexts).Return(embeddings, nil)

	results := &chromago.QueryResults{
		Documents: [][]string{
			{"doc1", "doc2"},
		},
	}

	mockCollection.On("EmbeddingFunction").Return(mockEmbeddingFunc)
	mockCollection.On("QueryWithOptions", ctx, mock.Anything).Return(results, nil)

	result, err := QueryVectorDbWithOptions(ctx, mockCollection, queryTexts)
	assert.NoError(t, err)
	assert.Equal(t, results, result)

	mockCollection.AssertExpectations(t)
	mockEmbeddingFunc.AssertExpectations(t)
}

func TestQueryVectorDbWithOptions_EmbeddingError(t *testing.T) {
	ctx := context.Background()
	mockCollection := new(MockCollection)
	mockEmbeddingFunc := new(MockEmbeddingFunc)

	queryTexts := []string{"query1", "query2"}
	embeddingError := errors.New("embedding error")

	mockEmbeddingFunc.On("EmbedDocuments", ctx, queryTexts).Return([]*types.Embedding{}, embeddingError)

	mockCollection.On("EmbeddingFunction").Return(mockEmbeddingFunc)

	result, err := QueryVectorDbWithOptions(ctx, mockCollection, queryTexts)
	assert.Nil(t, result)
	assert.Equal(t, embeddingError, err)

	mockCollection.AssertExpectations(t)
	mockEmbeddingFunc.AssertExpectations(t)
}

func TestQueryVectorDbWithOptions_QueryError(t *testing.T) {
	ctx := context.Background()
	mockCollection := new(MockCollection)
	mockEmbeddingFunc := new(MockEmbeddingFunc)

	queryError := errors.New("query error")
	queryTexts := []string{"query1", "query2"}
	embeddings := []*types.Embedding{
		{ArrayOfFloat32: &[]float32{0.1, 0.2, 0.3}},
	}
	mockEmbeddingFunc.On("EmbedDocuments", ctx, queryTexts).Return(embeddings, nil)

	mockCollection.On("EmbeddingFunction").Return(mockEmbeddingFunc)
	mockCollection.On("QueryWithOptions", ctx, mock.Anything).Return(&chromago.QueryResults{}, queryError)

	result, err := QueryVectorDbWithOptions(ctx, mockCollection, queryTexts)
	assert.Nil(t, result)
	assert.Equal(t, queryError, err)

	mockCollection.AssertExpectations(t)
	mockEmbeddingFunc.AssertExpectations(t)

}
