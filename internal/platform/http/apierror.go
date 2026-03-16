package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type errEntry struct {
	err    error
	status int
	msg    string
}

type ErrorMapper struct {
	entries []errEntry
}

func E(err error, status int, msg string) errEntry {
	return errEntry{err, status, msg}
}

func NewErrorMapper(entries ...errEntry) *ErrorMapper {
	return &ErrorMapper{entries: entries}
}

func (m *ErrorMapper) Respond(c *gin.Context, err error, fallback string) {
	for _, e := range m.entries {
		if errors.Is(err, e.err) {
			c.JSON(e.status, gin.H{"error": e.msg})
			return
		}
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": fallback})
}
