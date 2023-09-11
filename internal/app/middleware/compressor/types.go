package compressor

import (
	"compress/gzip"
	"github.com/pkg/errors"
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

// Close closes the compressor writer.
func (c *compressWriter) Close() error {
	if err := c.gzipWriter.Close(); err != nil {
		return errors.Wrap(err, "error while closing compressWriter")
	}
	return nil
}

type compressReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating new compressReader")
	}

	return &compressReader{
		reader:     r,
		gzipReader: gzipReader,
	}, nil
}

// Read reads compressed data from the compressor reader.
func (c *compressReader) Read(p []byte) (n int, err error) {
	read, err := c.gzipReader.Read(p)
	if err != nil {
		return read, errors.Wrap(err, "error while reading from gzipReader")
	}
	return read, nil
}

// Close closes the compressor reader.
func (c *compressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return errors.Wrap(err, "error while closing Reader")
	}
	if err := c.gzipReader.Close(); err != nil {
		return errors.Wrap(err, "error while closing gzipReader")
	}
	return nil
}
