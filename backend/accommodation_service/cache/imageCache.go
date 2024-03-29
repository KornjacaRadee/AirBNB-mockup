package cache

import (
	"accommodation_service/config"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type ImageCache struct {
	cli    *redis.Client
	logger *config.Logger
}

// Constructs Redis Client
func New(logger *config.Logger) *ImageCache {
	redisHost := os.Getenv("REDIS_HOST")
	redistPort := os.Getenv("REDIS_PORT")
	redisAddress := fmt.Sprintf("%s:%s", redisHost, redistPort)

	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	return &ImageCache{
		cli:    client,
		logger: logger,
	}
}

// Check connection
func (ic *ImageCache) Ping() {
	val, _ := ic.cli.Ping().Result()
	ic.logger.Printf("Redis ping: %s\n", val)
}

// Set key-value pair with default expiration
func (ic *ImageCache) Post(image *Image) error {
	key := constructKey(image.AccommodationId, image.Id)

	value, err := json.Marshal(image)
	if err != nil {
		return err
	}

	err = ic.cli.Set(key, value, 30*time.Second).Err()

	return err
}

// Get single image by key
func (ic *ImageCache) Get(accID, imageID string) (*Image, error) {
	key := constructKey(accID, imageID)

	val, err := ic.cli.Get(key).Bytes()
	if err != nil {
		return nil, err
	}

	image := &Image{}
	err = json.Unmarshal(val, &image)
	if err != nil {
		return nil, err
	}

	ic.logger.Println("Cache hit")
	return image, nil
}

// Get all images for accommodation
func (ic *ImageCache) GetAll(accID string) (Images, error) {
	key := constructKey(accID, "")

	vals, err := ic.cli.Get(key).Bytes()
	if err != nil {
		return nil, err
	}

	var images Images
	err = json.Unmarshal(vals, &images)
	if err != nil {
		return nil, err
	}

	ic.logger.Println("Cache hit")
	return images, nil
}

// Post all images for accommodation
func (ic *ImageCache) PostAll(accID string, images Images) error {
	key := constructKey(accID, "")

	value, err := json.Marshal(images)
	if err != nil {
		return err
	}

	err = ic.cli.Set(key, value, 120*time.Second).Err()

	return err
}

// Check if given key exists
func (ic *ImageCache) Exists(accID, imageID string) (bool, error) {
	key := constructKey(accID, imageID)

	exists, err := ic.cli.Exists(key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
