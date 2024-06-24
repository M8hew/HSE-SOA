package server

import (
	pb "content_service/proto"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDB struct {
}

func (db *MockDB) GetPostObj(id uint) (*PostIntRep, error) {
	return nil, nil
}

func (db *MockDB) CreatePost(post *PostIntRep) error {
	return nil
}

func (db *MockDB) UpdatePost(post *PostIntRep, id uint) error {
	return nil
}

func (db *MockDB) DeletePost(postAuthor int32, id uint) error {
	return nil
}

func (db *MockDB) GetPosts(offset int, batchSize int) ([]PostIntRep, error) {
	return nil, nil
}

func TestYourHandler(t *testing.T) {
	s := Server{dbWrapper: &MockDB{}}
	ctx := context.Background()

	t.Run("Post creation", func(t *testing.T) {
		resp, err := s.CreatePost(ctx, &pb.CreatePostRequest{
			Author:  0,
			Content: "Text",
		})

		assert.NoError(t, err, "Unexpected error")
		assert.Equal(t, uint64(0), resp.Id, "Content should match")
	})

	t.Run("Post deletion", func(t *testing.T) {
		_, err := s.DeletePost(ctx, &pb.DeletePostRequest{
			Id:     0,
			Author: 0,
		})

		assert.NoError(t, err, "Unexpected error")
	})
}
