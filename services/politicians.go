package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/OP-interesse-popular/interesse-popular-api/cache"
	"github.com/OP-interesse-popular/interesse-popular-api/models"
)

type Politicians struct{ cache cache.PoliticianCache }

func NewPoliticians(cache cache.PoliticianCache) *Politicians { return &Politicians{cache: cache} }

func (s *Politicians) List(ctx context.Context) ([]models.Politician, error) {
	items, err := s.cache.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cached politicians: %w", err)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].FullName < items[j].FullName })
	return items, nil
}
