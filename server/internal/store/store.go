package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"comic_reader/pkg/model"
)

type Store struct {
	mu       sync.RWMutex
	filePath string
	comics   map[string]model.Comic
}

func New(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori data: %w", err)
	}

	dataFile := filepath.Join(dataDir, "catalog.json")
	s := &Store{
		filePath: dataFile,
		comics:   make(map[string]model.Comic),
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("gagal membaca file katalog: %w", err)
	}

	var list []model.Comic
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("gagal parse file json katalog: %w", err)
	}

	for _, c := range list {
		s.comics[c.ID] = c
	}
	return nil
}

func (s *Store) save() error {
	list := make([]model.Comic, 0, len(s.comics))
	for _, c := range s.comics {
		list = append(list, c)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal serialisasi json katalog: %w", err)
	}

	return os.WriteFile(s.filePath, data, 0644)
}

func (s *Store) GetAll() []model.Comic {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]model.Comic, 0, len(s.comics))
	for _, c := range s.comics {
		list = append(list, c)
	}
	return list
}

func (s *Store) Get(id string) (model.Comic, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.comics[id]
	return c, ok
}

func (s *Store) Upsert(c model.Comic) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	s.comics[c.ID] = c
	return s.save()
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.comics, id)
	return s.save()
}
