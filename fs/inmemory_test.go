package fs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()

	const fakeKey = "fake-key"
	testData := struct {
		key       string
		fileName  string
		imageData []byte
	}{
		key:       "test-key",
		fileName:  "test-img.jpeg",
		imageData: []byte("fake-image-images"),
	}

	assert.Equal(t, len(store.List()), 0)

	// save the image
	store.SaveImage(testData.key, testData.fileName, testData.imageData)
	assert.Equal(t, len(store.List()), 1)

	// get the image
	img, err := store.GetImage(testData.key, testData.fileName)
	assert.NoError(t, err)
	assert.Equal(t, img.fileName, testData.fileName)
	assert.Equal(t, img.image, testData.imageData)
	assert.Equal(t, store.CountByKey(testData.key), 1)

	// get the image with the wrong input
	_, err = store.GetImage(fakeKey, testData.fileName)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrKeyNotFound)
	_, err = store.GetImage(testData.key, fakeKey)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrImageNotFound)

	// check we can get images for a given key
	images, exists := store.GetImagesForKey(testData.key)
	assert.True(t, exists)
	assert.Equal(t, len(images), 1)

	// delete the image
	err = store.DeleteImage(testData.key, testData.fileName)
	assert.NoError(t, err)
	_, err = store.GetImage(testData.key, testData.fileName)
	assert.Error(t, err)
	assert.Equal(t, len(store.List()), 0)

	// test clean
	store.SaveImage(testData.key, testData.fileName, testData.imageData)
	store.Clear()
	keys := store.List()
	assert.Equal(t, len(keys), 0)
}
