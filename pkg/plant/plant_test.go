package plant

import (
	"testing"
	"time"
)

func TestPlantIntegration(t *testing.T) {
	// Create a new fern
	p := NewPlant("Test Fern", Fern, false)

	// Start at the creation time
	currentTime := p.CreationTime

	// Log initial state
	t.Logf("Created plant: %s (Type: %s)", p.Name, p.Type)
	t.Logf("Initial state: Stage=%s, Water=%s, Growth=%d",
		p.GrowthStage, p.WaterLevelPercent(), p.Growth)

	// Simulate 30 days with regular watering
	for day := 1; day <= 30; day++ {
		currentTime = currentTime.Add(24 * time.Hour)

		// Water the plant when below 50%
		if p.WaterLevel < 50 {
			// Add water in multiple 5% increments to reach ~30%
			waterAdded := 0
			for i := 0; i < 6 && p.WaterLevel < 50; i++ {
				increment := p.Water(currentTime)
				waterAdded += increment
				if increment == 0 {
					break
				}
			}
			if waterAdded > 0 {
				t.Logf("Day %d: Watered plant, added %d%%", day, waterAdded)
			}
		}

		// Update the plant state
		p.Update(currentTime)

		// Log every 5 days
		if day%5 == 0 || day == 1 {
			t.Logf("Day %d: Stage=%s, Water=%s, Growth=%d",
				day, p.GrowthStage, p.WaterLevelPercent(), p.Growth)
		}
	}

	// Verify the plant has grown
	if p.GrowthStage == Seeding {
		t.Errorf("Plant didn't grow beyond Seeding stage after 30 days")
	}

	// Now test a dry period - don't water for 10 days
	t.Logf("\n--- Dry Period Test ---")

	// Variables to track when water hits zero
	var zeroWaterDay int
	var zeroWaterGrowth int64

	for day := 1; day <= 10; day++ {
		// Store previous water level to detect when it hits zero
		prevWaterLevel := p.WaterLevel

		currentTime = currentTime.Add(24 * time.Hour)
		p.Update(currentTime)

		// Detect when water first hits zero
		if prevWaterLevel > 0 && p.WaterLevel == 0 && zeroWaterDay == 0 {
			zeroWaterDay = day
			zeroWaterGrowth = p.Growth
			t.Logf("Water reached 0%% on dry day %d with growth at %d",
				zeroWaterDay, zeroWaterGrowth)
		}

		t.Logf("Dry Day %d: Stage=%s, Water=%s, Growth=%d",
			day, p.GrowthStage, p.WaterLevelPercent(), p.Growth)
	}

	// Verify growth stopped when water reached 0%
	if zeroWaterDay > 0 && p.Growth > zeroWaterGrowth {
		t.Errorf("Plant continued to grow after water reached 0%%: %d -> %d",
			zeroWaterGrowth, p.Growth)
	}

	// Resume watering and verify growth continues
	t.Logf("\n--- Recovery Test ---")

	// Add water multiple times to ensure the plant has enough water to grow
	totalWaterAdded := 0
	for i := 0; i < 20 && p.WaterLevel < 100; i++ {
		increment := p.Water(currentTime)
		totalWaterAdded += increment
		if increment == 0 {
			break
		}
	}
	t.Logf("Watered plant, added %d%%", totalWaterAdded)

	recoveryStartGrowth := p.Growth

	for day := 1; day <= 5; day++ {
		currentTime = currentTime.Add(24 * time.Hour)
		p.Update(currentTime)

		// Add more water if needed
		if day == 3 && p.WaterLevel < 50 {
			waterAdded := 0
			for i := 0; i < 10 && p.WaterLevel < 50; i++ {
				increment := p.Water(currentTime)
				waterAdded += increment
				if increment == 0 {
					break
				}
			}
			if waterAdded > 0 {
				t.Logf("Recovery Day %d: Added %d%% more water", day, waterAdded)
			}
		}

		t.Logf("Recovery Day %d: Stage=%s, Water=%s, Growth=%d",
			day, p.GrowthStage, p.WaterLevelPercent(), p.Growth)
	}

	// Verify growth resumed
	if p.Growth <= recoveryStartGrowth {
		t.Errorf("Plant didn't resume growing after watering: %d -> %d",
			recoveryStartGrowth, p.Growth)
	}
}
