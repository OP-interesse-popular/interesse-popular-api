package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/OP-interesse-popular/interesse-popular-api/models"
	"github.com/redis/go-redis/v9"
)

const parliamentariansKey = "opip:parliamentarians:v1"

var ErrEmpty = errors.New("parliamentarians cache is empty")

// ParliamentarianCache is the read boundary used by the application service.
type ParliamentarianCache interface {
	List(context.Context) ([]models.Parliamentarian, error)
}

type RedisParliamentarianCache struct{ client *redis.Client }

func NewRedisParliamentarianCache(url string) (ParliamentarianCache, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &RedisParliamentarianCache{client: redis.NewClient(options)}, nil
}

func (r *RedisParliamentarianCache) List(ctx context.Context) ([]models.Parliamentarian, error) {
	raw, err := r.client.Get(ctx, parliamentariansKey).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrEmpty
	}
	if err != nil {
		return nil, err
	}
	var items []models.Parliamentarian
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrEmpty
	}
	return items, nil
}
