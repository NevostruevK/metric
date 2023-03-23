package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/NevostruevK/metric/internal/util/metrics"
)

type saver struct {
	file    *os.File
	encoder *json.Encoder
}

func NewSaver(filename string) (*saver, error){
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil{
		fmt.Printf("Can't open %s for writing, error : %v\n",filename,err)
		return nil, err
	}
	return &saver{
		file: file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (s *saver) WriteMetric(rp RepositoryData) error{
	return s.encoder.Encode(rp)
}

func (s *saver) Close() {
	s.file.Close()
}

type loader struct{
	file *os.File
	decoder *json.Decoder
}

func NewLoader(filename string) (*loader, error){
	file, err := os.OpenFile(filename, os.O_RDONLY,0777)
	if err != nil{
		fmt.Printf("Can't open %s for reading, error : %v\n",filename,err)
		return nil, err
	}
	return &loader{
		file: file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (l *loader) ReadMetric() (*metrics.Metrics, error){
	m := metrics.Metrics{}
	err := l.decoder.Decode(&m)
	return &m, err
}

func (l *loader) Close() {
	l.file.Close()
}

