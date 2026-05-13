package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

type RawEvent struct {
	EventID     string    `json:"event_id"`
	DeveloperID string    `json:"developer_id"`
	MetricType  string    `json:"metric_type"`
	Value       int       `json:"value"`
	Repository  string    `json:"repository"`
	Timestamp   time.Time `json:"timestamp"`
}

type ProcessedEvent struct {
	RawEvent
	ProcessedAt time.Time `json:"processed_at"`
	ProcessorID string    `json:"processor_id"`
}

func main() {
	log.Println("Starting Processor...")
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	rawQueueURL := os.Getenv("QUEUE_RAW_EVENTS")
	processedQueueURL := os.Getenv("QUEUE_PROCESSED_EVENTS")

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		if awsEndpoint != "" {
			o.BaseEndpoint = aws.String(awsEndpoint)
		}
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					out, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
						QueueUrl: aws.String(rawQueueURL), WaitTimeSeconds: 10,
					})
					if err != nil {
						continue
					}
					for _, msg := range out.Messages {
						var raw RawEvent
						json.Unmarshal([]byte(*msg.Body), &raw)
						if _, err := uuid.Parse(raw.EventID); err == nil {
							proc := ProcessedEvent{RawEvent: raw, ProcessedAt: time.Now(), ProcessorID: "proc-1"}
							b, _ := json.Marshal(proc)
							sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
								QueueUrl: aws.String(processedQueueURL), MessageBody: aws.String(string(b)),
							})
						}
						sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
							QueueUrl: aws.String(rawQueueURL), ReceiptHandle: msg.ReceiptHandle,
						})
					}
				}
			}
		}()
	}
	<-ctx.Done()
	wg.Wait()
}
