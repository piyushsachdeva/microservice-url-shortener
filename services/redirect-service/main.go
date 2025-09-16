package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsbaivab/url-shortener/internal/adapters/cache"
	"github.com/itsbaivab/url-shortener/internal/adapters/repository/postgres"
	"github.com/itsbaivab/url-shortener/internal/core/domain"
	"github.com/itsbaivab/url-shortener/internal/core/services"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type RedirectServiceHandler struct {
	linkService  *services.LinkService
	statsService *services.StatsService
}

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "urlshortener")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Redis connection
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisCache := cache.NewRedisCache(redisHost+":"+redisPort, "", 0)

	// Initialize repositories and services
	linkRepo := postgres.NewPostgresLinkRepository(db)
	statsRepo := postgres.NewPostgresStatsRepository(db)
	linkService := services.NewLinkService(linkRepo, redisCache)
	statsService := services.NewStatsService(statsRepo, redisCache)

	// Initialize handler
	handler := &RedirectServiceHandler{
		linkService:  linkService,
		statsService: statsService,
	}

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Redirect endpoint
	router.GET("/redirect/:id", handler.Redirect)

	// Start server
	port := getEnv("SERVICE_PORT", "8002")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Redirect Service started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Redirect Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Redirect Service stopped")
}

func (h *RedirectServiceHandler) Redirect(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Get original URL
	originalURL, err := h.linkService.GetOriginalURL(c.Request.Context(), id)
	if err != nil || originalURL == nil || *originalURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Create stats entry asynchronously
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		stats := domain.Stats{
			Id:        uuid.New().String(),
			LinkID:    id,
			Platform:  domain.PlatformUnknown, // TODO: Detect platform from user agent
			CreatedAt: time.Now(),
		}

		if err := h.statsService.Create(ctx, stats); err != nil {
			log.Printf("Failed to create stats: %v", err)
		}
	}()

	// Redirect to original URL
	c.Redirect(http.StatusMovedPermanently, *originalURL)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
