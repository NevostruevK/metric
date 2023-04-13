package metrics

import (
	"fmt"

	"github.com/NevostruevK/metric/internal/util/sign"
)

func (m Metrics) CountHash(key string) (hash string, err error) {
	switch m.MType {
	case Gauge:
		hash, err = sign.Hash(fmt.Sprintf("%s:gauge:%f", m.ID, *(m.Value)), key)
	case Counter:
		hash, err = sign.Hash(fmt.Sprintf("%s:counter:%d", m.ID, *(m.Delta)), key)
	default:
		return "", fmt.Errorf(" can't sign metric %v , type %s isnot implemented", m, m.MType)
	}
	if err != nil {
		return "", fmt.Errorf(" can't sign metric %v , error %v", m, err)
	}
	return
}

func (m Metrics) CheckHash(key string) (bool, error) {
	hash, err := m.CountHash(key)
	if err != nil {
		return false, fmt.Errorf(" can't check hash for metric %v , error %v", m, err)
	}
	if m.Hash != hash {
/*		fmt.Println("CheckHash error")
		fmt.Println(m.Hash)
		fmt.Println(hash)
		fmt.Println("----------------------------")
*/		return false, nil
	}
	return true, nil
}

func (m *Metrics) SetHash(key string) error {
	hash, err := m.CountHash(key)
	if err != nil {
		return fmt.Errorf(" can't set hash for metric %v , error %v", m, err)
	}
	m.Hash = hash
//	fmt.Printf("Set hash for %s  hash %s\n", m, m.Hash)
	return nil
}
