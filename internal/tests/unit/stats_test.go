package unit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/itsbaivab/url-shortener/internal/adapters/cache"
	"github.com/itsbaivab/url-shortener/internal/adapters/handlers"
	"github.com/itsbaivab/url-shortener/internal/core/domain"
	"github.com/itsbaivab/url-shortener/internal/core/services"
	"github.com/itsbaivab/url-shortener/internal/tests/mock"
	"github.com/stretchr/testify/assert"
)

func TestStatsTest(t *testing.T) {
	mockStatsRepo := mock.NewMockStatsRepo()
	cache := cache.NewRedisCache("localhost:6379", "", 0)
	statsService := services.NewStatsService(mockStatsRepo, cache)

	mockLinkRepo := mock.NewMockLinkRepo()
	linkService := services.NewLinkService(mockLinkRepo, cache)

	apiHander := handlers.NewStatsFunctionHandler(linkService, statsService)

	t.Run("Stats Unit Test", func(t *testing.T) {
		request := events.APIGatewayV2HTTPRequest{
			RawPath: "/stats",
		}

		response, err := apiHander.Stats(context.Background(), request)
		if err != nil {
			t.Fatal(err)
		}

		var links []domain.Link
		err = json.Unmarshal([]byte(response.Body), &links)

		assert.Nil(t, err)
		assert.Equal(t, response.StatusCode, 200)
		assert.Equal(t, len(links), 3)
	})

}
