package store

import (
	"errors"
	"fmt"
	"github.com/williamnoble/kube-botany/plant"
	"maps"
	"slices"
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

	// LookupType Get plant type characteristics
	LookupType(plantType string) (plant.Variety, error)

	// ListSupportedVarieties lists the types of plant variety supported by the backend
	ListSupportedVarieties() []string

	// Variety returns the characterists of a particular variety of plant
	Variety(variety string) (plant.Variety, error)
}

type InMemoryStore struct {
	Plants          map[string]*plant.Plant
	PlantsByVariety map[string][]string
	Varieties       plant.Varieties
}

func NewInMemoryStore(populateStore bool) (PlantRepository, error) {
	props, err := plant.VarietiesFromJson()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory store: %w", err)
	}

	s := InMemoryStore{
		Plants:          make(map[string]*plant.Plant),
		Varieties:       props,
		PlantsByVariety: make(map[string][]string),
	}

	if populateStore {
		s.populateSamplePlants()
	}

	return &s, nil
}

func (s *InMemoryStore) NewPlant(
	namespacedName string,
	friendlyName string,
	plantType string,
	creationTime time.Time) (*plant.Plant, error) {
	// Provie a Random CurrentGrowth Rate, our plant requires 1000 CurrentGrowth to reach maturity. Hence,
	// We divide 1000 by the number of days required to reach maturity.
	// minDaysToMaturity = 28 // 4 weeks
	// maxDaysToMaturity = 42 // 6 weeks
	// Random days to maturity between min and max
	// daysToMaturity = rand.Intn(maxDaysToMaturity-minDaysToMaturity+1) + minDaysToMaturity
	// Calculate daily growth rate
	// growthRate = int64(1000 / daysToMaturity)
	//WaterConsumptionRatePerDay: rand.Intn(6) + 2,
	health := plant.Health{

		CurrentGrowth:     0,
		CurrentWaterLevel: 50, // 50% watered
		LastWatered:       time.Now(),
	}

	p := &plant.Plant{
		NamespacedName: namespacedName,
		FriendlyName:   friendlyName,
		Variety:        plantType,
		CreationTime:   creationTime,
		LastUpdated:    creationTime,
		Health:         health,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	s.Plants[namespacedName] = p
	s.PlantsByVariety[plantType] = append(s.PlantsByVariety[plantType], namespacedName)

	return p, nil
}

func (s *InMemoryStore) GetPlant(id string) (*plant.Plant, error) {
	p, ok := s.Plants[id]
	if !ok {
		return nil, errors.New("plant not found")
	}
	return p, nil

}

func (s *InMemoryStore) UpdatePlantById(namespacedName string) error {
	for id, p := range s.Plants {
		if id == namespacedName {
			t, err := s.LookupType(p.Variety)
			if err != nil {
				return fmt.Errorf("failed to update plant: %w", err)
			}
			p.Update(time.Now(), t)
			return nil
		}
	}
	return errors.New("plant not found")
}

func (s *InMemoryStore) DeletePlant(id string) error {
	p, err := s.GetPlant(id)
	if err != nil {
		return fmt.Errorf("failed to delete plant: %w", err)
	}

	delete(s.Plants, id)

	if typeIDs, ok := s.PlantsByVariety[p.Variety]; ok {
		s.PlantsByVariety[p.Variety] = slices.DeleteFunc(typeIDs, func(plantID string) bool {
			return plantID == id
		})
	}

	return nil

}

func (s *InMemoryStore) ListPlantsByType(plantType string) (map[string][]string, error) {
	_, ok := s.PlantsByVariety[plantType]
	if !ok {
		return make(map[string][]string), errors.New("plant type not found")
	}
	return s.PlantsByVariety, nil
}

func (s *InMemoryStore) ListAllPlants() map[string]*plant.Plant {
	return s.Plants
}

func (s *InMemoryStore) UpdatePlants(namespacedNames []string) error {
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

func (s *InMemoryStore) LookupType(plantType string) (plant.Variety, error) {
	v, ok := s.Varieties[plantType]
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
	varieties := slices.Collect(maps.Keys(s.Varieties))
	return varieties
}

func (s *InMemoryStore) Variety(variety string) (plant.Variety, error) {
	if _, ok := s.Varieties[variety]; !ok {
		return plant.Variety{}, errors.New("variety not found")
	}
	return s.Varieties[variety], nil
}
