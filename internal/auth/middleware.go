package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type sessionVerifier interface {
	VerifySessionValid(ctx context.Context, sessionId uuid.UUID) (*Session, error)
}

func Middleware(svc sessionVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionIdStr, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		sessionId, err := uuid.Parse(sessionIdStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		session, err := svc.VerifySessionValid(c.Request.Context(), sessionId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("session", session)
		c.Next()
	}
}
