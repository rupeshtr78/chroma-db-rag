package vectordb

import (
	"chroma-db/internal/constants"
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/rs/zerolog/log"
)

type RecordSet interface {
	WithRecord(recordOpts ...types.Option) *types.RecordSet
	BuildAndValidate(ctx context.Context) ([]*types.Record, error)
	AddTextToRecordSet(ctx context.Context, documents []string, metadata map[string]any) (*types.RecordSet, error)
}

type ChromagoRecordSet struct {
	RecordSet *types.RecordSet
}

func (rsw *ChromagoRecordSet) WithRecord(recordOpts ...types.Option) *types.RecordSet {
	return rsw.RecordSet.WithRecord(recordOpts...)
}

func (rsw *ChromagoRecordSet) BuildAndValidate(ctx context.Context) ([]*types.Record, error) {
	return rsw.RecordSet.BuildAndValidate(ctx)
}

func CreateRecordSet(embeddingFunction types.EmbeddingFunction) (*ChromagoRecordSet, error) {
	// Create a new record set with to hold the records to insert
	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(embeddingFunction),
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Err(err).Msg("Error creating record set")
		return nil, err
	}

	return &ChromagoRecordSet{RecordSet: rs}, nil
}

// AddRecordSetToCollection adds a record set to a collection
func (c *ChromagoCollection) AddRecordSetToCollection(ctx context.Context,
	recordSet *ChromagoRecordSet,
	docs []string,
	metadata constants.Metadata) (*chromago.Collection, error) {

	records, err := recordSet.AddTextToRecordSet(ctx, docs, metadata)
	if err != nil {
		log.Debug().Msgf("Error adding to record set: %v\n", err)
		return nil, err
	}

	_, err = recordSet.BuildAndValidate(ctx)
	if err != nil {
		log.Debug().Msgf("Error building and validating records: %v\n", err)
		return nil, err
	}

	coll, err := c.AddRecords(ctx, records)
	if err != nil {
		log.Err(err).Msg("Error adding records to collection")
		return nil, err
	}

	log.Debug().Msgf("Added %d records to collection", len(docs))

	return coll, nil
}

// internal/chromaclient/chroma_recordset.go
func (rs *ChromagoRecordSet) AddTextToRecordSet(ctx context.Context, documents []string,
	metadata map[string]any) (*types.RecordSet, error) {

	// Iterate over documents and metadata list and add records to the record set
	for _, doc := range documents {
		rs.WithRecord(
			types.WithDocument(doc),
			// types.WithMetadatas(metadata),
		)
	}

	return rs.RecordSet, nil
}

// AddPdfToRecordSet adds pdf documents to the record set
func AddPdfToRecordSet(ctx context.Context,
	collection *chromago.Collection,
	rs *types.RecordSet,
	documents []string,
	metadata map[string]any) (*types.RecordSet, error) {

	// Iterate over documents and metadata list and add records to the record set
	for i, doc := range documents {
		pageNum := i + 1
		key := fmt.Sprintf("%d", pageNum)
		metadataValue := metadata[key].(string)
		rs.WithRecord(
			types.WithDocument(doc),
			types.WithMetadata(key, metadataValue),
		)
	}

	return rs, nil
}
