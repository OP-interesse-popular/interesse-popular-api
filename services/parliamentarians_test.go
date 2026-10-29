package services

import (
	"context"
	"testing"

	"github.com/OP-interesse-popular/interesse-popular-api/models"
)

type cacheStub struct{ items []models.Parliamentarian }

func (stub cacheStub) List(context.Context) ([]models.Parliamentarian, error) {
	return stub.items, nil
}

func TestListCamaraReturnsAllCachedDeputies(t *testing.T) {
	service := NewParliamentarians(cacheStub{items: []models.Parliamentarian{
		{Name: "Deputado PL", Party: "PL", PartyDirection: "direita", PartyIdeology: "direita_liberal", LegislativeHouse: models.ChamberOfDeputies},
		{Name: "Senador PL", Party: "pl", PartyDirection: "direita", PartyIdeology: "extrema_direita", LegislativeHouse: models.FederalSenate},
		{Name: "Deputada não classificada", Party: "OUTRO", LegislativeHouse: models.ChamberOfDeputies},
	}})

	parliamentarians, err := service.ListCamara(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(parliamentarians) != 2 {
		t.Fatalf("parliamentarians length = %d, want 2", len(parliamentarians))
	}
	if parliamentarians[0].Name != "Deputada não classificada" || parliamentarians[1].Name != "Deputado PL" {
		t.Fatalf("unexpected parliamentarians: %#v", parliamentarians)
	}
}

func TestListSenadoReturnsOnlySenators(t *testing.T) {
	service := NewParliamentarians(cacheStub{items: []models.Parliamentarian{
		{Name: "Deputado", LegislativeHouse: models.ChamberOfDeputies},
		{Name: "Senador", LegislativeHouse: models.FederalSenate},
	}})

	parliamentarians, err := service.ListSenado(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(parliamentarians) != 1 || parliamentarians[0].Name != "Senador" {
		t.Fatalf("parliamentarians = %#v", parliamentarians)
	}
}
