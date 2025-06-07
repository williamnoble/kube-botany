package render

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williamnoble/kube-botany/pkg/repository"
	"testing"
	"time"
)

// the newInMemoryStore fn returns a new in-memory store for testing, filepath is incorrect if we use the std func.
func newInMemoryStore(t *testing.T) repository.PlantRepository {
	t.Helper()
	s, err := repository.NewInMemoryStore(false, "../plant/varieties.json")
	require.NoError(t, err)
	return s
}

func TestRenderTest(t *testing.T) {
	s := newInMemoryStore(t)
	r := NewASCIIRenderer()

	testBonsai, err := s.NewPlant("FooPlant", "MyBonsai", "bonsai", time.Now())
	assert.NoError(t, err)

	// plant is initially "Seeding".
	output := r.RenderText(testBonsai)
	assert.Contains(t, output, seeding)

	// add 50 days to ensure the plant is fully matured
	dayThirty := time.Now().Add(24 * time.Hour * 50)
	testBonsai.Update(dayThirty)
	output = r.RenderText(testBonsai)
	assert.Contains(t, output, maturing)
}
