package plant_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/williamnoble/kube-botany/plant"
	"github.com/williamnoble/kube-botany/repository/store"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

// the newInMemoryStore fn returns a new in-memory store for testing, filepath is incorrect if we use the std func.
func newInMemoryStore(t *testing.T) store.PlantRepository {
	t.Helper()
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	assert.NoError(t, err)
	s, err := store.NewInMemoryStore(false)
	assert.NoError(t, err)
	return s
}

func testPlant(t *testing.T) (*plant.Plant, time.Time) {
	s := newInMemoryStore(t)
	currentTime := time.Now()
	s.NewPlant("FooPlant", "MyBonsai", "bonsai", currentTime)
	p, err := s.GetPlant("FooPlant")
	assert.NoError(t, err)
	return p, currentTime
}

func TestWater(t *testing.T) {
	p, currentTime := testPlant(t)

	// water to 100%
	p.AddWater()
	assert.Equal(t, currentTime, p.LastUpdated)
	assert.Equal(t, 100, p.CurrentWaterLevel())

	// try to water above 100%
	p.AddWater()
	assert.Equal(t, 100, p.CurrentWaterLevel())
}

func TestUpdateWaterConsumption(t *testing.T) {
	p, currentTime := testPlant(t)

	// set the water level to 50, elapse two days and update the plant
	p.Health.CurrentWaterLevel = 50
	secondDay := currentTime.Add(24 * time.Hour * 2)
	p.Update(secondDay)

	// bonsai consumes 1 unit of water per day or 2 units in two days
	assert.Equal(t, 48, p.CurrentWaterLevel())
}

func TestUpdateGrowth(t *testing.T) {
	// growth is 0 when the plant is created
	p, currentTime := testPlant(t)

	// bonsai grows 5 units per day and 15 units in 3 days, it's still seeding
	// it fully matures in 47 days, if this seems too long, pick a sunflower:)
	dayThree := currentTime.Add(24 * time.Hour * 3)
	p.Update(dayThree)
	assert.Equal(t, int64(15), p.CurrentGrowth())
	assert.Equal(t, plant.Seeding.String(), p.GrowthStage())
	assert.Equal(t, p.DaysToMaturity(), 47)

	// bonsai grows 5 units per day and 150 units in 30 days; it's now growing
	// it fully matures in 20 days
	dayThirty := currentTime.Add(24 * time.Hour * 30)
	p.Update(dayThirty)
	assert.Equal(t, int64(150), p.CurrentGrowth())
	assert.Equal(t, plant.Growing.String(), p.GrowthStage())
	assert.Equal(t, p.DaysToMaturity(), 20)

	// bonsai grows 5 units per day and 250 units in 50 days, it's fully matured
	dayFifty := currentTime.Add(24 * time.Hour * 50)
	p.Update(dayFifty)
	assert.Equal(t, int64(250), p.CurrentGrowth())
	assert.Equal(t, plant.Maturing.String(), p.GrowthStage())
	assert.Equal(t, p.DaysToMaturity(), 0)
}

func TestUpdate(t *testing.T) {
	p, currentTime := testPlant(t)
	assert.Equal(t, currentTime, p.LastUpdated)
}
