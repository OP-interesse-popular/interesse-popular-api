package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/OP-interesse-popular/interesse-popular-api/models"
	"github.com/OP-interesse-popular/interesse-popular-api/services"
	"github.com/TiB-Software/go-commons/commonsmodels"
	commonserrors "github.com/TiB-Software/go-commons/errors"
	commonsutils "github.com/TiB-Software/go-commons/utils"
	"github.com/gin-gonic/gin"
)

type Politicians struct{ service *services.Politicians }

func NewPoliticians(service *services.Politicians) *Politicians {
	return &Politicians{service: service}
}

// List returns one item per person. Mandate filters must match the same mandate.
func (h *Politicians) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, apiErr := commonsutils.GetQueryPage(c)
		if apiErr != nil {
			commonsutils.SendErrorResponse(c, apiErr)
			return
		}
		inOffice, hasInOffice, err := inOfficeQuery(c)
		if err != nil {
			commonsutils.SendErrorResponse(c, commonserrors.NewApiError(http.StatusBadRequest, commonserrors.BadRequestError("inOffice must be true or false")))
			return
		}
		items, err := h.service.List(c.Request.Context())
		if err != nil {
			slog.Error("unable to retrieve politicians", "error", err, "method", c.Request.Method, "path", c.FullPath())
			commonsutils.SendErrorResponse(c, commonserrors.NewApiError(http.StatusBadGateway, commonserrors.InternalServerError(err.Error())))
			return
		}
		items = filterPoliticians(c, items, inOffice, hasInOffice)
		if err := sortPoliticians(items, c.Query("sort"), c.Query("order")); err != nil {
			commonsutils.SendErrorResponse(c, commonserrors.NewApiError(http.StatusBadRequest, commonserrors.BadRequestError(err.Error())))
			return
		}
		slog.Info("politicians retrieved", "count", len(items))
		commonsutils.SendResponse(c, paginatePoliticians(items, page, parliamentariansQueryLimit(c)), http.StatusOK)
	}
}

func inOfficeQuery(c *gin.Context) (bool, bool, error) {
	raw, exists := c.GetQuery("inOffice")
	if !exists {
		return false, false, nil
	}
	value, err := strconv.ParseBool(raw)
	return value, true, err
}

func filterPoliticians(c *gin.Context, items []models.Politician, inOffice, hasInOffice bool) []models.Politician {
	query := strings.ToLower(strings.TrimSpace(c.Query("q")))
	parties := queryValues(c.QueryArray("party"))
	states := queryValues(c.QueryArray("state"))
	houses := queryValues(c.QueryArray("house"))
	offices := queryValues(c.QueryArray("office"))
	statuses := queryValues(c.QueryArray("status"))
	directions := queryValues(c.QueryArray("direction"))
	ideologies := queryValues(c.QueryArray("ideology"))
	electionYears := queryIntValues(c.QueryArray("electionYear"))
	filtered := make([]models.Politician, 0, len(items))
	for _, politician := range items {
		if query != "" && !strings.Contains(strings.ToLower(politician.FullName), query) {
			continue
		}
		if hasInOffice && hasActiveMandate(politician) != inOffice {
			continue
		}
		if len(parties)+len(states)+len(houses)+len(offices)+len(statuses)+len(directions)+len(ideologies)+len(electionYears) > 0 && !hasMatchingMandate(politician, parties, states, houses, offices, statuses, directions, ideologies, electionYears) {
			continue
		}
		filtered = append(filtered, politician)
	}
	return filtered
}

func hasActiveMandate(politician models.Politician) bool {
	for _, mandate := range politician.Mandates {
		if mandate.Status == models.MandateStatusActive {
			return true
		}
	}
	return false
}

func hasMatchingMandate(politician models.Politician, parties, states, houses, offices, statuses, directions, ideologies map[string]struct{}, electionYears map[int]struct{}) bool {
	for _, mandate := range politician.Mandates {
		if (len(parties) == 0 || containsQueryValue(parties, mandate.PartyAcronym)) &&
			(len(states) == 0 || containsQueryValue(states, mandate.State)) &&
			(len(houses) == 0 || containsQueryValue(houses, mandate.LegislativeHouse)) &&
			(len(offices) == 0 || containsQueryValue(offices, mandate.Office)) &&
			(len(statuses) == 0 || containsQueryValue(statuses, mandate.Status)) &&
			(len(directions) == 0 || containsQueryValue(directions, mandate.PartyDirection)) &&
			(len(ideologies) == 0 || containsQueryValue(ideologies, mandate.PartyIdeology)) &&
			(len(electionYears) == 0 || (mandate.ElectionYear != nil && containsIntQueryValue(electionYears, *mandate.ElectionYear))) {
			return true
		}
	}
	return false
}

func sortPoliticians(items []models.Politician, field, order string) error {
	field = strings.TrimSpace(field)
	if field == "" {
		field = "name"
	}
	order = strings.ToLower(strings.TrimSpace(order))
	if order == "" {
		order = "asc"
	}
	if field != "name" && field != "startsAt" && field != "endsAt" {
		return fmt.Errorf("sort must be name, startsAt or endsAt")
	}
	if order != "asc" && order != "desc" {
		return fmt.Errorf("order must be asc or desc")
	}
	sort.SliceStable(items, func(i, j int) bool {
		var less bool
		switch field {
		case "startsAt":
			less = latestMandateDate(items[i], true).Before(latestMandateDate(items[j], true))
		case "endsAt":
			less = latestMandateDate(items[i], false).Before(latestMandateDate(items[j], false))
		default:
			less = strings.ToLower(items[i].FullName) < strings.ToLower(items[j].FullName)
		}
		if items[i].FullName == items[j].FullName {
			return false
		}
		if order == "desc" {
			return !less
		}
		return less
	})
	return nil
}

func latestMandateDate(politician models.Politician, startsAt bool) time.Time {
	var result time.Time
	for _, mandate := range politician.Mandates {
		date := mandate.StartsAt
		if !startsAt {
			date = mandate.EndsAt
		}
		if date != nil && date.After(result) {
			result = *date
		}
	}
	return result
}

func queryIntValues(values []string) map[int]struct{} {
	result := make(map[int]struct{})
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			parsed, err := strconv.Atoi(strings.TrimSpace(item))
			if err == nil && parsed > 0 {
				result[parsed] = struct{}{}
			}
		}
	}
	return result
}

func containsIntQueryValue(values map[int]struct{}, value int) bool {
	_, ok := values[value]
	return ok
}

func containsQueryValue(values map[string]struct{}, value string) bool {
	_, ok := values[strings.ToLower(strings.TrimSpace(value))]
	return ok
}

func paginatePoliticians(items []models.Politician, page, limit int32) commonsmodels.PaginatedResponse[models.Politician] {
	total := len(items)
	pageCount := int64((total + int(limit) - 1) / int(limit))
	offset := int(commonsutils.CalculateOffset(page, limit))
	if offset >= total {
		return commonsmodels.PaginatedResponse[models.Politician]{Items: []models.Politician{}, PageCount: pageCount, Page: int64(page)}
	}
	end := offset + int(limit)
	if end > total {
		end = total
	}
	return commonsmodels.PaginatedResponse[models.Politician]{Items: items[offset:end], PageCount: pageCount, Page: int64(page)}
}
