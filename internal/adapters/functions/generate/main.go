package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsbaivab/url-shortener/internal/adapters/cache"
	"github.com/itsbaivab/url-shortener/internal/adapters/handlers"
	"github.com/itsbaivab/url-shortener/internal/adapters/repository"
	"github.com/itsbaivab/url-shortener/internal/config"
	"github.com/itsbaivab/url-shortener/internal/core/services"
)

func main() {
	appConfig := config.NewConfig()
	redisAddress, redisPassword, redisDB := appConfig.GetRedisParams()
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)
	linkTableName := appConfig.GetLinkTableName()
	statsTableName := appConfig.GetStatsTableName()

	linkRepo := repository.NewLinkRepository(context.TODO(), linkTableName)
	linkService := services.NewLinkService(linkRepo, cache)

	statsRepo := repository.NewStatsRepository(context.TODO(), statsTableName)
	statsService := services.NewStatsService(statsRepo, cache)

	handler := handlers.NewGenerateLinkFunctionHandler(linkService, statsService)
	lambda.Start(handler.CreateShortLink)
}
