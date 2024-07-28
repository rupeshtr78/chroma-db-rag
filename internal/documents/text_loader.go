package documents

import (
	"chroma-db/internal/constants"
	"context"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func TextLoader(file string) ([]string, constants.Metadata, error) {

	f, err := os.Open(file)
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	defer f.Close()

	// Read the first byte to check if the file is empty
	i, err := f.Read([]byte{0})
	if err != nil || i == 0 {
		return nil, nil, err
	}

	// read the file into a byte slice
	fBuf := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		fBuf = append(fBuf, buf[:n]...)
	}

	// Convert the byte slice to a slice of strings
	str := strings.TrimSpace(string(fBuf))

	fileInfo, err := os.Stat(file)
	if err != nil {
		log.Error().Err(err).Msg("Error getting file info")
		return nil, nil, err
	}
	if fileInfo.Size() == 0 {
		log.Warn().Msg("Empty file")
		return nil, nil, nil
	}

	metaData := map[string]string{
		"file": fileInfo.Name(),
	}

	// Convert metaData to Metadata
	meta := make(constants.Metadata)
	for k, v := range metaData {
		meta[k] = v
	}
	return []string{str}, meta, nil

}

func TextLoaderV2(ctx context.Context, file string) ([]string, constants.Metadata, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	loader := documentloaders.NewText(f)
	// docs, err := loader.Load(ctx)

	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(1024),
		textsplitter.WithChunkOverlap(512),
	)
	docs, err := loader.LoadAndSplit(ctx, splitter)

	if err != nil || len(docs) == 0 {
		return nil, nil, err
	}

	// returning a slice of strings and a Metadata map
	strSlice := make([]string, 0)
	metadata := make(constants.Metadata)
	for _, doc := range docs {
		strSlice = append(strSlice, doc.PageContent)
		if doc.Metadata == nil {
			continue
		}
		for k, v := range doc.Metadata {
			log.Debug().Msgf("TextLoaderV2: Metadata: %s: %s", k, v)
			metadata[k] = v
		}
	}

	log.Info().Msgf("TextLoaderV2: Successfully loaded text data from %s", file)
	log.Debug().Msgf("TextLoaderV2: Metadata: %v", metadata)

	return strSlice, metadata, nil
}
