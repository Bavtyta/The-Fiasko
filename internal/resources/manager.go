package resources

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// ResourceManager manages game resources with caching
type ResourceManager struct {
	images map[string]*ebiten.Image
}

// NewResourceManager creates a new ResourceManager
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		images: make(map[string]*ebiten.Image),
	}
}

// LoadImage loads an image from the given path with caching
// Returns cached image if already loaded, otherwise loads and caches it
// Returns error if loading fails (NO fallback)
func (r *ResourceManager) LoadImage(path string) (*ebiten.Image, error) {
	// Check cache first
	if img, ok := r.images[path]; ok {
		return img, nil
	}

	// Load file
	file, err := os.Open(path)
	if err != nil {
		log.Println("Failed to load image:", path, err)
		return nil, err
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("Failed to decode image:", path, err)
		return nil, err
	}

	// Create Ebiten Image
	ebitenImg := ebiten.NewImageFromImage(img)

	// Cache it
	r.images[path] = ebitenImg
	log.Println("Loaded image:", path)

	return ebitenImg, nil
}
