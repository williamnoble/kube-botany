package repository

import (
	"errors"
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/fs"
	"github.com/williamnoble/kube-botany/pkg/plant"
	"maps"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

type PlantRepository interface {
	// NewPlant Create a new plant
	NewPlant(id, friendlyName, plantType string, creationTime time.Time) (*plant.Plant, error)

	// GetPlant Retrieve a plant by ID
	GetPlant(id string) (*plant.Plant, error)

	// DeletePlant Delete a plant
	DeletePlant(id string) error

	// ListPlantsByType List plants by type
	ListPlantsByType(plantType string) (map[string][]string, error)

	// ListAllPlants List all plants
	ListAllPlants() map[string]*plant.Plant

	// UpdatePlants Update all plants' state
	UpdatePlants(ids []string) error

	// UpdatePlantById Updates a specific plant's state
	UpdatePlantById(id string) error

	// GetVarietyUnsafe Get plant type characteristics. This is not thread-safe.
	GetVarietyUnsafe(plantType string) (plant.Variety, error)

	// ListSupportedVarieties lists the types of plant variety supported by the backend
	ListSupportedVarieties() []string

	// Variety returns the characteristics of a particular variety of plant
	Variety(variety string) (plant.Variety, error)

	// ImageExists returns true when an image exists for the given key
	ImageExists(key string, fileName string) bool

	// SetImage saves an image using the given key
	SetImage(id string, fileName string, image []byte)
}

type InMemoryStore struct {
	Plants          map[string]*plant.Plant
	PlantsByVariety map[string][]string
	Varieties       plant.Varieties
	ImageStore      fs.ImageStore
	mu              sync.RWMutex // Mutex for thread-safe access to plants
}

func NewInMemoryStore(populateStore bool, filePaths ...string) (PlantRepository, error) {
	defaultPath := filepath.Join("pkg/plant/", "varieties.json")
	varietiesFilePath := defaultPath

	if len(filePaths) > 0 && filePaths[0] != "" {
		varietiesFilePath = filePaths[0]
	}

	props, err := plant.VarietiesFromJson(varietiesFilePath)
	if err != nil {
		return nil, fmt.Errorf("store: failed to create in-memory store: %w", err)
	}

	s := InMemoryStore{
		Plants:          make(map[string]*plant.Plant),
		Varieties:       props,
		PlantsByVariety: make(map[string][]string),
		ImageStore:      fs.NewInMemoryImageStore(),
	}

	if populateStore {
		s.populateSamplePlants()
	}

	return &s, nil
}

func (s *InMemoryStore) NewPlant(
	id string,
	friendlyName string,
	varietyType string,
	creationTime time.Time) (*plant.Plant, error) {
	health := plant.Health{
		CurrentGrowth:     0,
		CurrentWaterLevel: 50, // 50% watered
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	variety, err := s.GetVarietyUnsafe(varietyType)
	if err != nil {
		return nil, err
	}

	p := &plant.Plant{
		Id:           id,
		FriendlyName: friendlyName,
		Variety:      &variety,
		CreationTime: creationTime,
		LastUpdated:  creationTime,
		Health:       health,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	s.Plants[id] = p
	s.PlantsByVariety[varietyType] = append(s.PlantsByVariety[varietyType], id)

	return p, nil
}

func (s *InMemoryStore) GetPlant(id string) (*plant.Plant, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.Plants[id]
	if !ok {
		return nil, errors.New("plant not found")
	}
	return p, nil
}

func (s *InMemoryStore) UpdatePlantById(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.Plants[id]
	if !ok {
		return errors.New("plant not found")
	}

	p.Update(time.Now())
	return nil
}

func (s *InMemoryStore) DeletePlant(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.Plants[id]
	if !ok {
		return errors.New("plant not found")
	}

	delete(s.Plants, id)

	if typeIDs, ok := s.PlantsByVariety[p.Variety.Type]; ok {
		s.PlantsByVariety[p.Variety.Type] = slices.DeleteFunc(typeIDs, func(plantID string) bool {
			return plantID == id
		})
	}

	return nil
}

func (s *InMemoryStore) ListPlantsByType(plantType string) (map[string][]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plantIDs, ok := s.PlantsByVariety[plantType]
	if !ok {
		return make(map[string][]string), errors.New("plant type not found")
	}

	result := make(map[string][]string)
	result[plantType] = plantIDs
	return result, nil
}

func (s *InMemoryStore) ListAllPlants() map[string]*plant.Plant {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Plants
}

func (s *InMemoryStore) UpdatePlants(ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var failedErrs []error

	for _, id := range ids {
		p, ok := s.Plants[id]
		if !ok {
			failedErrs = append(failedErrs, fmt.Errorf("plant not found: %s", id))
			continue
		}
		p.Update(time.Now())
	}

	if len(failedErrs) > 0 {
		return fmt.Errorf("failed to update plants: %v", errors.Join(failedErrs...))
	}

	return nil
}

// GetVarietyUnsafe is NOT threadsafe. This function should ONLY be called by another function which has
// a mutex lock on the underlying data.
func (s *InMemoryStore) GetVarietyUnsafe(variety string) (plant.Variety, error) {
	v, ok := s.Varieties[variety]
	if !ok {
		return plant.Variety{}, errors.New("plant type not found")
	}
	return v, nil
}

func (s *InMemoryStore) populateSamplePlants() {
	_, _ = s.NewPlant(
		"DefaultBonsai123",
		"my-bonsai",
		"bonsai",
		time.Now(),
	)
	_, _ = s.NewPlant(
		"DefaultSunflower234",
		"my-sunflower",
		"sunflower",
		time.Now(),
	)

}

func (s *InMemoryStore) ListSupportedVarieties() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	varieties := slices.Collect(maps.Keys(s.Varieties))
	return varieties
}

func (s *InMemoryStore) Variety(variety string) (plant.Variety, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.Varieties[variety]; !ok {
		return plant.Variety{}, errors.New("variety not found")
	}
	return s.Varieties[variety], nil
}

func (s *InMemoryStore) ImageExists(key string, fileName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, err := s.ImageStore.GetImage(key, fileName)
	if err != nil {
		return false
	}
	return true
}

func (s *InMemoryStore) SetImage(id string, fileName string, image []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ImageStore.SaveImage(id, fileName, image)
}
