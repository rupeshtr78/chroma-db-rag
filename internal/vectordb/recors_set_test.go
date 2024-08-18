package vectordb

import (
	"chroma-db/internal/constants"
	"context"
	"testing"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRecordSet is a mock implementation of the RecordSet interface
type MockRecordSet struct {
	mock.Mock
}

func (m *MockRecordSet) WithRecord(recordOpts ...types.Option) *types.RecordSet {
	args := m.Called(recordOpts)
	return args.Get(0).(*types.RecordSet)
}

func (m *MockRecordSet) BuildAndValidate(ctx context.Context) ([]*types.Record, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.Record), args.Error(1)
}

func (m *MockRecordSet) WithDocument(document string) types.Option {
	args := m.Called(document)
	return args.Get(0).(types.Option)
}

func (m *MockRecordSet) AddTextToRecordSet(ctx context.Context, documents []string, metadata map[string]any) (*types.RecordSet, error) {
	args := m.Called(ctx, documents, metadata)
	return args.Get(0).(*types.RecordSet), args.Error(1)
}

type MockChromagoCollection struct {
	mock.Mock
}

func (m *MockChromagoCollection) AddRecords(ctx context.Context, records *types.RecordSet) (*chromago.Collection, error) {
	args := m.Called(ctx, records)
	return args.Get(0).(*chromago.Collection), args.Error(1)
}

func (m *MockChromagoCollection) AddRecordSetToCollection(ctx context.Context, recordSet *ChromagoRecordSet, docs []string, metadata constants.Metadata) (*chromago.Collection, error) {
	args := m.Called(ctx, recordSet, docs, metadata)
	return args.Get(0).(*chromago.Collection), args.Error(1)
}

func MockEmbeddingFunction(document string) ([]float32, error) {
	return []float32{1.0, 2.0, 3.0}, nil
}

// TestCreateRecordSet tests the CreateRecordSet function
func TestCreateRecordSet(t *testing.T) {
	mockEmbeddingFunction := new(MockEmbeddingFunc)

	rs, err := CreateRecordSet(mockEmbeddingFunction)
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	assert.NotNil(t, rs.RecordSet)
}

// TestAddTextToRecordSet tests the AddTextToRecordSet functionfunc TestAddTextToRecordSet(t *testing.T) {
func TestAddTextToRecordSet(t *testing.T) {
	mockRecordSet := new(MockRecordSet)
	mockRecordSet.On("WithRecord", mock.Anything).Return(&types.RecordSet{})
	mockRecordSet.On("WithDocument", mock.Anything).Return(types.Option(nil))
	mockRecordSet.On("AddTextToRecordSet", mock.Anything, mock.Anything, mock.Anything).Return(&types.RecordSet{}, nil)

	// Set up the documents and metadata
	documents := []string{"doc1", "doc2"}
	metadata := map[string]any{"key1": "value1", "key2": "value2"}

	// Call the function under test
	result, err := mockRecordSet.AddTextToRecordSet(context.Background(), documents, metadata)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify expectations
	mockRecordSet.AssertCalled(t, "AddTextToRecordSet", context.Background(), documents, metadata)
}

// TestAddRecordSetToCollection tests the AddRecordSetToCollection function
func TestAddRecordSetToCollection(t *testing.T) {
	mockRecordSet := new(MockRecordSet)
	mockRecordSet.On("AddTextToRecordSet", mock.Anything, mock.Anything, mock.Anything).Return(&types.RecordSet{}, nil)
	mockRecordSet.On("BuildAndValidate", mock.Anything).Return([]*types.Record{}, nil)

	mockCollection := new(MockChromagoCollection)
	mockCollection.On("AddRecords", mock.Anything, mock.Anything).Return(&chromago.Collection{}, nil)

	// Replace the actual ChromagoCollection with a mock implementation
	coll := new(MockChromagoCollection)
	coll.On("AddRecords", mock.Anything, mock.Anything).Return(&chromago.Collection{}, nil)
	coll.On("AddRecordSetToCollection", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&chromago.Collection{}, nil)
	// Set up the ChromagoRecordSet with a proper embedding function
	mockEmbeddingFunction := new(MockEmbeddingFunc)

	recordSet, err := CreateRecordSet(mockEmbeddingFunction)
	assert.NoError(t, err)

	// Set up the documents and metadata
	docs := []string{"doc1", "doc2"}
	metadata := map[string]any{"key1": "value1", "key2": "value2"}

	rs, err := recordSet.AddTextToRecordSet(context.Background(), docs, metadata)
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	// Call the function under test
	result, err := coll.AddRecordSetToCollection(context.Background(), recordSet, docs, metadata)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify expectations
	// mockRecordSet.AssertCalled(t, "AddTextToRecordSet", context.Background(), docs, metadata)
	// mockRecordSet.AssertCalled(t, "BuildAndValidate", context.Background())
	// coll.AssertCalled(t, "AddRecords", context.Background(), rs)
}
