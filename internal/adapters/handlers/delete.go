package handlers

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/itsbaivab/url-shortener/internal/core/services"
)

type DeleteFunctionHandler struct {
	statsService *services.StatsService
	linkService  *services.LinkService
}

func NewDeleteFunctionHandler(l *services.LinkService, s *services.StatsService) *DeleteFunctionHandler {
	return &DeleteFunctionHandler{linkService: l, statsService: s}
}

func (s *DeleteFunctionHandler) Delete(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]

	err := s.linkService.Delete(ctx, id)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	err = s.statsService.Delete(ctx, id)
	if err != nil {
		// don't fail the delete operation when stats service is down; log and continue
		// this decouples the critical link deletion from analytics cleanup
		// which prevents the entire links flow from crashing if stats pod is deleted
		// the link is successfully deleted even if stats cleanup fails
		log.Printf("failed to delete stats for link '%s': %v", id, err)
	}

	return events.APIGatewayProxyResponse{StatusCode: 204}, nil
}
