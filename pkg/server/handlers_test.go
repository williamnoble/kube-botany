package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williamnoble/kube-botany/pkg/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// the newInMemoryTestStore fn returns a new in-memory store for testing, filepath is incorrect if we use the std func.
func newInMemoryTestStore(t *testing.T) repository.PlantRepository {
	t.Helper()
	s, err := repository.NewInMemoryStore(false, "../plant/varieties.json")
	require.NoError(t, err)
	return s
}

func TestListPlants(t *testing.T) {
	t.Parallel()
	s := newInMemoryTestStore(t)
	req := httptest.NewRequest(http.MethodGet, "/api/plants", nil)
	_, err := s.NewPlant("TestPlant", "TestBonsai", "bonsai", time.Now())
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	server := &Server{store: s}
	server.Routes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "\"id\":\"TestPlant\"")
}

func TestGetPlant(t *testing.T) {
	t.Parallel()
	s := newInMemoryTestStore(t)
	req := httptest.NewRequest(http.MethodGet, "/api/plants/TestPlant", nil)
	_, err := s.NewPlant("TestPlant", "TestBonsai", "bonsai", time.Now())
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	server := &Server{store: s}
	server.Routes().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "\"id\":\"TestPlant\"")
}

func TestCreatePlant(t *testing.T) {
	t.Parallel()
	s := newInMemoryTestStore(t)
	testRequestBody := struct {
		Id           string `json:"id"`
		FriendlyName string `json:"friendly_name"`
		Variety      string `json:"variety"`
	}{
		Id:           "TestPlant",
		FriendlyName: "TestBonsai",
		Variety:      "bonsai",
	}
	js, err := json.Marshal(testRequestBody)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/api/plants", bytes.NewReader(js))
	rr := httptest.NewRecorder()

	server := &Server{store: s}
	server.Routes().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	fmt.Println(rr.Body.String())
}

func TestDeletePlant(t *testing.T) {
	t.Parallel()
	s := newInMemoryTestStore(t)
	req := httptest.NewRequest(http.MethodGet, "/api/plants/TestPlant", nil)
	// create a test plant
	_, err := s.NewPlant("TestPlant", "TestBonsai", "bonsai", time.Now())
	require.NoError(t, err)
	server := &Server{store: s}
	// delete a test plant
	req = httptest.NewRequest(http.MethodDelete, "/api/plants/TestPlant", nil)
	rr := httptest.NewRecorder()
	server.Routes().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	//assert.Contains(t, rr.Body.String(), "\"id\":\"TestPlant\"")
}
