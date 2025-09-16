package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/itsbaivab/url-shortener/internal/core/domain"
	"github.com/itsbaivab/url-shortener/internal/core/services"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type RequestBody struct {
	Long string `json:"long"`
}
type GenerateLinkFunctionHandler struct {
	linkService  *services.LinkService
	statsService *services.StatsService
}

func NewGenerateLinkFunctionHandler(l *services.LinkService, s *services.StatsService) *GenerateLinkFunctionHandler {
	return &GenerateLinkFunctionHandler{linkService: l, statsService: s}
}

func (h *GenerateLinkFunctionHandler) CreateShortLink(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody
	err := json.Unmarshal([]byte(req.Body), &requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid JSON"}`,
		}, nil
	}

	if requestBody.Long == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "URL cannot be empty"}`,
		}, nil
	}
	if len(requestBody.Long) < 15 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "URL must be at least 15 characters long"}`,
		}, nil
	}
	if !IsValidLink(requestBody.Long) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid URL format"}`,
		}, nil
	}

	link := domain.Link{
		Id:          GenerateShortURLID(8),
		OriginalURL: requestBody.Long,
		CreatedAt:   time.Now(),
	}

	err = h.linkService.Create(ctx, link)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to create link"}`,
		}, err
	}

	js, err := json.Marshal(link)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to marshal response"}`,
		}, err
	}

	err = h.statsService.Create(ctx, domain.Stats{
		Id:        uuid.NewString(),
		LinkID:    link.Id,
		Platform:  domain.PlatformTwitter,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("failed to create stats: ", err)
	}

	sendMessageToQueue(ctx, link)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func sendMessageToQueue(ctx context.Context, link domain.Link) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err.Error())
		return
	}

	sqsClient := sqs.NewFromConfig(cfg)
	queueUrl := os.Getenv("QueueUrl")

	if queueUrl == "" {
		log.Println("QueueUrl is not set")
		return
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String("The system generated a short URL with the ID " + link.Id),
	})

	if err != nil {
		fmt.Printf("Failed to send message to SQS, %v", err.Error())
	}
}

func GenerateShortURLID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[charIndex.Int64()]
	}
	return string(result)
}

func IsValidLink(link string) bool {
	_, err := url.Parse(link)
	return err == nil && (len(link) > 7) && (link[:7] == "http://" || link[:8] == "https://")
}
