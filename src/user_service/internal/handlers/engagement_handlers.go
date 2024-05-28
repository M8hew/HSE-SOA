package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	pb "user_service/api/proto"

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
	EventID     uint32    `json:"id"`
	Type        EventType `json:"event_type"`
	PostID      int64     `json:"post_id"`
	EventAuthor int       `json:"event_author"`
	PostAuthor  int32     `json:"post_author"`
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

func (s *ServerHandler) engagementEventHandler(ctx echo.Context, postId int64, e EventType) error {
	userId, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}

	log.Println("UserAuthorized successfully")

	post, err := s.contentService.GetPost(context.Background(), &pb.GetPostRequest{
		Id:     uint64(postId),
		Author: int32(userId),
	})
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	if post == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
	}

	return s.sendToKafka(EngagementEvent{
		EventID:     uuid.New().ID(),
		Type:        e,
		PostID:      postId,
		EventAuthor: userId,
		PostAuthor:  post.Author,
	})
}

// Send a post like
// (POST /posts/{post_id}/like)
func (s *ServerHandler) PostPostsPostIdLike(ctx echo.Context, postId int64) error {
	log.Println("Like post request")

	return s.engagementEventHandler(ctx, postId, Like)
}

// Send a post view
// (POST /posts/{post_id}/view)
func (s *ServerHandler) PostPostsPostIdView(ctx echo.Context, postId int64) error {
	log.Println("View post request")

	return s.engagementEventHandler(ctx, postId, View)
}
