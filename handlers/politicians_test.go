package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OP-interesse-popular/interesse-popular-api/models"
	"github.com/gin-gonic/gin"
)

func TestFilterPoliticiansMatchesOneMandate(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?party=PL&state=SP&status=elected_pending", nil)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = request
	items := []models.Politician{
		{FullName: "Politico A", Mandates: []models.ParliamentaryMandate{
			{PartyAcronym: "PL", State: "SP", Status: models.MandateStatusActive},
			{PartyAcronym: "MISSAO", State: "SP", Status: models.MandateStatusElectedPending},
		}},
		{FullName: "Politico B", Mandates: []models.ParliamentaryMandate{{PartyAcronym: "PL", State: "SP", Status: models.MandateStatusElectedPending}}},
	}
	result := filterPoliticians(context, items, false, false)
	if len(result) != 1 || result[0].FullName != "Politico B" {
		t.Fatalf("must match every mandate filter on one mandate, got %#v", result)
	}
}

func TestFilterPoliticiansInOffice(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?inOffice=true", nil)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = request
	items := []models.Politician{
		{FullName: "Em exercicio", Mandates: []models.ParliamentaryMandate{{Status: models.MandateStatusActive}}},
		{FullName: "Eleito", Mandates: []models.ParliamentaryMandate{{Status: models.MandateStatusElectedPending}}},
	}
	result := filterPoliticians(context, items, true, true)
	if len(result) != 1 || result[0].FullName != "Em exercicio" {
		t.Fatalf("inOffice filter = %#v", result)
	}
}
