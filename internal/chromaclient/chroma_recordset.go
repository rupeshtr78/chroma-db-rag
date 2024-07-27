package chromaclient

import (
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
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

// internal/chromaclient/chroma_recordset.go
func AddToRecordSet(ctx context.Context,
	collection *chromago.Collection,
	rs *types.RecordSet,
	documents []string,
	metadata map[string]any) (*types.RecordSet, error) {

	// Iterate over documents and metadata list and add records to the record set
	for _, doc := range documents {
		rs.WithRecord(types.WithDocument(doc),
			types.WithMetadatas(metadata),
		)
	}

	// Add the records to the collection
	_, err := collection.AddRecords(ctx, rs)
	if err != nil {
		log.Err(err).Msg("Error adding records")
		return nil, err
	}
	return rs, nil
}

func QueryRecords(ctx context.Context, collection *chromago.Collection, query []string) error {
	// Query the collection
	qr, qerr := collection.Query(ctx,
		query,
		5,
		nil,
		nil,
		nil)

	if qerr != nil {
		log.Err(qerr).Msg("Error querying collection")
		return qerr
	}

	fmt.Printf("qr: %v\n", qr.Documents[0][0]) //this should result in the document about dogs
	return nil
}

// AddRecords adds records to the record set and collection
// TODO fix document and metadata
// func AddRecords(ctx context.Context, rs *types.RecordSet, collection *chromago.Collection) error {
// 	// Add a few records to the record set
// 	rs.WithRecord(types.WithDocument("My name is John. And I have two dogs."), types.WithMetadata("key1", "value1"))
// 	rs.WithRecord(types.WithDocument("My name is Jane. I am a data scientist."), types.WithMetadata("key2", "value2"))

// 	// Add the records to the collection
// 	_, err := collection.AddRecords(context.Background(), rs)
// 	if err != nil {
// 		log.Err(err).Msg("Error adding records")
// 		return err
// 	}
// 	return err
// }
