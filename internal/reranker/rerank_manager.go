package reranker

import (
	"chroma-db/internal/constants"
	"context"
	"net/http"
	"time"
)

type Reranker interface {
	RerankQueryResult(context.Context, []string, []string) (string, error)
}

func NewReRanker(ctx context.Context, protocol constants.Protocol) (Reranker, error) {
	if protocol == constants.GRPC {
		client, err := GetGrpcRerankClient(ctx, constants.GrpcRerankServer)
		if err != nil {
			return nil, err
		}
		return client, nil
	} else if protocol == constants.HTTP {
		c := NewHttpClient()
		return GetHttpRerankClient(c, constants.HuggingFaceRerankUrl, constants.HuggingFaceRerankModel, "", nil), nil
	}

	return nil, nil

}

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second, // Set a timeout
	}

}
