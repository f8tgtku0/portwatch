package notify

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/user/portwatch/internal/state"
)

type mockSNSClient struct {
	captured *sns.PublishInput
	err      error
}

func (m *mockSNSClient) Publish(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.captured = params
	return &sns.PublishOutput{}, m.err
}

func TestSNSNotifier_Send_OpenedPort(t *testing.T) {
	mock := &mockSNSClient{}
	n := &SNSNotifier{client: mock, topicARN: "arn:aws:sns:us-east-1:123:test"}

	err := n.Send(state.Change{Port: 8080, Type: state.Opened})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.captured == nil {
		t.Fatal("expected Publish to be called")
	}
	if *mock.captured.Subject != "Port opened: 8080" {
		t.Errorf("unexpected subject: %s", *mock.captured.Subject)
	}
}

func TestSNSNotifier_Send_ClosedPort(t *testing.T) {
	mock := &mockSNSClient{}
	n := &SNSNotifier{client: mock, topicARN: "arn:aws:sns:us-east-1:123:test"}

	err := n.Send(state.Change{Port: 22, Type: state.Closed})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *mock.captured.Subject != "Port closed: 22" {
		t.Errorf("unexpected subject: %s", *mock.captured.Subject)
	}
}

func TestSNSNotifier_Send_PublishError(t *testing.T) {
	mock := &mockSNSClient{err: errors.New("publish failed")}
	n := &SNSNotifier{client: mock, topicARN: "arn:aws:sns:us-east-1:123:test"}

	err := n.Send(state.Change{Port: 443, Type: state.Opened})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewSNSNotifier_ImplementsNotifier(t *testing.T) {
	n := &SNSNotifier{client: &mockSNSClient{}, topicARN: "arn:test"}
	var _ Notifier = n
}
