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

	metadata["title"] = strings.Split(file.Name(), "/")[1]

	numPages := pdfReader.NumPage()
	fmt.Printf("Number of pages: %d\n", numPages)

	for i := 0; i < numPages; i++ {
		page := pdfReader.Page(i + 1)
		if err != nil {
			return nil, nil, err
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return nil, nil, err
		}
		pdfStrings = append(pdfStrings, text)

	}

	return pdfStrings, metadata, nil
}
