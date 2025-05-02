package kafka_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/messaging/kafka"
)

func TestNewKafkaProducer(t *testing.T) {
	// Test case 1: Test creating a new Kafka producer
	messenger, err := kafka.NewKafkaMessenger("test-topic", fmt.Sprintf("test-group-%d", time.Now().UnixNano()), []string{"localhost:9092"})
	if err != nil {
		t.Errorf("Failed to create Kafka producer: %v", err)
	}
	if err := messenger.Close(); err != nil {
		t.Errorf("Failed to close Kafka producer: %v", err)
	}
}

func TestPubSub(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messenger, err := kafka.NewKafkaMessenger("test-topic", fmt.Sprintf("test-group-%d", time.Now().UnixNano()), []string{"localhost:9092"})
	if err != nil {
		t.Fatalf("Failed to create Kafka messenger: %v", err)
	}
	defer func() {
		if closeErr := messenger.Close(); closeErr != nil {
			t.Errorf("Failed to close Kafka messenger: %v", closeErr)
		}
	}()

	messageChannel, errorChannel, err := messenger.StartConsumer(ctx)
	if err != nil {
		t.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	time.Sleep(2 * time.Second)

	messageToPublish := "Hello, Kafka from Test!"
	testKey := "test-key"

	t.Logf("Test: Publishing message '%s' with key '%s'", messageToPublish, testKey)
	if err := messenger.ProduceMessage(ctx, testKey, messageToPublish); err != nil {
		t.Fatalf("Failed to produce message: %v", err)
	}
	t.Log("Test: Message published.")

	t.Log("Test: Waiting for message from consumer...")
	select {
	case receivedMsg := <-messageChannel:
		t.Logf("Test: Received message '%s'", receivedMsg)
		if receivedMsg != messageToPublish {
			t.Errorf("Expected message: '%s', got: '%s'", messageToPublish, receivedMsg)
		} else {
			t.Log("Test: Received message matches published message.")
		}

	case consumerErr := <-errorChannel:
		t.Fatalf("Consumer reported error: %v", consumerErr)

	case <-ctx.Done():
		t.Fatalf("Test context cancelled or timed out while waiting for message. Error: %v", ctx.Err())
	}

	t.Log("Test: Stopping consumer...")
	if stopErr := messenger.StopConsumer(); stopErr != nil {
		t.Errorf("Failed to stop Kafka consumer: %v", stopErr)
	}
	t.Log("Test: Consumer stopped.")

	t.Log("TestPubSub finished successfully.")
}
