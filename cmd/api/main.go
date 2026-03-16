package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	platformHTTP "github.com/luizhvicari/backend/internal/platform/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/luizhvicari/backend/docs"
	"github.com/luizhvicari/backend/internal/auth"
	"github.com/luizhvicari/backend/internal/platform/config"
	"github.com/luizhvicari/backend/internal/platform/crypto"
	platformDb "github.com/luizhvicari/backend/internal/platform/db"
	"github.com/luizhvicari/backend/internal/user"
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
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database connection", "error", err)
		}
	}()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     envConfig.CacheHost + ":" + strconv.Itoa(envConfig.CachePort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("failed to close redis connection", "error", err)
		}
	}()

	r := gin.New()
	r.Use(platformHTTP.RequestLoggerMiddleware(logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{envConfig.CorsAllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.GET("/health", healthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("server running", "port", *port)

	queries := platformDb.New(db)
	hasher := crypto.NewHasher()

	userRepository := user.NewRepository(queries)
	userService := user.NewService(userRepository, hasher)

	authRepository := auth.NewRepository(redisClient)
	authService := auth.NewService(userService, authRepository, hasher)

	authHandler := auth.NewHandler(authService, envConfig.CookieSecure)

	v1 := r.Group("/v1")
	authHandler.Register(v1.Group("/auth"))

	err = r.Run(":" + *port)
	if err != nil {
		logger.Error("Failed to start server", "error", err)
	}

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
