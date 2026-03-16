package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jackc/pgx/v5"
	_ "github.com/luizhvicari/backend/docs"
	"github.com/luizhvicari/backend/pkg/config"
)

// @title Godo List
// @version 1.0
// @description A Simple To-Do List API built with Go, Gin, and PostgreSQL.
// @BasePath /
// @schemes http

func main() {

	port := flag.String("port", "8080", "server port")
	flag.Parse()

	envConfig := config.Load()

	logger := slog.Default()

	databaseUrl := "postgres://" + envConfig.DatabaseUser + ":" + envConfig.DatabasePassword + "@" + envConfig.DatabaseHost + ":" + strconv.Itoa(envConfig.DatabasePort) + "/" + envConfig.DatabaseName
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		return
	}
	defer conn.Close(context.Background())

	redisClient := redis.NewClient(&redis.Options{
		Addr:     envConfig.CacheHost + ":" + strconv.Itoa(envConfig.CachePort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer redisClient.Close()

	r := gin.Default()
	r.GET("/health", healthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("server running", "port", *port)

	r.Run(":" + *port)

}

// healthCheck godoc
// @Summary Health check
// @Description Returns API health status
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
