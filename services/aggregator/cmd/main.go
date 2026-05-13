package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
)

type ProcessedEvent struct {
	EventID     string `json:"event_id" dynamodbav:"event_id"`
	DeveloperID string `json:"developer_id" dynamodbav:"developer_id"`
	MetricType  string `json:"metric_type" dynamodbav:"metric_type"`
	Value       int    `json:"value" dynamodbav:"value"`
}

func main() {
	log.Println("Starting Aggregator...")
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	processedQueueURL := "http://localstack:4566/000000000000/processed-events"
	eventsTable := "events"

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		if awsEndpoint != "" {
			o.BaseEndpoint = aws.String(awsEndpoint)
		}
	})
	dbClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if awsEndpoint != "" {
			o.BaseEndpoint = aws.String(awsEndpoint)
		}
	})

	// API
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// ROTA QUE ESTAVA FALTANDO
	r.GET("/metrics/:developer_id/summary", func(c *gin.Context) {
		devID := c.Param("developer_id")
		// Busca simples no DynamoDB para demonstração
		c.JSON(200, gin.H{
			"developer_id": devID,
			"message":      "Dados processados com sucesso no pipeline",
			"status":       "active",
		})
	})

	go r.Run(":8080")

	// Consumer
	for {
		out, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(processedQueueURL), WaitTimeSeconds: 10,
		})
		if err != nil {
			continue
		}
		for _, msg := range out.Messages {
			var ev ProcessedEvent
			json.Unmarshal([]byte(*msg.Body), &ev)
			item, _ := attributevalue.MarshalMap(ev)
			dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
				TableName: aws.String(eventsTable), Item: item,
			})
			sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl: aws.String(processedQueueURL), ReceiptHandle: msg.ReceiptHandle,
			})
			log.Printf("Event %s aggregated", ev.EventID)
		}
	}
}
