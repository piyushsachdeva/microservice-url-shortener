package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsbaivab/url-shortener/internal/adapters/cache"
	"github.com/itsbaivab/url-shortener/internal/adapters/repository/postgres"
	"github.com/itsbaivab/url-shortener/internal/core/domain"
	"github.com/itsbaivab/url-shortener/internal/core/services"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LinkServiceHandler struct {
	linkService *services.LinkService
}

type CreateLinkRequest struct {
	Long string `json:"long" binding:"required"`
}

type DeleteLinkRequest struct {
	ID string `json:"id" binding:"required"`
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

	// Initialize repository and service
	linkRepo := postgres.NewPostgresLinkRepository(db)
	linkService := services.NewLinkService(linkRepo, redisCache)

	// Initialize handler
	handler := &LinkServiceHandler{
		linkService: linkService,
	}

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Link endpoints
	router.PUT("/generate", handler.CreateLink)
	router.GET("/links", handler.GetAllLinks)
	router.DELETE("/delete", handler.DeleteLink)

	// Start server
	port := getEnv("SERVICE_PORT", "8001")
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

	log.Printf("Link Service started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Link Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Link Service stopped")
}

func (h *LinkServiceHandler) CreateLink(c *gin.Context) {
	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate URL
	if len(req.Long) < 15 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL must be at least 15 characters long"})
		return
	}

	// Generate short link
	link := domain.Link{
		Id:          generateShortURLID(8),
		OriginalURL: req.Long,
		CreatedAt:   time.Now(),
	}

	if err := h.linkService.Create(c.Request.Context(), link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the full link object so frontend can display id, original_url, created_at
	c.JSON(http.StatusOK, link)
}

func (h *LinkServiceHandler) GetAllLinks(c *gin.Context) {
	links, err := h.linkService.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("Error getting all links: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get links"})
		return
	}

	c.JSON(http.StatusOK, links)
}

func (h *LinkServiceHandler) DeleteLink(c *gin.Context) {
	var req DeleteLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.linkService.Delete(c.Request.Context(), req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func generateShortURLID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		charIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[charIndex.Int64()]
	}
	return string(result)
}
