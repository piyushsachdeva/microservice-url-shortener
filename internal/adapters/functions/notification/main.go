package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsbaivab/url-shortener/internal/adapters/handlers"
)

func main() {
	log.Print("Starting Lambda")
	lambda.Start(handlers.SlackHandler)
}
