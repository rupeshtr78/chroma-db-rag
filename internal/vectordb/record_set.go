package vectordb

import (
	"github.com/amikos-tech/chroma-go/types"
	"github.com/rs/zerolog/log"
)

func CreateRecordSet(embeddingFunction types.EmbeddingFunction) (*types.RecordSet, error) {
	// Create a new record set with to hold the records to insert
	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(embeddingFunction),
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Err(err).Msg("Error creating record set")
		return nil, err
	}

	return rs, nil
}
