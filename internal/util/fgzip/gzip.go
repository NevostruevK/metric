// Package fgzip модуль компрессии/декомпрессии данных по алгоритму gzip
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

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	_, err = gz.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = gz.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}
