package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaMessenger struct {
	Topic         string
	GroupID       string
	BrokerAddress string

	producer *kafka.Writer
	consumer *kafka.Reader

	messageChannel chan string
	errorChannel   chan error
	consumerCancel context.CancelFunc
	consumerDone   chan struct{}
}

const (
	defaultNumPartitions     = 3
	defaultReplicationFactor = 1
	defaultBatchSize         = 100
	defaultBatchTimeout      = 1 * time.Second
)

// NewKafkaMessenger creates a new KafkaMessenger with the given topic, group ID, and broker address.
// It ensures the topic exists with the default number of partitions and replication factor, and
// sets up a producer and consumer for the topic.
// If the topic cannot be created or there is a problem communicating with the kafka broker,
// it returns an error.
//
// Parameters:
//   - topic: the name of the Kafka topic to produce and consume messages from
//   - groupID: the group ID for the Kafka consumer
//   - brokerAddress: the address of the Kafka broker to connect to
//
// Returns:
//   - *KafkaMessenger: a new KafkaMessenger instance
//   - error: an error if the topic cannot be created or if there is a problem communicating with the kafka broker
//
// Example usage:
//   messenger, err := NewKafkaMessenger("my-topic", "my-group", "localhost:9092")
//   if err != nil {
//     log.Fatal(err)
//   }
//   defer messenger.Close()
func NewKafkaMessenger(topic string, groupID string, brokerAddress string) (*KafkaMessenger, error) {
	messenger := &KafkaMessenger{
		Topic:         topic,
		GroupID:       groupID,
		BrokerAddress: brokerAddress,
	}

	messenger.producer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{brokerAddress},
		Topic:        topic,
		Logger:       log.New(log.Writer(), "KAFKA-PRODUCER: ", log.LstdFlags),
		BatchSize:    defaultBatchSize,
		BatchTimeout: defaultBatchTimeout,
	})

	messenger.consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: groupID,
		Logger:  log.New(log.Writer(), "KAFKA-CONSUMER: ", log.LstdFlags),
	})

	if err := messenger.EnsureTopicExists(defaultNumPartitions, defaultReplicationFactor); err != nil {
		if closeErr := messenger.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close producer: %w", closeErr)
		}
		if closeErr := messenger.consumer.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close consumer: %w", closeErr)
		}
		return nil, fmt.Errorf("failed to ensure topic %s: %w", topic, err)
	}
	return messenger, nil
}

// Close closes the underlying Kafka producer and consumer.
//
// If the producer is not nil, it attempts to close the producer and logs an error if it fails.
// If the consumer is not nil, it attempts to close the consumer and logs an error if it fails.
// If any of the close operations fail, it returns an error that includes all the errors.
//
// Returns:
//   - error: an error if any of the close operations fail
//
// Example usage:
//   err := messenger.Close()
//   if err != nil {
//     log.Fatal(err)
//   }
func (m *KafkaMessenger) Close() (err error) {
	if m.producer != nil {
		log.Println("KAFKA-PRODUCER: Closing producer")
		if closeErr := m.producer.Close(); closeErr != nil {
			log.Printf("KAFKA-PRODUCER: Failed to close producer: %v", closeErr)
			err = fmt.Errorf("failed to close producer: %w", closeErr)
		}
		m.producer = nil
	}

	if m.consumer != nil {
		log.Println("KAFKA-CONSUMER: Closing consumer")
		if closeErr := m.consumer.Close(); closeErr != nil {
			log.Printf("KAFKA-CONSUMER: Failed to close consumer: %v", closeErr)
			if err == nil {
				err = fmt.Errorf("failed to close consumer: %w", closeErr)
			} else {
				err = fmt.Errorf("%w; and also failed to close consumer: %w", err, closeErr)
			}
		}
		m.consumer = nil
	}

	return
}

// ProduceMessage produces a message to the Kafka topic associated with this messenger.
// Parameters:
//   - ctx: the context to use for the write operation
//   - key: the key to use for the message (may be empty)
//   - message: the content of the message
//
// Returns:
//   - error: an error if the message could not be written to the topic
//
// Example:
//   err := messenger.ProduceMessage(context.Background(), "key", "message")
//   if err != nil {
// 		handle error
//   }
func (m *KafkaMessenger) ProduceMessage(ctx context.Context, key, message string) error {
	if m.producer == nil {
		return fmt.Errorf("producer is not initialized")
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
	}

	if err := m.producer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// EnsureTopicExists ensures that the topic exists with the given number of partitions and replication factor.
// If the topic already exists, it does nothing.
// If the topic does not exist, it creates it with the given number of partitions and replication factor.
// It returns an error if the topic cannot be created or if there is a problem communicating with the kafka broker.
//
// Parameters:
//   - numPartitions: the number of partitions to create the topic with
//   - replicationFactor: the replication factor to create the topic with
//
// Returns:
//   - error: an error if the topic cannot be created or if there is a problem communicating with the kafka broker
//
// Example usage:
//   err := messenger.EnsureTopicExists(3, 3)
//   if err != nil {
//     log.Fatal(err)
//   }
func (m *KafkaMessenger) EnsureTopicExists(numPartitions, replicationFactor int) error {
	if m.BrokerAddress == "" {
		return fmt.Errorf("broker address is empty")
	}

	conn, err := kafka.DialContext(context.Background(), "tcp", m.BrokerAddress)
	if err != nil {
		return fmt.Errorf("failed to dial kafka broker to create topic: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller connection: %w", err)
	}

	controllerConn, err := kafka.DialContext(context.Background(), "tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to dial controller: %w", err)
	}
	defer controllerConn.Close()

	if existingConfigs, err := controllerConn.ReadPartitions(m.Topic); err == nil && len(existingConfigs) > 0 {
		log.Printf("KAFKA-PRODUCER: Topic %s already exists", m.Topic)
		return nil
	} else if err != kafka.UnknownTopicOrPartition {
		return fmt.Errorf("failed to read partitions for topic %s: %w", m.Topic, err)
	}

	topicConfigs := []kafka.TopicConfig{{
		Topic:             m.Topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}}

	if err := controllerConn.CreateTopics(topicConfigs...); err != nil {
		if err == kafka.TopicAlreadyExists {
			log.Printf("KAFKA-PRODUCER: Topic %s already exists", m.Topic)
			return nil
		}
		return fmt.Errorf("failed to create topic: %w", err)
	}

	log.Printf("KAFKA-PRODUCER: Created topic %s with %d partitions and replication factor %d\n", m.Topic, numPartitions, replicationFactor)
	return nil
}

// StartConsumer initializes the Kafka consumer and starts a goroutine to read messages from the specified topic.
// It returns channels for messages and errors, and an error if the consumer fails to start.
//
// Parameters:
//   - ctx: context.Context for managing the consumer lifecycle
//
// Returns:
//   - <-chan string: a channel for receiving messages as strings
//   - <-chan error: a channel for receiving errors that occur during message consumption
//   - error: an error if the consumer is already started or not initialized
//
// The function checks if the consumer is already started or not initialized, and returns an error in such cases.
// It sets up message and error channels, and a done channel for the consumer lifecycle. The consumer goroutine
// reads messages from the Kafka topic, sends them to the message channel, and handles errors by sending them
// to the error channel. It also commits message offsets after successful handling. The goroutine listens for
// context cancellation to gracefully shut down.
func (m *KafkaMessenger) StartConsumer(ctx context.Context) (<-chan string, <-chan error, error) {
	if m.consumer == nil {
		return nil, nil, fmt.Errorf("consumer is not initialized")
	}

	if m.messageChannel != nil || m.errorChannel != nil {
		return nil, nil, fmt.Errorf("consumer is already started")
	}

	m.messageChannel = make(chan string, 10)
	m.errorChannel = make(chan error, 10)
	m.consumerDone = make(chan struct{})

	consumerCtx, cancel := context.WithCancel(ctx)
	m.consumerCancel = cancel
	log.Printf("KAFKA-CONSUMER: Started consumer goroutine for topic '%s' in group '%s'\n", m.Topic, m.GroupID)
	log.Printf("KAFKA-CONSUMER: Consuming from topic '%s' in group '%s'\n", m.Topic, m.GroupID)

	go func() {
		defer close(m.consumerDone)
		log.Printf("KAFKA-CONSUMER: Goroutine started for topic '%s'\n", m.Topic)

		for {
			select {
			case <-consumerCtx.Done():
				log.Printf("KAFKA-CONSUMER: Goroutine received cancellation signal for topic '%s' in group '%s'\n", m.Topic, m.GroupID)
				close(m.messageChannel)
				close(m.errorChannel)
				log.Printf("KAFKA-CONSUMER: Goroutine for topic '%s' is exiting\n", m.Topic)
				return
			default:
				msg, err := m.consumer.ReadMessage(consumerCtx)
				if err != nil {
					if err == context.Canceled || err == context.DeadlineExceeded {
						log.Printf("KAFKA-CONSUMER: ReadMessage cancelled/timed out for topic '%s'", m.Topic)
						continue
					}
					log.Printf("KAFKA-CONSUMER: Failed to read message from topic '%s': %v", m.Topic, err)
					m.errorChannel <- fmt.Errorf("failed to read message: %w", err)
					continue
				}

				select {
				case m.messageChannel <- string(msg.Value):
					log.Printf("KAFKA-CONSUMER: Sent message to channel from offset %d", msg.Offset)
				case <-consumerCtx.Done():
					log.Printf("KAFKA-CONSUMER: Context cancelled while sending message for topic '%s'", m.Topic)
					continue
				}

				if err := m.consumer.CommitMessages(consumerCtx, msg); err != nil {
					log.Printf("KAFKA-CONSUMER: Failed to commit message for topic '%s': %v", m.Topic, err)
					m.errorChannel <- fmt.Errorf("failed to commit message: %w", err)
				} else {
					log.Printf("KAFKA-CONSUMER: Committed message offset %d", msg.Offset)
				}
			}
		}
	}()

	return m.messageChannel, m.errorChannel, nil
}

// StopConsumer signals the consumer goroutine to stop and waits for it to exit.
// It closes the message and error channels and resets the consumer state.
// If the consumer is not running, it returns an error.
//
// Returns:
//   - error: an error if the consumer is not running
func (m *KafkaMessenger) StopConsumer() error {
    if m.consumerCancel == nil {
        return fmt.Errorf("kafka consumer is not running")
    }

    log.Printf("KAFKA-CONSUMER: Signaling consumer goroutine to stop for topic '%s'", m.Topic)
    m.consumerCancel()

    <-m.consumerDone

    log.Printf("KAFKA-CONSUMER: Consumer goroutine for topic '%s' has stopped.", m.Topic)

    m.messageChannel = nil
    m.errorChannel = nil
    m.consumerCancel = nil
    m.consumerDone = nil

	return nil
}
