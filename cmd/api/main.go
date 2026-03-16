package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/luizhvicari/backend/pkg/config"
)

func main() {

	envConfig := config.Load()

	logger := slog.Default()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	logger.Info("Server running on port " + envConfig.ServerPort)

	r.Run()

}
