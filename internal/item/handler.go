package item

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luizhvicari/backend/internal/auth"
	platformHTTP "github.com/luizhvicari/backend/internal/platform/http"
)

var errMapper = platformHTTP.NewErrorMapper(
	platformHTTP.E(ErrorItemNotFound, http.StatusNotFound, ErrorItemNotFound.Error()),
	platformHTTP.E(ErrorInvalidFilterParams, http.StatusBadRequest, ErrorInvalidFilterParams.Error()),
	platformHTTP.E(ErrorItemPositionTaken, http.StatusConflict, ErrorItemPositionTaken.Error()),
	platformHTTP.E(ErrorItemNotBelongsToStep, http.StatusUnprocessableEntity, ErrorItemNotBelongsToStep.Error()),
)

type service interface {
	CreateItem(ctx context.Context, params CreateItemParams) error
	GetItemByID(ctx context.Context, id, stepID, ownerID uuid.UUID) (*Item, error)
	UpdateItem(ctx context.Context, id, stepID, ownerID uuid.UUID, params UpdateItemParams) error
	DeleteItem(ctx context.Context, id, stepID, ownerID uuid.UUID) error
	ListItemsByStepID(ctx context.Context, stepID, ownerID uuid.UUID, filter ListItemsFilter) (*ListItemsResult, error)
	RepositionItems(ctx context.Context, params RepositionItemsParams) error
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("", h.Create)
	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
	rg.PUT("/reposition", h.Reposition)
}

// Create godoc
// @Summary Create a new item
// @Tags items
// @Accept json
// @Param project_id path string true "Project ID"
// @Param step_id path string true "Step ID"
// @Param request body CreateItemRequest true "Item data"
// @Success 201
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 409
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{step_id}/items [post]
func (h *Handler) Create(c *gin.Context) {
	stepID, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step_id"})
		return
	}

	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.CreateItem(c.Request.Context(), CreateItemParams{
		StepID:      stepID,
		OwnerID:     session.UserId,
		Name:        req.Name,
		Description: req.Description,
		Priority:    req.Priority,
		Position:    req.Position,
	}); err != nil {
		errMapper.Respond(c, err, "failed to create item")
		return
	}

	c.Status(http.StatusCreated)
}

// List godoc
// @Summary List items for a step
// @Tags items
// @Produce json
// @Param project_id path string true "Project ID"
// @Param step_id path string true "Step ID"
// @Param name query string false "Filter by name"
// @Param sort query string false "Sort field: name, position, priority, created_at, updated_at"
// @Param direction query string false "Sort direction: asc, desc"
// @Param limit query int false "Max results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} ListItemsResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{step_id}/items [get]
func (h *Handler) List(c *gin.Context) {
	stepID, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step_id"})
		return
	}

	var req ListItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	result, err := h.service.ListItemsByStepID(c.Request.Context(), stepID, session.UserId, ListItemsFilter(req))
	if err != nil {
		errMapper.Respond(c, err, "failed to list items")
		return
	}

	results := make([]ItemResponse, len(result.Items))
	for i, item := range result.Items {
		results[i] = toResponse(item)
	}

	c.JSON(http.StatusOK, ListItemsResponse{Total: result.Total, Results: results})
}

// GetByID godoc
// @Summary Get an item by ID
// @Tags items
// @Produce json
// @Param project_id path string true "Project ID"
// @Param step_id path string true "Step ID"
// @Param id path string true "Item ID"
// @Success 200 {object} ItemResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{step_id}/items/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	stepID, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step_id"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	item, err := h.service.GetItemByID(c.Request.Context(), id, stepID, session.UserId)
	if err != nil {
		errMapper.Respond(c, err, "failed to get item")
		return
	}

	c.JSON(http.StatusOK, toResponse(item))
}

// Update godoc
// @Summary Update an item
// @Tags items
// @Accept json
// @Param project_id path string true "Project ID"
// @Param step_id path string true "Step ID"
// @Param id path string true "Item ID"
// @Param request body UpdateItemRequest true "Item data"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 409
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{step_id}/items/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	stepID, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step_id"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.UpdateItem(c.Request.Context(), id, stepID, session.UserId, UpdateItemParams(req)); err != nil {
		errMapper.Respond(c, err, "failed to update item")
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete an item
// @Tags items
// @Param project_id path string true "Project ID"
// @Param step_id path string true "Step ID"
// @Param id path string true "Item ID"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{step_id}/items/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	stepID, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step_id"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.DeleteItem(c.Request.Context(), id, stepID, session.UserId); err != nil {
		errMapper.Respond(c, err, "failed to delete item")
		return
	}

	c.Status(http.StatusNoContent)
}

// Reposition godoc
// @Summary Reposition items within a step
// @Tags items
// @Accept json
// @Param project_id path string true "Project ID"
// @Param step_id path string true "Step ID"
// @Param request body RepositionItemsRequest true "Items with new positions"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 422 {object} map[string]string
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{step_id}/items/reposition [put]
func (h *Handler) Reposition(c *gin.Context) {
	stepID, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step_id"})
		return
	}

	var req RepositionItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]ItemReposition, len(req.Items))
	for i, item := range req.Items {
		id, err := uuid.Parse(item.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id: " + item.ID})
			return
		}
		items[i] = ItemReposition{ID: id, Position: item.Position}
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.RepositionItems(c.Request.Context(), RepositionItemsParams{
		StepID:  stepID,
		OwnerID: session.UserId,
		Items:   items,
	}); err != nil {
		errMapper.Respond(c, err, "failed to reposition items")
		return
	}

	c.Status(http.StatusNoContent)
}
