package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/OP-interesse-popular/interesse-popular-api/models"
	"github.com/redis/go-redis/v9"
)

const politiciansKey = "opip:politicians:v1"

var ErrPoliticiansEmpty = errors.New("politicians cache is empty")

type PoliticianCache interface {
	List(context.Context) ([]models.Politician, error)
}

type RedisPoliticianCache struct{ client *redis.Client }

func NewRedisPoliticianCache(url string) (PoliticianCache, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &RedisPoliticianCache{client: redis.NewClient(options)}, nil
}

func (r *RedisPoliticianCache) List(ctx context.Context) ([]models.Politician, error) {
	raw, err := r.client.Get(ctx, politiciansKey).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrPoliticiansEmpty
	}
	if err != nil {
		return nil, err
	}
	var items []models.Politician
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrPoliticiansEmpty
	}
	return items, nil
}
