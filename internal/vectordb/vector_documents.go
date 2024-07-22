package vectordb

import (
	"context"
	"errors"
	"log"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

// Add documents to the vector store.
func AddDocuments(ctx context.Context,
	store *chroma.Store,
	documents []schema.Document,
	nameSpace string,
	embedder embeddings.Embedder) error {

	// TODO - Refector this to use the vectorstores.Option
	vecOptions := make([]vectorstores.Option, 3)
	vecOptions = append(vecOptions, vectorstores.WithEmbedder(embedder))
	vecOptions = append(vecOptions, vectorstores.WithNameSpace(nameSpace))
	// vecOptions = append(vecOptions, vectorstores.WithDeduplicater(fn func(ctx context.Context, doc schema.Document) bool)

	// 	// Add documents to the vector store. returns the ids of the added documents.
	docIds, errAd := store.AddDocuments(ctx,
		documents,
		vecOptions...,
	)
	if errAd != nil {
		log.Default().Printf("add documents: %v\n", errAd)
		return chroma.ErrAddDocument
	}
	if len(docIds) != len(documents) {
		log.Default().Printf("add documents: expected %d ids, got %d\n", len(documents), len(docIds))
		return chroma.ErrAddDocument
	}

	return nil

}

// Search for documents in the vector store.
func SearchVectorDb(ctx context.Context,
	store *chroma.Store,
	query string,
	numDocuments int,
	namespace string,
	options ...vectorstores.Option) ([]schema.Document, error) {

	// vecOpts := make([]vectorstores.Option, 2)
	// vecOpts = append(vecOpts, vectorstores.WithNameSpace(namespace))
	// vecOpts = append(vecOpts, vectorstores.WithScoreThreshold(constants.ScoreThreshold))
	// vecOpts = append(vecOpts, vectorstores.WithDeduplicater()))

	// 	// Search for similar documents in the vector store.
	// 	// returns the most similar documents to the query.
	similarDocs, errSs := store.SimilaritySearch(ctx,
		query,
		numDocuments,
		options...)
	if errSs != nil {
		log.Default().Printf("similarity search: %v\n", errSs)
		return nil, errSs
	}
	if len(similarDocs) == 0 {
		log.Default().Printf("similarity search: no similar documents found\n")
		return nil, errors.New("no similar documents found")
	}

	return similarDocs, nil

}

// Delete collection and all documents in the vector store.
// each store has only one collection.?
func DeleteCollection(ctx context.Context,
	store *chroma.Store) error {
	// Delete collection and all documents in the vector store.
	errDc := store.RemoveCollection()
	if errDc != nil {
		log.Default().Printf("delete collection: %v\n", errDc)
		return errors.New("delete collection failed")
	}

	return nil

}

// func(ctx context.Context, doc schema.Document) bool

// func dupplicateChecker(ctx context.Context, doc schema.Document) bool {
// 	// TODO - Implement the deduplicater function
// 	// Implement deduplication logic here
// 	// check if a document with the same ID already exists in the store
// 	// existingDocs, err := store.GetDocuments(ctx, []string{doc.ID})
// 	// if err != nil || len(existingDocs) > 0 {
// 	// 	return true
// 	// }
// 	return false
// }
