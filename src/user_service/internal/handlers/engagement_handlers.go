package handlers

import (
	"encoding/json"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EventType int

const (
	Like EventType = iota
	View
)

type EngagementEvent struct {
	EventID uint32    `json:"id"`
	Type    EventType `json:"type"`
	PostID  int64     `json:"post"`
	Author  int       `json:"author"`
}

func setKafkaUp() (*kafka.Producer, error) {
	kafkaServer := os.Getenv("KAFKA_HOST") + ":" + os.Getenv("KAFKA_BOOTSTRAP_PORT")
	kafkaConfig := kafka.ConfigMap{"bootstrap.servers": kafkaServer}

	kafkaProd, err := kafka.NewProducer(&kafkaConfig)
	if err != nil {
		return nil, err
	}

	return kafkaProd, nil
}

func (s *ServerHandler) sendToKafka(event EngagementEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	topic := os.Getenv("KAFKA_TOPIC")
	return s.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)
}

// Send a post like
// (POST /posts/{post_id}/like)
func (s *ServerHandler) PostPostsPostIdLike(ctx echo.Context, postId int64) error {
	log.Println("Like post request")

	userId, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}

	log.Println("UserAuthorized successfully")

	return s.sendToKafka(EngagementEvent{
		EventID: uuid.New().ID(),
		Type:    Like,
		PostID:  postId,
		Author:  userId,
	})
}

// Send a post view
// (POST /posts/{post_id}/view)
func (s *ServerHandler) PostPostsPostIdView(ctx echo.Context, postId int64) error {
	log.Println("View post request")

	userId, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}

	log.Println("UserAuthorized successfully")

	return s.sendToKafka(EngagementEvent{
		EventID: uuid.New().ID(),
		Type:    View,
		PostID:  postId,
		Author:  userId,
	})
}
