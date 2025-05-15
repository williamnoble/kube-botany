package plant

//
//import (
//	"github.com/stretchr/testify/assert"
//	"testing"
//	"time"
//)
//
//func TestPlantIntegration(t *testing.T) {
//	currentTime := time.Now()
//	p := NewPlant("MyBonsai", "DefaultBonsai123", "bonsai", currentTime, false)
//
//	// Log initial state
//	t.Logf("Created plant: %s (TypeCharacteristics: %s)", p.FriendlyName, p.TypeCharacteristics)
//	t.Logf("Initial state: Stage=%s, AddWater=%d%%, CurrentGrowth=%d",
//		p.GrowthStage, p.WaterLevel, p.Growth)
//
//	// Simulate 30 days with regular watering
//	for day := 1; day <= 30; day++ {
//		currentTime = currentTime.Add(24 * time.Hour)
//
//		// AddWater the plant when below 50%
//		if p.WaterLevel < 50 {
//			// Add water in multiple 5% increments to reach ~30%
//			waterAdded := 0
//			for i := 0; i < 6 && p.WaterLevel < 50; i++ {
//				increment := p.AddWater(currentTime)
//				waterAdded += increment
//				if increment == 0 {
//					break
//				}
//			}
//			if waterAdded > 0 {
//				t.Logf("Day %d: Watered plant, added %d%%", day, waterAdded)
//			}
//		}
//
//		p.Update(currentTime)
//
//		// Log every 5 days
//		if day%5 == 0 || day == 1 {
//			t.Logf("Day %d: Stage=%s, AddWater=%d%%, CurrentGrowth=%d",
//				day, p.GrowthStage, p.WaterLevel, p.Growth)
//		}
//	}
//
//	// Verify the plant has grown
//	if p.GrowthStage == Seeding {
//		t.Errorf("Plant didn't grow beyond Seeding stage after 30 days")
//	}
//
//	// Now test a dry period - don't water for 10 days
//	t.Logf("\n--- Dry Period Test ---")
//
//	//  track when water hits zero
//	var zeroWaterDay int
//	var zeroWaterGrowth int64
//
//	for day := 1; day <= 10; day++ {
//		// Store previous water level to detect when it hits zero
//		prevWaterLevel := p.WaterLevel
//
//		currentTime = currentTime.Add(24 * time.Hour)
//		p.Update(currentTime)
//
//		// Detect when water first hits zero
//		if prevWaterLevel > 0 && p.WaterLevel == 0 && zeroWaterDay == 0 {
//			zeroWaterDay = day
//			zeroWaterGrowth = p.Growth
//			t.Logf("AddWater reached 0%% on dry day %d with growth at %d",
//				zeroWaterDay, zeroWaterGrowth)
//		}
//
//		t.Logf("Dry Day %d: Stage=%s, AddWater=%d%%, CurrentGrowth=%d",
//			day, p.GrowthStage, p.WaterLevel, p.Growth)
//	}
//
//	// Verify growth stopped when water reached 0%
//	if zeroWaterDay > 0 && p.Growth > zeroWaterGrowth {
//		t.Errorf("Plant continued to grow after water reached 0%%: %d -> %d",
//			zeroWaterGrowth, p.Growth)
//	}
//
//	// Resume watering and verify growth continues
//	t.Logf("\n--- Recovery Test ---")
//
//	// Add water multiple times to ensure the plant has enough water to grow
//	totalWaterAdded := 0
//	for i := 0; i < 20 && p.WaterLevel < 100; i++ {
//		increment := p.AddWater(currentTime)
//		totalWaterAdded += increment
//		if increment == 0 {
//			break
//		}
//	}
//	t.Logf("Watered plant, added %d%%", totalWaterAdded)
//
//	recoveryStartGrowth := p.Growth
//
//	for day := 1; day <= 5; day++ {
//		currentTime = currentTime.Add(24 * time.Hour)
//		p.Update(currentTime)
//
//		// Add more water if needed
//		if day == 3 && p.WaterLevel < 50 {
//			waterAdded := 0
//			for i := 0; i < 10 && p.WaterLevel < 50; i++ {
//				increment := p.AddWater(currentTime)
//				waterAdded += increment
//				if increment == 0 {
//					break
//				}
//			}
//			if waterAdded > 0 {
//				t.Logf("Recovery Day %d: Added %d%% more water", day, waterAdded)
//			}
//		}
//
//		t.Logf("Recovery Day %d: Stage=%s, AddWater=%d%%, CurrentGrowth=%d",
//			day, p.GrowthStage, p.WaterLevel, p.Growth)
//	}
//
//	// Verify growth resumed
//	if p.Growth <= recoveryStartGrowth {
//		t.Errorf("Plant didn't resume growing after watering: %d -> %d",
//			recoveryStartGrowth, p.Growth)
//	}
//}
//
//func TestWater(t *testing.T) {
//	p := NewPlant("test", "My Bonsai", "bonsai", time.Now(), false)
//
//	t.Run("Adding water when below maximum", func(t *testing.T) {
//		p.WaterLevel = 80
//		currentTime := time.Now().Add(24 * time.Hour) // 1 day later
//		added := p.AddWater(currentTime)
//
//		assert.GreaterOrEqual(t, added, 1)
//		assert.LessOrEqual(t, added, 5)
//		assert.Equal(t, currentTime, p.LastWatered)
//		assert.LessOrEqual(t, p.WaterLevel, 100)   // Should never exceed 100
//		assert.GreaterOrEqual(t, p.WaterLevel, 81) // Should have increased from 80
//
//	})
//	t.Run("Adding water when at maximum", func(t *testing.T) {
//		// Incrementally Add Water
//
//		// Test Water full
//		p.WaterLevel = 100
//		currentTime := time.Now().Add(24 * time.Hour)
//		added := p.AddWater(currentTime)
//
//		assert.Equal(t, 0, added)          // Should add nothing
//		assert.Equal(t, 100, p.WaterLevel) // Still at max
//		assert.Equal(t, currentTime, p.LastWatered)
//	})
//}
//
//func TestUpdateWaterConsumption(t *testing.T) {
//	creationTime := time.Now()
//	p := NewPlant("test", "My Bonsai", "bonsai", creationTime, false)
//	p.WaterLevel = 50
//	p.WaterConsumptionRatePerDay = 5
//	twoDay := creationTime.Add(48 * time.Hour) // 2 days later
//	p.updateWaterConsumption(twoDay)
//
//	assert.Equal(t, 40, p.WaterLevel) // 50 - (2 * 5) = 40
//	assert.Equal(t, twoDay, p.LastWatered)
//}
//
//func TestUpdateGrowth(t *testing.T) {
//	creationTime := time.Now()
//	p := NewPlant("DefaultBonsai123", "My Bonsai", "bonsai", creationTime, false)
//	p.WaterLevel = 100 // Optimal water level
//	p.GrowthRate = 10
//
//	threeDay := creationTime.Add(72 * time.Hour) // 3 days later
//	p.updateGrowth(threeDay)
//
//	expectedGrowth := int64(10 * 3) // 10 growth per day * 3 days
//	assert.Equal(t, expectedGrowth, p.Growth)
//	assert.Equal(t, threeDay, p.LastUpdated)
//}
//
//func TestUpdateGrowthWithLowWater(t *testing.T) {
//	// Setup
//	creationTime := time.Now()
//	p := NewPlant("DefaultBonsai123", "My Bonsai", "bonsai", creationTime, false)
//	p.WaterLevel = 25 // 50% growth rate (25/50)
//	p.GrowthRate = 10 // Fixed rate for testing
//
//	twoDay := creationTime.Add(48 * time.Hour) // 2 days later
//	p.updateGrowth(twoDay)
//
//	assert.Equal(t, int64(10), p.Growth) // 10 growth per day * 2 days * 0.5 (water multiplier) = 10
//}
//
//func TestUpdateGrowthStage(t *testing.T) {
//	testCases := []struct {
//		name        string
//		growth      int64
//		expectStage GrowthStage
//	}{
//		{"Seeding stage", 50, Seeding},
//		{"Sprouting stage", 200, Sprouting},
//		{"Growing stage", 600, Growing},
//		{"Maturing stage", 1100, Maturing},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			p := NewPlant("test", "My Bonsai", "bonsai", time.Now(), false)
//			p.Growth = tc.growth
//			p.updateGrowthStage()
//			assert.Equal(t, tc.expectStage, p.GrowthStage)
//		})
//	}
//}
//
//func TestUpdate(t *testing.T) {
//	// Setup
//	creationTime := time.Now()
//	p := NewPlant("test", "My Bonsai", "bonsai", creationTime, false)
//	p.WaterLevel = 60
//	p.WaterConsumptionRatePerDay = 5
//	p.GrowthRate = 20
//
//	// Action - 2 days later
//	twoDay := creationTime.Add(48 * time.Hour)
//	p.Update(twoDay)
//
//	// Assertions
//	assert.Equal(t, 50, p.WaterLevel)    // 60 - (5 * 2) = 50
//	assert.Equal(t, int64(40), p.Growth) // 20 * 2 = 40
//	assert.Equal(t, twoDay, p.LastWatered)
//	assert.Equal(t, twoDay, p.LastUpdated)
//}
//
//func TestImage(t *testing.T) {
//	// Setup
//	p := NewPlant("test-plant", "My Bonsai", "bonsai", time.Now(), false)
//
//	// Action
//	img := p.Image()
//
//	// Assertions
//	today := time.Now().Format("2006-01-02")
//	expectedImg := today + "-test-plant.png"
//	assert.Equal(t, expectedImg, img)
//}
