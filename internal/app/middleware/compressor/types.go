package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
)

type compressWriter struct {
	writer     http.ResponseWriter
	gzipWriter *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer:     w,
		gzipWriter: gzip.NewWriter(w),
	}
}

// Close closes the compressor writer
func (c *compressWriter) Close() error {
	return c.gzipWriter.Close()
}

type compressReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		reader:     r,
		gzipReader: gzipReader,
	}, nil
}

// Read reads compressed data from the compressor reader
func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.gzipReader.Read(p)
}

// Close closes the compressor reader
func (c *compressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.gzipReader.Close()
}
