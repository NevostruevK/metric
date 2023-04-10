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
	init    bool
}

func NewSaver(filename string) (*saver, error) {
	if filename == "" {
		return &saver{init: false}, nil
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Printf("Can't open %s for writing, error : %v\n", filename, err)
		return &saver{init: false}, err
	}
	return &saver{
		file:    file,
		encoder: json.NewEncoder(file),
		init:    true,
	}, nil
}

func (s *saver) WriteMetric(rp RepositoryData) error {
	return s.encoder.Encode(rp)
}

func (s *saver) Close() (err error) {
	if !s.init {
		return nil
	}
	if err = s.file.Close(); err != nil {
		fmt.Printf("Closing file %s  with the error : %v\n", s.file.Name(), err)
	}
	return err
}

type loader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewLoader(filename string) (*loader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Printf("Can't open %s for reading, error : %v\n", filename, err)
		return nil, err
	}
	return &loader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (l *loader) ReadMetric() (*metrics.Metrics, error) {
	m := metrics.Metrics{}
	err := l.decoder.Decode(&m)
	return &m, err
}

func (l *loader) Close() (err error) {
	if err = l.file.Close(); err != nil {
		fmt.Printf("Closing file %s  with the error : %v\n", l.file.Name(), err)
	}
	return err
}
