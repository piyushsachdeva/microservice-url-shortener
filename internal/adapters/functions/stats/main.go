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
	linkTableName := appConfig.GetLinkTableName()
	statsTableName := appConfig.GetStatsTableName()

	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)

	linkRepo := repository.NewLinkRepository(context.TODO(), linkTableName)
	statsRepo := repository.NewStatsRepository(context.TODO(), statsTableName)

	linkService := services.NewLinkService(linkRepo, cache)
	statsService := services.NewStatsService(statsRepo, cache)

	handler := handlers.NewStatsFunctionHandler(linkService, statsService)

	lambda.Start(handler.Stats)
}
