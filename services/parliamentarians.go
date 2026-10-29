package services

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/OP-interesse-popular/interesse-popular-api/cache"
	"github.com/OP-interesse-popular/interesse-popular-api/models"
)

type Parliamentarians struct {
	cache cache.ParliamentarianCache
}

func NewParliamentarians(cache cache.ParliamentarianCache) *Parliamentarians {
	return &Parliamentarians{cache: cache}
}

func (s *Parliamentarians) ListRightWing(ctx context.Context) ([]models.Parliamentarian, error) {
	all, err := s.cache.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cached parliamentarians: %w", err)
	}
	deputies := filter(all, "camara")
	senators := filter(all, "senado")
	parliamentarians := append(deputies, senators...)
	return sortParliamentarians(parliamentarians), nil
}

func (s *Parliamentarians) ListRightWingCamara(ctx context.Context) ([]models.Parliamentarian, error) {
	all, err := s.cache.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cached parliamentarians: %w", err)
	}
	return filter(all, "camara"), nil
}

func (s *Parliamentarians) ListRightWingSenado(ctx context.Context) ([]models.Parliamentarian, error) {
	all, err := s.cache.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cached parliamentarians: %w", err)
	}
	return filter(all, "senado"), nil
}

func filter(all []models.Parliamentarian, house string) []models.Parliamentarian {
	result := make([]models.Parliamentarian, 0)
	for _, item := range all {
		if item.LegislativeHouse == house && isRightWing(item) {
			result = append(result, item)
		}
	}
	return sortParliamentarians(result)
}

func isRightWing(item models.Parliamentarian) bool {
	if strings.ToLower(strings.TrimSpace(item.PartyDirection)) != "direita" {
		return false
	}
	ideology := strings.ToLower(strings.TrimSpace(item.PartyIdeology))
	return ideology == "direita_liberal" || ideology == "extrema_direita"
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
