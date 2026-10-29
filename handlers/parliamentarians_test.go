package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OP-interesse-popular/interesse-popular-api/models"
	"github.com/gin-gonic/gin"
)

func TestParliamentariansQueryLimit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?limit=30", nil)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = request
	if limit := parliamentariansQueryLimit(context); limit != 30 {
		t.Fatalf("limit = %d, want 30", limit)
	}
}

func TestPaginate(t *testing.T) {
	items := []models.Parliamentarian{{Name: "A"}, {Name: "B"}, {Name: "C"}}
	result := paginate(items, 2, 2)
	if result.Page != 2 || result.PageCount != 2 {
		t.Fatalf("pagination = %#v", result)
	}
	if len(result.Items) != 1 || result.Items[0].Name != "C" {
		t.Fatalf("items = %#v", result.Items)
	}
}

func TestFilterByParties(t *testing.T) {
	items := []models.Parliamentarian{{Name: "A", Party: "PL"}, {Name: "B", Party: "PSOL"}}
	result := filterByParties(items, []string{"pl, PSOL"})
	if len(result) != 2 {
		t.Fatalf("items = %#v", result)
	}
}

func TestFilterByStates(t *testing.T) {
	items := []models.Parliamentarian{{Name: "A", State: "SP"}, {Name: "B", State: "RJ"}}
	result := filterByStates(items, []string{"sp,rj"})
	if len(result) != 2 {
		t.Fatalf("items = %#v", result)
	}
}
