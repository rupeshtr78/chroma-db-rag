package reranker

import (
	"context"
	"log"

	pb "chroma-db/internal/grpc/generated"
	"chroma-db/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcRerank(ctx context.Context, targetServer string, query string, texts []string) (*pb.RerankResponse, error) {
	conn, err := grpc.NewClient(targetServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewRerankClient(conn)

	req := &pb.RerankRequest{
		Query:      query,
		Texts:      texts,
		RawScores:  false,
		ReturnText: true,
	}

	res, err := c.Rerank(ctx, req)
	if err != nil {
		log.Fatalf("error while calling Rerank: %v", err)
	}
	// Print the response including the text field in Rank messages
	for _, rank := range res.Ranks {
		t := rank.Text
		logger.Log.Debug().Msgf("Rank: index=%v, text=%s, score=%f", rank.Index, *t, rank.Score)
	}
	logger.Log.Debug().Msgf("Rerank Metadata: %v", res.Metadata)

	return res, nil
}
