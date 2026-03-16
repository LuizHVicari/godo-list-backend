package project

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luizhvicari/backend/internal/auth"
	platformHTTP "github.com/luizhvicari/backend/internal/platform/http"
)

var errMapper = platformHTTP.NewErrorMapper(
	platformHTTP.E(ErrorProjectNotFound, http.StatusNotFound, ErrorProjectNotFound.Error()),
	platformHTTP.E(ErrorInvalidFilterParams, http.StatusBadRequest, ErrorInvalidFilterParams.Error()),
)

type service interface {
	CreateProject(ctx context.Context, params CreateProjectParams) error
	GetProjectById(ctx context.Context, id, ownerID uuid.UUID) (*Project, error)
	UpdateProject(ctx context.Context, id, ownerID uuid.UUID, params UpdateProjectParams) error
	DeleteProject(ctx context.Context, id, ownerID uuid.UUID) error
	ListProjectsByOwnerID(ctx context.Context, ownerID uuid.UUID, filter ListProjectsFilter) (*ListProjectsResult, error)
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
	rg.GET("/:project_id", h.GetByID)
	rg.PUT("/:project_id", h.Update)
	rg.DELETE("/:project_id", h.Delete)
}

// Create godoc
// @Summary Create a new project
// @Tags projects
// @Accept json
// @Param request body CreateProjectRequest true "Project data"
// @Success 201
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /v1/projects [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.CreateProject(c.Request.Context(), CreateProjectParams{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     session.UserId,
	}); err != nil {
		errMapper.Respond(c, err, "failed to create project")
		return
	}

	c.Status(http.StatusCreated)
}

// List godoc
// @Summary List projects for the authenticated user
// @Tags projects
// @Produce json
// @Param name query string false "Filter by name"
// @Param sort query string false "Sort field: name, created_at, updated_at"
// @Param direction query string false "Sort direction: asc, desc"
// @Param limit query int false "Max results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} ListProjectsResponse
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /v1/projects [get]
func (h *Handler) List(c *gin.Context) {
	var req ListProjectsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	result, err := h.service.ListProjectsByOwnerID(c.Request.Context(), session.UserId, ListProjectsFilter{
		Name:      req.Name,
		Sort:      req.Sort,
		Direction: req.Direction,
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		errMapper.Respond(c, err, "failed to list projects")
		return
	}

	results := make([]ProjectResponse, len(result.Projects))
	for i, p := range result.Projects {
		results[i] = toResponse(p)
	}

	c.JSON(http.StatusOK, ListProjectsResponse{Total: result.Total, Results: results})
}

// GetByID godoc
// @Summary Get a project by ID
// @Tags projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} ProjectResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	project, err := h.service.GetProjectById(c.Request.Context(), id, session.UserId)
	if err != nil {
		errMapper.Respond(c, err, "failed to get project")
		return
	}

	c.JSON(http.StatusOK, toResponse(project))
}

// Update godoc
// @Summary Update a project
// @Tags projects
// @Accept json
// @Param id path string true "Project ID"
// @Param request body UpdateProjectRequest true "Project data"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.UpdateProject(c.Request.Context(), id, session.UserId, UpdateProjectParams{
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		errMapper.Respond(c, err, "failed to update project")
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete a project
// @Tags projects
// @Param id path string true "Project ID"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /v1/projects/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.DeleteProject(c.Request.Context(), id, session.UserId); err != nil {
		errMapper.Respond(c, err, "failed to delete project")
		return
	}

	c.Status(http.StatusNoContent)
}
