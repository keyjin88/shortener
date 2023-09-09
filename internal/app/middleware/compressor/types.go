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
	var err error
	if c.gzipWriter != nil {
		err = c.gzipWriter.Close()
		c.gzipWriter = nil
	}
	return errors.Wrap(err, "error while closing gzipWriter")
}

type compressReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating compressor reader")
	}

	return &compressReader{
		reader:     r,
		gzipReader: gzipReader,
	}, nil
}

// Read reads compressed data from the compressor reader.
func (c *compressReader) Read(p []byte) (n int, err error) {
	n, err = c.gzipReader.Read(p)
	if err != nil {
		return n, errors.Wrap(err, "failed to read compressed data")
	}
	return n, nil
}

// Close closes the compressor reader.
func (c *compressReader) Close() error {
	var err error
	if c.reader != nil {
		err = c.reader.Close()
		c.reader = nil
	}
	if c.gzipReader != nil {
		gzipErr := c.gzipReader.Close()
		c.gzipReader = nil
		if err == nil {
			err = gzipErr
		}
	}
	return err
}
