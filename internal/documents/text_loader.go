package documents

import (
	"context"
	"os"
	"strings"

	"github.com/tmc/langchaingo/documentloaders"
)

func TextLoader(file string) ([]string, Metadata, error) {

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
	meta := make(Metadata)
	for k, v := range metaData {
		meta[k] = v
	}
	return []string{str}, meta, nil

}

func TextLoaderV2(file string) ([]string, Metadata, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	loader := documentloaders.NewText(f)
	docs, err := loader.Load(context.Background())
	if err != nil {
		return nil, nil, err
	}

	if len(docs) == 0 {
		return nil, nil, nil
	}

	metaData := map[string]string{
		"file": file,
	}
	meta := make(Metadata)
	for k, v := range metaData {
		meta[k] = v
	}

	// Assuming you want to keep the original behavior of returning a slice of strings
	return []string{docs[0].PageContent}, meta, nil
}
