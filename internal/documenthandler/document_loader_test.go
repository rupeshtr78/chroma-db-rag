package documenthandler

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileReader is a mock implementation of the FileReader interface
type MockFileReader struct {
	mock.Mock
	FileReader
}

func (m *MockFileReader) ReadFile(filePath string) (*os.File, error) {
	args := m.Called(filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*os.File), args.Error(1)
}

func TestTextLoader_LoadDocument(t *testing.T) {
	ctx := context.Background()

	t.Run("FileReadError", func(t *testing.T) {
		mockReader := new(MockFileReader)
		mockReader.On("ReadFile", "/path/to/file").Return(nil, errors.New("read error"))

		// Use type assertion to treat MockFileReader as FileReader
		loader := &TextLoader{fileReader: mockReader.FileReader}
		_, _, err := loader.LoadDocument(ctx, "/path/to/file")
		assert.Error(t, err)
	})

	t.Run("EmptyFilePath", func(t *testing.T) {
		mockReader := new(MockFileReader)
		loader := &TextLoader{fileReader: mockReader.FileReader}
		_, _, err := loader.LoadDocument(ctx, "")
		assert.Error(t, err)
	})
}
