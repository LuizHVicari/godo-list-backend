package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type service interface {
	SignIn(ctx context.Context, email string, password string) (*Session, error)
	SignOut(ctx context.Context, sessionId uuid.UUID) error
	SignUp(ctx context.Context, email string, password string) error
	VerifySessionValid(ctx context.Context, sessionId uuid.UUID) (*Session, error)
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/sign-up", h.SignUp)
	rg.POST("/sign-in", h.SignIn)
	rg.POST("/sign-out", h.SignOut)
}

// SignUp godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Param request body SignUpRequest true "Sign up credentials"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/auth/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var request SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.SignUp(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign up"})
		return
	}

	c.Status(http.StatusCreated)
}

// SignIn godoc
// @Summary Authenticate user and create session
// @Tags auth
// @Accept json
// @Param request body SignInRequest true "Sign in credentials"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /v1/auth/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var request SignInRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.service.SignIn(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.SetCookie("session_id", session.ID.String(), sessionTtlSeconds, "/", "", true, true)
	c.Status(http.StatusOK)
}

// SignOut godoc
// @Summary Invalidate current session
// @Tags auth
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/auth/sign-out [post]
func (h *Handler) SignOut(c *gin.Context) {
	sessionIdStr, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id cookie is required"})
		return
	}

	sessionId, err := uuid.Parse(sessionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}

	err = h.service.SignOut(c.Request.Context(), sessionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign out"})
		return
	}

	c.SetCookie("session_id", "", -1, "/", "", true, true)
	c.Status(http.StatusNoContent)
}
