package documenthandler

import (
	"chroma-db/internal/constants"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// GoLangLoader struct
type GoLangLoader struct{}

// LoadDocument goes through all directories and loads all Go files along with metadata
func (t *GoLangLoader) LoadDocument(ctx context.Context, rootDir string) (map[string][]string, map[string]constants.Metadata, error) {
	if rootDir == "" {
		return nil, nil, fmt.Errorf("GoLangLoader: Root directory path is empty")
	}

	allFiles := make(map[string][]string, 0)
	allMetadata := make(map[string]constants.Metadata)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".go") {
			strSlice, metadata, err := t.processGoFile(path)
			if err != nil {
				return err
			}

			packageName, ok := metadata["package"]
			if !ok {
				packageName = "unknown"
			}

			allFiles[path] = strSlice
			allMetadata[path] = metadata

			log.Debug().Msgf("GoLangLoader: Loaded Go code from %s, package: %s", path, packageName)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return allFiles, allMetadata, nil
}

// processGoFile processes a single Go file and extracts necessary information
func (t *GoLangLoader) processGoFile(filePath string) ([]string, constants.Metadata, error) {
	if filePath == "" {
		return nil, nil, fmt.Errorf("GoLangLoader: File path is empty")
	}

	if !strings.HasSuffix(filePath, ".go") {
		return nil, nil, nil
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}

	strSlice := []string{}
	metadata := constants.Metadata{}

	metadata["package"] = node.Name.Name

	for _, decl := range node.Decls {
		if fn, isFn := decl.(*ast.FuncDecl); isFn {
			strSlice = append(strSlice, fn.Name.Name)
		}
	}

	metadata["filename"] = filepath.Base(filePath)
	metadata["filesize"] = fmt.Sprintf("%d bytes", getFileSize(filePath))

	log.Debug().Msgf("GoLangLoader: Loaded Go file from %s with package %s", filePath, node.Name.Name)

	return strSlice, metadata, nil
}

// getFileSize returns the size of a file
func getFileSize(filePath string) int64 {
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fi.Size()
}
