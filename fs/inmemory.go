package fs

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrImageNotFound = errors.New("image not found")
)

type store interface {
	// store actions are grouped by Key (i.e., the NamespacedName) or FileName
	save(namespacedName string, image []byte)
	delete(namespacedName string, filename string)
}

type Image struct {
	fileName string
	image    []byte
}

type InMemoryStore struct {
	data map[string][]Image
	mu   sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string][]Image),
	}
}

func (s *InMemoryStore) GetImage(key, fileName string) (*Image, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	images, exists := s.data[key]
	if !exists {
		return nil, ErrKeyNotFound
	}
	for _, img := range images {
		if img.fileName == fileName {
			return &img, nil
		}
	}
	return nil, ErrImageNotFound
}

func (s *InMemoryStore) SaveImage(key string, fileName string, imageData []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	image := Image{
		fileName: fileName,
		image:    imageData,
	}
	s.data[key] = append(s.data[key], image)
}

func (s *InMemoryStore) DeleteImage(key, fileName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	images, exists := s.data[key]
	if !exists {
		return ErrKeyNotFound
	}
	for i, img := range images {
		if img.fileName == fileName {
			s.data[key] = append(images[:i], images[i+1:]...)
			if len(s.data[key]) == 0 {
				delete(s.data, key)
			}
			return nil
		}
	}
	return ErrImageNotFound
}

func (s *InMemoryStore) DeleteKey(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return true
	}
	return false
}

func (s *InMemoryStore) GetImagesForKey(key string) ([]Image, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	images, exists := s.data[key]
	return images, exists
}

func (s *InMemoryStore) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

func (s *InMemoryStore) CountByKey(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if images, exists := s.data[key]; exists {
		return len(images)
	}
	return 0
}

func (s *InMemoryStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string][]Image)
}
