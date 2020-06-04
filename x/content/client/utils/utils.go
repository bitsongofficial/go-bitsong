package utils

import (
	"bytes"
	"compress/gzip"
)

var (
	gzipIdent = []byte("\x1F\x8B\x08")
	hlsIdent  = []byte{35, 69, 88, 84, 77, 51, 85, 10, 35} // #EXTM3U + 2byte
)

// IsGzip returns checks if the file contents are gzip compressed
func IsGzip(input []byte) bool {
	return bytes.Equal(input[:3], gzipIdent)
}

// IsHls checks if the file contents are a m3u8 format
func IsHls(input []byte) bool {
	return bytes.Equal(input[:9], hlsIdent)
}

// GzipIt compresses the input ([]byte)
func GzipIt(input []byte) ([]byte, error) {
	// Create gzip writer.
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(input)
	if err != nil {
		return nil, err
	}
	err = w.Close() // You must close this first to flush the bytes to the buffer.
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
