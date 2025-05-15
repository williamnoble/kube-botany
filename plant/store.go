package plant

import (
	"math/rand"
	"time"
)

type Store struct {
	Plants       map[string]Plant
	PlantTypes   map[string]TypeCharacteristics
	PlantsByType map[string][]string
}

func NewStore() *Store {
	s := Store{
		Plants:       make(map[string]Plant),
		PlantTypes:   make(map[string]TypeCharacteristics),
		PlantsByType: make(map[string][]string),
	}
	return &s
}

func (s *Store) NewPlant(
	NamespacedName string,
	FriendlyName string,
	plantType string,
	creationTime time.Time,
	canDie bool) *Plant {

	// Provide a Random CurrentGrowth Rate, our plant requires 1000 CurrentGrowth to reach maturity. Hence,
	// We divide 1000 by the number of days required to reach maturity.
	minDaysToMaturity := 28 // 4 weeks
	maxDaysToMaturity := 42 // 6 weeks
	// Random days to maturity between min and max
	daysToMaturity := rand.Intn(maxDaysToMaturity-minDaysToMaturity+1) + minDaysToMaturity
	// Calculate daily growth rate
	growthRate := int64(1000 / daysToMaturity)

	return &Plant{
		NamespacedName: NamespacedName,
		FriendlyName:   FriendlyName,
		Type:           plantType,
		CreationTime:   creationTime,
		GrowthStage:    Seeding,
		CurrentGrowth:  0,
		LastUpdated:    creationTime,

		CurrentWaterLevel: 50,
		LastWatered:       creationTime,

		GrowthRate:                 growthRate,
		WaterConsumptionRatePerDay: rand.Intn(6) + 2,
	}
}
