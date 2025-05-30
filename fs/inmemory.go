package fs

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrImageNotFound = errors.New("image not found")
)

type ImageStore interface {
	GetImage(key, fileName string) (*ImageMetadata, error)
	SaveImage(key string, fileName string, imageData []byte)
	DeleteImage(key, fileName string) error
	DeleteKey(key string) bool
	GetImagesForKey(key string) ([]ImageMetadata, bool)
	List() []string
	CountByKey(key string) int
	Clear()
}

type ImageMetadata struct {
	fileName string
	image    []byte
}

type InMemoryImageStore struct {
	images map[string][]ImageMetadata
	mu     sync.RWMutex
}

func NewInMemoryImageStore() ImageStore {
	return &InMemoryImageStore{
		images: make(map[string][]ImageMetadata),
	}
}

func (s *InMemoryImageStore) GetImage(key, fileName string) (*ImageMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	images, exists := s.images[key]
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

func (s *InMemoryImageStore) SaveImage(key string, fileName string, imageData []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	image := ImageMetadata{
		fileName: fileName,
		image:    imageData,
	}
	s.images[key] = append(s.images[key], image)
}

func (s *InMemoryImageStore) DeleteImage(key, fileName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	images, exists := s.images[key]
	if !exists {
		return ErrKeyNotFound
	}
	for i, img := range images {
		if img.fileName == fileName {
			s.images[key] = append(images[:i], images[i+1:]...)
			if len(s.images[key]) == 0 {
				delete(s.images, key)
			}
			return nil
		}
	}
	return ErrImageNotFound
}

func (s *InMemoryImageStore) DeleteKey(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.images[key]; exists {
		delete(s.images, key)
		return true
	}
	return false
}

func (s *InMemoryImageStore) GetImagesForKey(key string) ([]ImageMetadata, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	images, exists := s.images[key]
	return images, exists
}

func (s *InMemoryImageStore) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.images))
	for key := range s.images {
		keys = append(keys, key)
	}
	return keys
}

func (s *InMemoryImageStore) CountByKey(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if images, exists := s.images[key]; exists {
		return len(images)
	}
	return 0
}

func (s *InMemoryImageStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.images = make(map[string][]ImageMetadata)
}
