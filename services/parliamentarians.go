package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/OP-interesse-popular/interesse-popular-api/cache"
	"github.com/OP-interesse-popular/interesse-popular-api/models"
)

type Parliamentarians struct {
	cache cache.ParliamentarianCache
}

func NewParliamentarians(cache cache.ParliamentarianCache) *Parliamentarians {
	return &Parliamentarians{cache: cache}
}

func (s *Parliamentarians) ListCamara(ctx context.Context) ([]models.Parliamentarian, error) {
	all, err := s.cache.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cached parliamentarians: %w", err)
	}
	return filterByHouse(all, models.ChamberOfDeputies), nil
}

func (s *Parliamentarians) ListSenado(ctx context.Context) ([]models.Parliamentarian, error) {
	all, err := s.cache.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cached parliamentarians: %w", err)
	}
	return filterByHouse(all, models.FederalSenate), nil
}

func filterByHouse(all []models.Parliamentarian, house string) []models.Parliamentarian {
	result := make([]models.Parliamentarian, 0)
	for _, item := range all {
		if item.LegislativeHouse == house {
			result = append(result, item)
		}
	}
	return sortParliamentarians(result)
}

func sortParliamentarians(parliamentarians []models.Parliamentarian) []models.Parliamentarian {
	sort.Slice(parliamentarians, func(i, j int) bool {
		if parliamentarians[i].LegislativeHouse == parliamentarians[j].LegislativeHouse {
			return parliamentarians[i].Name < parliamentarians[j].Name
		}
		return parliamentarians[i].LegislativeHouse < parliamentarians[j].LegislativeHouse
	})
	return parliamentarians
}
