package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/itsbaivab/url-shortener/internal/core/domain"
	"github.com/itsbaivab/url-shortener/internal/core/services"
)

type RedirectFunctionHandler struct {
	linkService  *services.LinkService
	statsService *services.StatsService
}

func NewRedirectFunctionHandler(l *services.LinkService, s *services.StatsService) *RedirectFunctionHandler {
	return &RedirectFunctionHandler{linkService: l, statsService: s}
}

func (h *RedirectFunctionHandler) Redirect(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	pathSegments := strings.Split(req.RawPath, "/")
	if len(pathSegments) < 2 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid URL path"}`,
		}, nil
	}

	shortLinkKey := pathSegments[len(pathSegments)-1]
	longLink, err := h.linkService.GetOriginalURL(ctx, shortLinkKey)
	if err != nil || *longLink == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Link not found"}`,
		}, nil
	}

	if err := h.statsService.Create(ctx, domain.Stats{
		Id:        uuid.NewString(),
		LinkID:    shortLinkKey,
		CreatedAt: time.Now(),
		Platform:  domain.PlatformTwitter, // * TODO: Get platform from request
	}); err != nil {
		// don't fail the redirect when stats service is down; log and continue
		// this decouples the critical redirect path from analytics availability
		// which prevents the entire links flow from crashing if stats pod is deleted
		// keep the redirect fast and reliable.
		// Note: stats will be best-effort; failures are non-fatal.
		// Log the error for observability.
		// Using standard log to keep parity with other handlers.
		log.Printf("failed to record stats for link '%s': %v", shortLinkKey, err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMovedPermanently,
		Headers: map[string]string{
			"Location": *longLink,
		},
	}, nil
}
