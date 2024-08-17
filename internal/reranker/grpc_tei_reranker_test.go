package reranker

import (
	"context"
	"testing"

	pb "chroma-db/internal/grpc/generated"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// MockRerankClient is a mock implementation of the RerankClient interface
type MockRerankClient struct {
	mock.Mock
}

func (m *MockRerankClient) Rerank(ctx context.Context, in *pb.RerankRequest, opts ...grpc.CallOption) (*pb.RerankResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.RerankResponse), args.Error(1)
}

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Setup code (if any)

	// Run tests
	m.Run()

	// Teardown code (if any)
}

func TestGetGrpcRerankClient(t *testing.T) {
	ctx := context.Background()
	targetServer := "localhost:50051"

	client, err := GetGrpcRerankClient(ctx, targetServer)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Ensure singleton behavior
	client2, err := GetGrpcRerankClient(ctx, targetServer)
	assert.NoError(t, err)
	assert.Equal(t, client, client2)

	client.Close()
}

func TestGrpcRerank(t *testing.T) {
	ctx := context.Background()
	query := "test query"
	texts := []string{"text1", "text2"}

	// Mock response
	mockResponse := &pb.RerankResponse{
		Ranks: []*pb.Rank{
			{Index: 0, Text: proto.String("text1"), Score: 0.9},
			{Index: 1, Text: proto.String("text2"), Score: 0.8},
		},
		Metadata: &pb.Metadata{},
	}

	// Mock client
	mockClient := new(MockRerankClient)
	mockClient.On("Rerank", ctx, &pb.RerankRequest{
		Query:      query,
		Texts:      texts,
		RawScores:  false,
		ReturnText: true,
	}).Return(mockResponse, nil)

	res, err := mockClient.Rerank(ctx, &pb.RerankRequest{
		Query:      query,
		Texts:      texts,
		RawScores:  false,
		ReturnText: true,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, mockResponse, res)

	mockClient.AssertExpectations(t)
}

func TestRerankQueryResult(t *testing.T) {
	ctx := context.Background()
	query := "test query"
	texts := []string{"text1", "text2"}

	// Mock response
	mockResponse := &pb.RerankResponse{
		Ranks: []*pb.Rank{
			{Index: 0, Text: proto.String("text1"), Score: 0.9},
		},
	}

	// Mock client
	mockClient := new(MockRerankClient)
	mockClient.On("Rerank", ctx, &pb.RerankRequest{
		Query:      query,
		Texts:      texts,
		RawScores:  false,
		ReturnText: true,
	}).Return(mockResponse, nil)

	// Replace the actual gRPC client with the mock client
	res, err := mockClient.Rerank(ctx, &pb.RerankRequest{
		Query:      query,
		Texts:      texts,
		RawScores:  false,
		ReturnText: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "text1", res.Ranks[0].GetText())

	mockClient.AssertExpectations(t)
}
