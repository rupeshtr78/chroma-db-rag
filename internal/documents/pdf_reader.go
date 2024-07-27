package documents

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
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

func ParsePDF(path string) ([]string, Metadata, error) {
	pdfStrings := []string{}
	metadata := Metadata{}

	file, pdfReader, err := pdf.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	fileName := strings.Split(file.Name(), "/")[1]

	numPages := pdfReader.NumPage()
	log.Debug().Msgf("Number of pdf pages: %d", numPages)

	for i := 0; i < numPages; i++ {
		page := pdfReader.Page(i + 1)

		text, err := page.GetPlainText(nil)
		if err != nil {
			return nil, nil, err
		}
		pdfStrings = append(pdfStrings, text)
		metadata[fmt.Sprintf("page_%d", i+1)] = map[string]string{
			"file_name": fileName,
			"page_num":  fmt.Sprintf("%d", i+1),
		}

	}

	log.Debug().Msgf("PDF List Length: %v", len(pdfStrings))
	log.Debug().Msgf("Metadata Length: %v", len(metadata))

	return pdfStrings, metadata, nil
}
