package step

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luizhvicari/backend/internal/auth"
	platformHTTP "github.com/luizhvicari/backend/internal/platform/http"
)

var errMapper = platformHTTP.NewErrorMapper(
	platformHTTP.E(ErrorStepNotFound, http.StatusNotFound, ErrorStepNotFound.Error()),
	platformHTTP.E(ErrorInvalidFilterParams, http.StatusBadRequest, ErrorInvalidFilterParams.Error()),
	platformHTTP.E(ErrorStepPositionTaken, http.StatusConflict, ErrorStepPositionTaken.Error()),
	platformHTTP.E(ErrorStepNotBelongsToProject, http.StatusUnprocessableEntity, ErrorStepNotBelongsToProject.Error()),
)

type service interface {
	CreateStep(ctx context.Context, params CreateStepParams) error
	GetStepByID(ctx context.Context, id, projectID, ownerID uuid.UUID) (*Step, error)
	UpdateStep(ctx context.Context, id, projectID, ownerID uuid.UUID, params UpdateStepParams) error
	DeleteStep(ctx context.Context, id, projectID, ownerID uuid.UUID) error
	ListStepsByProjectID(ctx context.Context, projectID, ownerID uuid.UUID, filter ListStepsFilter) (*ListStepsResult, error)
	RepositionSteps(ctx context.Context, params RepositionStepsParams) error
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
	rg.GET("/:step_id", h.GetByID)
	rg.PUT("/:step_id", h.Update)
	rg.DELETE("/:step_id", h.Delete)
	rg.PUT("/reposition", h.Reposition)
}

// Create godoc
// @Summary Create a new step
// @Tags steps
// @Accept json
// @Param project_id path string true "Project ID"
// @Param request body CreateStepRequest true "Step data"
// @Success 201
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 409
// @Failure 500
// @Router /v1/projects/{project_id}/steps [post]
func (h *Handler) Create(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	var req CreateStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.CreateStep(c.Request.Context(), CreateStepParams{
		ProjectID: projectID,
		OwnerID:   session.UserId,
		Name:      req.Name,
		Position:  req.Position,
	}); err != nil {
		errMapper.Respond(c, err, "failed to create step")
		return
	}

	c.Status(http.StatusCreated)
}

// List godoc
// @Summary List steps for a project
// @Tags steps
// @Produce json
// @Param project_id path string true "Project ID"
// @Param name query string false "Filter by name"
// @Param sort query string false "Sort field: name, position, created_at, updated_at"
// @Param direction query string false "Sort direction: asc, desc"
// @Param limit query int false "Max results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} ListStepsResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{project_id}/steps [get]
func (h *Handler) List(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	var req ListStepsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	result, err := h.service.ListStepsByProjectID(c.Request.Context(), projectID, session.UserId, ListStepsFilter{
		Name:      req.Name,
		Sort:      req.Sort,
		Direction: req.Direction,
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		errMapper.Respond(c, err, "failed to list steps")
		return
	}

	results := make([]StepResponse, len(result.Steps))
	for i, s := range result.Steps {
		results[i] = toResponse(s)
	}

	c.JSON(http.StatusOK, ListStepsResponse{Total: result.Total, Results: results})
}

// GetByID godoc
// @Summary Get a step by ID
// @Tags steps
// @Produce json
// @Param project_id path string true "Project ID"
// @Param id path string true "Step ID"
// @Success 200 {object} StepResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	id, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	step, err := h.service.GetStepByID(c.Request.Context(), id, projectID, session.UserId)
	if err != nil {
		errMapper.Respond(c, err, "failed to get step")
		return
	}

	c.JSON(http.StatusOK, toResponse(step))
}

// Update godoc
// @Summary Update a step
// @Tags steps
// @Accept json
// @Param project_id path string true "Project ID"
// @Param id path string true "Step ID"
// @Param request body UpdateStepRequest true "Step data"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 409
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	id, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}

	var req UpdateStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.UpdateStep(c.Request.Context(), id, projectID, session.UserId, UpdateStepParams{
		Name:        req.Name,
		Position:    req.Position,
		IsCompleted: req.IsCompleted,
	}); err != nil {
		errMapper.Respond(c, err, "failed to update step")
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete a step
// @Tags steps
// @Param project_id path string true "Project ID"
// @Param id path string true "Step ID"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	id, err := uuid.Parse(c.Param("step_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.DeleteStep(c.Request.Context(), id, projectID, session.UserId); err != nil {
		errMapper.Respond(c, err, "failed to delete step")
		return
	}

	c.Status(http.StatusNoContent)
}

// Reposition godoc
// @Summary Reposition steps within a project
// @Tags steps
// @Accept json
// @Param project_id path string true "Project ID"
// @Param request body RepositionStepsRequest true "Steps with new positions"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 422 {object} map[string]string
// @Failure 500
// @Router /v1/projects/{project_id}/steps/reposition [put]
func (h *Handler) Reposition(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	var req RepositionStepsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	steps := make([]StepReposition, len(req.Steps))
	for i, s := range req.Steps {
		id, err := uuid.Parse(s.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id: " + s.ID})
			return
		}
		steps[i] = StepReposition{ID: id, Position: s.Position}
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.RepositionSteps(c.Request.Context(), RepositionStepsParams{
		ProjectID: projectID,
		OwnerID:   session.UserId,
		Steps:     steps,
	}); err != nil {
		errMapper.Respond(c, err, "failed to reposition steps")
		return
	}

	c.Status(http.StatusNoContent)
}
