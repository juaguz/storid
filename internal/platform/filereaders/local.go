package filereaders

import (
	"context"
	"io"
	"os"
)

type LocalFileReader struct{}

func NewLocalFileReader() *LocalFileReader {
	return &LocalFileReader{}

}

func (l *LocalFileReader) Open(_ context.Context, filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}
