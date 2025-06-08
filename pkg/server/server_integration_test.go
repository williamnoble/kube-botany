package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williamnoble/kube-botany/pkg/repository"
	"github.com/williamnoble/kube-botany/pkg/types"
	"net/http"
	"testing"
	"time"
)

// TestServerIntegration tests the server as a whole, including initialization,
// starting, handling requests, and shutting down.
func TestServerIntegration(t *testing.T) {
	// Create a new in-memory store for testing
	inMemoryStore, err := repository.NewInMemoryStore(false, "../plant/varieties.json")
	require.NoError(t, err)

	// Create a new server
	svr, err := NewServer(inMemoryStore)
	require.NoError(t, err)

	// Use a test port
	testPort := "8081"

	// Start the server in a goroutine
	go func() {
		err := svr.Start(testPort)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("server.Start() error = %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test creating a plant
	t.Run("Create Plant", func(t *testing.T) {
		plantData := map[string]string{
			"id":            "TestPlant1",
			"friendly_name": "Test Bonsai",
			"variety":       "bonsai",
		}
		jsonData, err := json.Marshal(plantData)
		require.NoError(t, err)

		resp, err := http.Post("http://localhost:"+testPort+"/api/plants",
			"application/json",
			bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Test getting the plant
	t.Run("Get Plant", func(t *testing.T) {
		resp, err := http.Get("http://localhost:" + testPort + "/api/plants/TestPlant1")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var plant types.PlantDTO
		err = json.NewDecoder(resp.Body).Decode(&plant)
		require.NoError(t, err)
		assert.Equal(t, "TestPlant1", plant.Id)
		assert.Equal(t, "Test Bonsai", plant.FriendlyName)
	})

	// Test listing plants
	t.Run("List Plants", func(t *testing.T) {
		resp, err := http.Get("http://localhost:" + testPort + "/api/plants")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var plants []types.PlantDTO
		err = json.NewDecoder(resp.Body).Decode(&plants)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(plants), 1)
	})

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = svr.Shutdown(ctx)
	require.NoError(t, err)
}

// TestNewServer tests the NewServer function in isolation
func TestNewServer(t *testing.T) {
	// Create a new in-memory store for testing
	inMemoryStore, err := repository.NewInMemoryStore(false, "../plant/varieties.json")
	require.NoError(t, err)

	// Create a new server
	svr, err := NewServer(inMemoryStore)
	require.NoError(t, err)

	// Verify the server was created correctly
	assert.NotNil(t, svr)
	assert.NotNil(t, svr.Logger)
	assert.NotNil(t, svr.store)
	assert.NotNil(t, svr.renderer)
	assert.NotEmpty(t, svr.templates)
	assert.Equal(t, "pkg/static", svr.staticDir)
}
