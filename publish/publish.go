package publish

import (
	"context"
	er "errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/shivasaicharanruthala/webapp/errors"
	"github.com/shivasaicharanruthala/webapp/log"
	"os"
)

type SNSPublishAPI interface {
	Publish(message []byte, attributes map[string]interface{}) (*sns.PublishOutput, error)
}

func New(logger *log.CustomLogger) (SNSPublishAPI, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Event{
		logger: logger,
		client: sns.NewFromConfig(cfg),
	}, nil
}

type Event struct {
	logger *log.CustomLogger
	client *sns.Client
}

func (e *Event) Publish(message []byte, attributes map[string]interface{}) (*sns.PublishOutput, error) {
	topicArn := os.Getenv("TOPIC_ARN")
	msg := string(message)

	if msg == "" || topicArn == "" {
		lm := log.Message{Level: "ERROR", ErrorMessage: "Message or Topic not provided"}
		e.logger.Log(&lm)

		return nil, errors.NewCustomError(er.New("message or Topic not provided"))
	}

	input := &sns.PublishInput{
		Message:  &msg,
		TopicArn: &topicArn,
	}

	return e.client.Publish(context.TODO(), input)
}
