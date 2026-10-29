package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/OP-interesse-popular/interesse-popular-api/models"
	"github.com/OP-interesse-popular/interesse-popular-api/services"
	"github.com/TiB-Software/go-commons/commonsmodels"
	commonserrors "github.com/TiB-Software/go-commons/errors"
	commonsutils "github.com/TiB-Software/go-commons/utils"
	"github.com/gin-gonic/gin"
)

type Parliamentarians struct {
	service *services.Parliamentarians
}

func NewParliamentarians(service *services.Parliamentarians) *Parliamentarians {
	return &Parliamentarians{service: service}
}

func (h *Parliamentarians) ListCamara() gin.HandlerFunc {
	return h.list(h.service.ListCamara)
}
func (h *Parliamentarians) ListSenado() gin.HandlerFunc {
	return func(c *gin.Context) {
		parliamentarians, err := h.service.ListSenado(c.Request.Context())
		if err != nil {
			h.sendListError(c, err)
			return
		}
		parliamentarians = applyFilters(c, parliamentarians)
		slog.Info("parliamentarians retrieved", "house", "senate", "count", len(parliamentarians))
		commonsutils.SendResponse(c, commonsmodels.ResponseList[models.Parliamentarian]{Items: parliamentarians, Total: len(parliamentarians)}, http.StatusOK)
	}
}

func (h *Parliamentarians) list(listParliamentarians func(context.Context) ([]models.Parliamentarian, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, apiErr := commonsutils.GetQueryPage(c)
		if apiErr != nil {
			commonsutils.SendErrorResponse(c, apiErr)
			return
		}
		limit := parliamentariansQueryLimit(c)

		parliamentarians, err := listParliamentarians(c.Request.Context())
		if err != nil {
			h.sendListError(c, err)
			return
		}

		parliamentarians = applyFilters(c, parliamentarians)
		slog.Info("parliamentarians retrieved", "count", len(parliamentarians))
		commonsutils.SendResponse(c, paginate(parliamentarians, page, limit), http.StatusOK)
	}
}

func (h *Parliamentarians) sendListError(c *gin.Context, err error) {
	slog.Error("unable to retrieve parliamentarians",
		"error", err,
		"method", c.Request.Method,
		"path", c.FullPath(),
	)
	commonsutils.SendErrorResponse(c, commonserrors.NewApiError(http.StatusBadGateway, commonserrors.InternalServerError(err.Error())))
}

func applyFilters(c *gin.Context, parliamentarians []models.Parliamentarian) []models.Parliamentarian {
	parliamentarians = filterByParties(parliamentarians, c.QueryArray("party"))
	parliamentarians = filterByStates(parliamentarians, c.QueryArray("state"))
	parliamentarians = filterByPartyDirections(parliamentarians, c.QueryArray("direction"))
	return filterByPartyIdeologies(parliamentarians, c.QueryArray("ideology"))
}

func parliamentariansQueryLimit(c *gin.Context) int32 {
	const defaultLimit int32 = 30
	const maxLimit int32 = 100
	value := c.Query("limit")
	if value == "" {
		return defaultLimit
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil || parsed < 1 || parsed > int64(maxLimit) {
		return defaultLimit
	}
	return int32(parsed)
}

func filterByParties(items []models.Parliamentarian, values []string) []models.Parliamentarian {
	parties := make(map[string]struct{})
	for _, value := range values {
		for _, party := range strings.Split(value, ",") {
			if party = strings.ToUpper(strings.TrimSpace(party)); party != "" {
				parties[party] = struct{}{}
			}
		}
	}
	if len(parties) == 0 {
		return items
	}
	filtered := make([]models.Parliamentarian, 0)
	for _, item := range items {
		if _, ok := parties[strings.ToUpper(strings.TrimSpace(item.Party))]; ok {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func filterByStates(items []models.Parliamentarian, values []string) []models.Parliamentarian {
	states := make(map[string]struct{})
	for _, value := range values {
		for _, state := range strings.Split(value, ",") {
			if state = strings.ToUpper(strings.TrimSpace(state)); len(state) == 2 {
				states[state] = struct{}{}
			}
		}
	}
	if len(states) == 0 {
		return items
	}
	filtered := make([]models.Parliamentarian, 0)
	for _, item := range items {
		if _, ok := states[strings.ToUpper(strings.TrimSpace(item.State))]; ok {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func filterByPartyDirections(items []models.Parliamentarian, values []string) []models.Parliamentarian {
	directions := queryValues(values)
	if len(directions) == 0 {
		return items
	}
	filtered := make([]models.Parliamentarian, 0)
	for _, item := range items {
		if _, ok := directions[strings.ToLower(strings.TrimSpace(item.PartyDirection))]; ok {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func filterByPartyIdeologies(items []models.Parliamentarian, values []string) []models.Parliamentarian {
	ideologies := queryValues(values)
	if len(ideologies) == 0 {
		return items
	}
	filtered := make([]models.Parliamentarian, 0)
	for _, item := range items {
		if _, ok := ideologies[strings.ToLower(strings.TrimSpace(item.PartyIdeology))]; ok {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func queryValues(values []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			if item = strings.ToLower(strings.TrimSpace(item)); item != "" {
				result[item] = struct{}{}
			}
		}
	}
	return result
}

func paginate(items []models.Parliamentarian, page, limit int32) commonsmodels.PaginatedResponse[models.Parliamentarian] {
	total := len(items)
	pageCount := int64((total + int(limit) - 1) / int(limit))
	offset := int(commonsutils.CalculateOffset(page, limit))
	if offset >= total {
		return commonsmodels.PaginatedResponse[models.Parliamentarian]{Items: []models.Parliamentarian{}, PageCount: pageCount, Page: int64(page)}
	}
	end := offset + int(limit)
	if end > total {
		end = total
	}
	return commonsmodels.PaginatedResponse[models.Parliamentarian]{Items: items[offset:end], PageCount: pageCount, Page: int64(page)}
}
