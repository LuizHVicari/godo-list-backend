package comment

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luizhvicari/backend/internal/auth"
	platformHTTP "github.com/luizhvicari/backend/internal/platform/http"
)

var errMapper = platformHTTP.NewErrorMapper(
	platformHTTP.E(ErrorCommentNotFound, http.StatusNotFound, ErrorCommentNotFound.Error()),
	platformHTTP.E(ErrorItemNotFound, http.StatusNotFound, ErrorItemNotFound.Error()),
	platformHTTP.E(ErrorForbidden, http.StatusForbidden, ErrorForbidden.Error()),
)

type service interface {
	CreateComment(ctx context.Context, params CreateCommentParams) (*Comment, error)
	GetCommentByID(ctx context.Context, id, itemID, ownerID uuid.UUID) (*Comment, error)
	UpdateComment(ctx context.Context, id, authorID uuid.UUID, content string) error
	DeleteComment(ctx context.Context, id, authorID uuid.UUID) error
	ListCommentsByItemID(ctx context.Context, itemID, ownerID uuid.UUID, filter ListCommentsFilter) (*ListCommentsResult, error)
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{service: service}
}

// Register registers flat mutation routes: POST, PUT /:id, DELETE /:id
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("", h.Create)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}

// RegisterReads registers nested read routes under /projects/:project_id/steps/:step_id/items/:item_id/comments
func (h *Handler) RegisterReads(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)
}

// Create godoc
// @Summary Create a comment on an item
// @Tags comments
// @Accept json
// @Produce json
// @Param request body CreateCommentRequest true "Comment data"
// @Success 201 {object} CommentResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/comments [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	comment, err := h.service.CreateComment(c.Request.Context(), CreateCommentParams{
		ItemID:   itemID,
		AuthorID: session.UserId,
		Content:  req.Content,
	})
	if err != nil {
		errMapper.Respond(c, err, "failed to create comment")
		return
	}

	c.JSON(http.StatusCreated, toResponse(comment))
}

// List godoc
// @Summary List comments for an item
// @Tags comments
// @Produce json
// @Param project_id path string true "Project ID"
// @Param step_id path string true "Step ID"
// @Param item_id path string true "Item ID"
// @Param limit query int false "Max results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} ListCommentsResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{step_id}/items/{item_id}/comments [get]
func (h *Handler) List(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	var req ListCommentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	result, err := h.service.ListCommentsByItemID(c.Request.Context(), itemID, session.UserId, ListCommentsFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		errMapper.Respond(c, err, "failed to list comments")
		return
	}

	results := make([]CommentResponse, len(result.Comments))
	for i, comment := range result.Comments {
		results[i] = toResponse(comment)
	}

	c.JSON(http.StatusOK, ListCommentsResponse{Total: result.Total, Results: results})
}

// GetByID godoc
// @Summary Get a comment by ID
// @Tags comments
// @Produce json
// @Param project_id path string true "Project ID"
// @Param step_id path string true "Step ID"
// @Param item_id path string true "Item ID"
// @Param id path string true "Comment ID"
// @Success 200 {object} CommentResponse
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /v1/projects/{project_id}/steps/{step_id}/items/{item_id}/comments/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	comment, err := h.service.GetCommentByID(c.Request.Context(), id, itemID, session.UserId)
	if err != nil {
		errMapper.Respond(c, err, "failed to get comment")
		return
	}

	c.JSON(http.StatusOK, toResponse(comment))
}

// Update godoc
// @Summary Update a comment
// @Tags comments
// @Accept json
// @Param id path string true "Comment ID"
// @Param request body UpdateCommentRequest true "Comment data"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /v1/comments/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.UpdateComment(c.Request.Context(), id, session.UserId, req.Content); err != nil {
		errMapper.Respond(c, err, "failed to update comment")
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete a comment
// @Tags comments
// @Param id path string true "Comment ID"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /v1/comments/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	session := c.MustGet("session").(*auth.Session)

	if err := h.service.DeleteComment(c.Request.Context(), id, session.UserId); err != nil {
		errMapper.Respond(c, err, "failed to delete comment")
		return
	}

	c.Status(http.StatusNoContent)
}
