package reranker

import (
	"context"
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

var grpcClient *GrpcClient
var once sync.Once
var log = logger.Log

// GetGrpcClientInstance returns a singleton instance of the GrpcClient
func GetGrpcClientInstance(ctx context.Context, targetServer string) (*GrpcClient, error) {
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

	return res, nil
}
