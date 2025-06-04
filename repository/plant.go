package repository

import (
	"errors"
	"fmt"
	"github.com/williamnoble/kube-botany/fs"
	"github.com/williamnoble/kube-botany/plant"
	"maps"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

type PlantRepository interface {
	// NewPlant Create a new plant
	NewPlant(namespacedName, friendlyName, plantType string, creationTime time.Time) (*plant.Plant, error)

	// GetPlant Retrieve a plant by ID
	GetPlant(id string) (*plant.Plant, error)

	// DeletePlant Delete a plant
	DeletePlant(id string) error

	// ListPlantsByType List plants by type
	ListPlantsByType(plantType string) (map[string][]string, error)

	// ListAllPlants List all plants
	ListAllPlants() map[string]*plant.Plant

	// UpdatePlants Update all plants' state
	UpdatePlants(namespacedNames []string) error

	// UpdatePlantById Updates a specific plant's state
	UpdatePlantById(namespacedName string) error

	// GetVarietyUnsafe Get plant type characteristics. This is not thread-safe.
	GetVarietyUnsafe(plantType string) (plant.Variety, error)

	// ListSupportedVarieties lists the types of plant variety supported by the backend
	ListSupportedVarieties() []string

	// Variety returns the characterists of a particular variety of plant
	Variety(variety string) (plant.Variety, error)
}

type InMemoryStore struct {
	Plants          map[string]*plant.Plant
	PlantsByVariety map[string][]string
	Varieties       plant.Varieties
	ImageStore      fs.ImageStore
	mu              sync.RWMutex // Mutex for thread-safe access to plants
}

func NewInMemoryStore(populateStore bool) (PlantRepository, error) {
	filePath := filepath.Join("./plant/", "varieties.json")
	props, err := plant.VarietiesFromJson(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory store: %w", err)
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
	namespacedName string,
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
		NamespacedName: namespacedName,
		FriendlyName:   friendlyName,
		Variety:        &variety,
		CreationTime:   creationTime,
		LastUpdated:    creationTime,
		Health:         health,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	s.Plants[namespacedName] = p
	s.PlantsByVariety[varietyType] = append(s.PlantsByVariety[varietyType], namespacedName)

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

func (s *InMemoryStore) UpdatePlantById(namespacedName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, p := range s.Plants {
		if id == namespacedName {
			p.Update(time.Now())
			return nil
		}
	}
	return errors.New("plant not found")
}

func (s *InMemoryStore) DeletePlant(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, err := s.GetPlant(id)
	if err != nil {
		return fmt.Errorf("failed to delete plant: %w", err)
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

	_, ok := s.PlantsByVariety[plantType]
	if !ok {
		return make(map[string][]string), errors.New("plant type not found")
	}
	return s.PlantsByVariety, nil
}

func (s *InMemoryStore) ListAllPlants() map[string]*plant.Plant {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Plants
}

func (s *InMemoryStore) UpdatePlants(namespacedNames []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var failedErrs []error

	for _, namespacedName := range namespacedNames {
		err := s.UpdatePlantById(namespacedName)
		if err != nil {
			failedErrs = append(failedErrs, err)
		}
	}

	if len(failedErrs) > 0 {
		return fmt.Errorf("failedErrs to update plants: %v", errors.Join(failedErrs...))
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
