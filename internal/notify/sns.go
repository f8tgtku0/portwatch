package notify

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/user/portwatch/internal/state"
)

// snsClient abstracts the SNS publish call for testing.
type snsClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSNotifier sends port-change alerts to an AWS SNS topic.
type SNSNotifier struct {
	client   snsClient
	topicARN string
}

// NewSNSNotifier creates an SNSNotifier using the default AWS credential chain.
func NewSNSNotifier(topicARN string) (*SNSNotifier, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns: load aws config: %w", err)
	}
	return &SNSNotifier{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

// Send publishes a notification message to the configured SNS topic.
func (n *SNSNotifier) Send(change state.Change) error {
	action := "closed"
	if change.Type == state.Opened {
		action = "opened"
	}
	subject := fmt.Sprintf("Port %s: %d", action, change.Port)
	body := fmt.Sprintf("Port %d has been %s on host.", change.Port, action)

	_, err := n.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(n.topicARN),
		Subject:  aws.String(subject),
		Message:  aws.String(body),
	})
	if err != nil {
		return fmt.Errorf("sns: publish: %w", err)
	}
	return nil
}
