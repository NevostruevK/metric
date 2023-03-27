package fgzip

import (
	"bytes"
	"compress/gzip"
	"fmt"
)
func printResult(b bytes.Buffer){
	var sum int
	for _, v:= range b.Bytes(){
		sum += int(v)
//		data := binary.BigEndian.Uint64(mySlice)
	}
	fmt.Printf("Size Sum (%d - %d) \n",b.Len(),sum)
	fmt.Println("--------------------------------")
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
//	fmt.Printf("Compressed data %v\n",data)
//	printResult(b)
	return b.Bytes(), nil
}

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
//	fmt.Printf("Decompressed data %v\n",data)
//	printResult(b)
	return b.Bytes(), nil
}