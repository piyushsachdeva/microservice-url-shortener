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
	"github.com/itsbaivab/url-shortener/internal/adapters/cache"
	"github.com/itsbaivab/url-shortener/internal/adapters/repository/postgres"
	"github.com/itsbaivab/url-shortener/internal/core/services"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type StatsServiceHandler struct {
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
	handler := &StatsServiceHandler{
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

	// Stats endpoints
	router.GET("/stats", handler.GetStats)
	router.GET("/stats/:id", handler.GetStatsByLinkID)

	// Start server
	port := getEnv("SERVICE_PORT", "8003")
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

	log.Printf("Stats Service started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Stats Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Stats Service stopped")
}

func (h *StatsServiceHandler) GetStats(c *gin.Context) {
	// Get all links with their stats
	links, err := h.linkService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enhance links with stats
	for i, link := range links {
		stats, err := h.statsService.GetStatsByLinkID(c.Request.Context(), link.Id)
		if err != nil {
			log.Printf("Error getting stats for link '%s': %v", link.Id, err)
			continue
		}
		links[i].Stats = stats
	}

	c.JSON(http.StatusOK, links)
}

func (h *StatsServiceHandler) GetStatsByLinkID(c *gin.Context) {
	linkID := c.Param("id")
	if linkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Link ID parameter is required"})
		return
	}

	stats, err := h.statsService.GetStatsByLinkID(c.Request.Context(), linkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
