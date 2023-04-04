package fgzip

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func Decompress(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed init decompress reader: %v", err)
	}
	defer gz.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(gz)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}
	return b.Bytes(), nil
}
