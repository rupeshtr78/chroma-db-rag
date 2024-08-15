package reranker

import (
	"context"
	"strings"
	"sync"

	pb "chroma-db/internal/grpc/generated"
	"chroma-db/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Singleton GrpcClient
type GrpcClient struct {
	client *grpc.ClientConn
}

func (g *GrpcClient) Close() {
	g.client.Close()
}

var grpcClient *GrpcClient
var once sync.Once
var log = logger.Log

// GetGrpcRerankClient returns a singleton instance of the GrpcClient
func GetGrpcRerankClient(ctx context.Context, targetServer string) (*GrpcClient, error) {
	var err error
	once.Do(func() {
		grpcClient = &GrpcClient{}
		grpcClient.client, err = grpc.NewClient(targetServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	return grpcClient, err
}

func GrpcRerank(ctx context.Context, grpcConn *GrpcClient, targetServer string, query string, texts []string) (*pb.RerankResponse, error) {
	c := pb.NewRerankClient(grpcConn.client)

	req := &pb.RerankRequest{
		Query:      query,
		Texts:      texts,
		RawScores:  false,
		ReturnText: true,
	}

	res, err := c.Rerank(ctx, req)
	if err != nil {
		log.Error().Msgf("Failed to rerank: %v", err)
	}
	// get the response including the text field in Rank messages
	for _, rank := range res.Ranks {
		t := rank.Text
		logger.Log.Debug().Msgf("Rank: index=%v, text=%s, score=%f", rank.Index, *t, rank.Score)
	}
	logger.Log.Debug().Msgf("Rerank Metadata: %v", res.Metadata)

	return res, err
}

func (g *GrpcClient) RerankQueryResult(ctx context.Context, queryTexts []string, texts []string) (string, error) {
	queryString := strings.Builder{}
	for _, text := range queryTexts {
		queryString.WriteString(text)
	}
	response, err := GrpcRerank(ctx, g, "", queryString.String(), texts)
	if err != nil {
		return "", err
	}
	// top reranked result return the text of the first index
	firstIndex := response.Ranks[0]
	return *firstIndex.Text, nil
}
