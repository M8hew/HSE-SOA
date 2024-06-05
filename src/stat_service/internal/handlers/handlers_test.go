package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	pb "stat_service/proto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDBWrapper struct {
}

func (db MockDBWrapper) countLikesViews(id uint64) (likes_, views_ int64, err error) {
	return 0, 0, nil
}

func (db MockDBWrapper) getTopPosts(e eventType) ([]post, error) {
	return []post{}, nil
}

func (db MockDBWrapper) getTopUsers() ([]userInfo, error) {
	return []userInfo{}, nil
}

func TestLivelinessHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/liveliness", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(liveHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "OK", rr.Body.String(), "handler returned unexpected body")
}

func TestGRPCHandlers(t *testing.T) {
	s := Server{
		dbWrapper: MockDBWrapper{},
	}

	// GetViewsLikes test
	resp, err := s.GetViewsLikes(context.Background(), &pb.GetViewsLikesRequest{Id: 1})
	require.NoError(t, err)
	assert.Equal(t, int64(0), resp.Likes)
	assert.Equal(t, int64(0), resp.Views)

	// GetTopPosts test
	resp2, err := s.GetTopPosts(context.Background(), &pb.GetTopPostsRequest{SortBy: pb.GetTopPostsRequest_LIKES})
	require.NoError(t, err)
	assert.Equal(t, 0, len(resp2.Posts))

	// GetTopUsers test
	resp3, err := s.GetTopUsers(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, 0, len(resp3.TopUsers))
}
