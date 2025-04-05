package plant

import (
	"testing"
	"time"
)

func TestFernPlantIntegrationTest(t *testing.T) {
	p := NewPlant("My Fern", Fern)

	// Simulate upgradeGrowth over time
	currentTime := p.CreationTime

	// Log initial state
	t.Logf("Created plant: %s (Type: %s)", p.Name, p.Type)
	t.Logf("Initial state: Stage=%s, Health=%.1f%%, Water=%s",
		p.GrowthStage, p.Health, p.WaterLevelPercent())

	// Simulate 60 days with proper watering
	for day := 1; day <= 60; day++ {
		currentTime = currentTime.Add(24 * time.Hour)

		// Water the plant more consistently to maintain optimal conditions
		if p.WaterLevel < p.Characteristics.OptimalWaterMin+10 {
			waterAmount := (p.Characteristics.OptimalWaterMax - p.WaterLevel) * 0.7
			if waterAmount < 0 {
				waterAmount = 0
			}
			_ = p.Water(waterAmount, currentTime)

		}

		// Update the plant
		p.Update(currentTime)

		// Log progress every 5 days
		if day%5 == 0 || day == 1 {
			t.Logf("Day %d: Stage=%s, Health=%.1f%%, Water=%s, Growth=%.1f",
				day, p.GrowthStage, p.Health, p.WaterLevelPercent(), p.Growth)
		}

		// Check if the plant died
		if p.GrowthStage == Dead {
			t.Fatalf("Plant died on day %d with updateHealth %.1f%% and water %s",
				day, p.Health, p.WaterLevelPercent())
		}
	}

	// Verify the plant has matured
	if p.GrowthStage != Maturing {
		t.Errorf("Expected plant to reach Maturing stage, but got %s with upgradeGrowth %.1f",
			p.GrowthStage, p.Growth)
	}
}
