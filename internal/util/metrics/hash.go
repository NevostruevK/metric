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
		return "", fmt.Errorf(" can't sign metric %v , error %w", m, err)
	}
	return
}

func (m Metrics) CheckHash(key string) (bool, error) {
	hash, err := m.CountHash(key)
	if err != nil {
		return false, fmt.Errorf(" can't check hash for metric %v , error %w", m, err)
	}
	if m.Hash != hash {
		return false, nil
	}
	return true, nil
}

func (m *Metrics) SetHash(key string) error {
	hash, err := m.CountHash(key)
	if err != nil {
		return fmt.Errorf(" can't set hash for metric %v , error %w", m, err)
	}
	m.Hash = hash
	return nil
}
