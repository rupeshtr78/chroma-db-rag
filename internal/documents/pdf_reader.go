package documents

import (
	"chroma-db/internal/constants"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/rs/zerolog/log"
)

func Pdfmain() {
	pdfPath := "test/Model Params.pdf" // Update this path

	s, m, err := ParsePDF(pdfPath)
	if err != nil {
		log.Err(err).Msg("Failed to parse PDF")
	}

	for _, str := range s {
		fmt.Println(str)
	}

	fmt.Println(m)
}

func ParsePDF(path string) ([]string, constants.Metadata, error) {
	pdfStrings := []string{}
	metadata := constants.Metadata{}

	file, pdfReader, err := pdf.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	fileName := strings.Split(file.Name(), "/")[1]

	numPages := pdfReader.NumPage()
	log.Debug().Msgf("Number of pdf pages: %d", numPages)

	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		page := pdfReader.Page(pageNum)

		text, err := page.GetPlainText(nil)
		if err != nil {
			return nil, nil, err
		}
		pdfStrings = append(pdfStrings, text)
		metadata[fmt.Sprintf("%d", pageNum)] = fileName

	}

	log.Debug().Msgf("PDF List Length: %v", len(pdfStrings))
	log.Debug().Msgf("Metadata Length: %v", len(metadata))

	return pdfStrings, metadata, nil
}
