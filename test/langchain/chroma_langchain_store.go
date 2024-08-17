// pOC using langchain
package langchain

import (
	"context"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

// CreateChromaStore creates a new langchain **chroma.Store** with the given parameters
// Poc using langchain
func CreateChromaStore(ctx context.Context,
	chromaUrl string,
	nameSpace string,
	embedder embeddings.Embedder,
	distanceFunction types.DistanceFunction) (*chroma.Store, error) {

	store, err := chroma.New(
		chroma.WithChromaURL(chromaUrl),
		chroma.WithNameSpace(nameSpace),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction(distanceFunction), // default is cosine l2 ip
	)
	if err != nil {
		return nil, err
	}

	return &store, nil

}
