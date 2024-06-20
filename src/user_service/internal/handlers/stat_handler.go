package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"user_service/api"

	pb "user_service/api/proto"

	"github.com/labstack/echo/v4"
)

// Get total number of views and likes for a post
// (GET /posts/{postId}/stats)
func (s *ServerHandler) GetPostsPostIdStats(ctx echo.Context, postId int64) error {
	log.Println("Get top posts stats request")

	stat, err := s.statService.GetViewsLikes(context.Background(), &pb.GetViewsLikesRequest{Id: uint64(postId)})
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if stat == nil {
		return ctx.JSON(http.StatusOK, api.PostStat{
			PostId: &postId,
			Likes:  new(int),
			Views:  new(int),
		})
	}

	likes_ := int(stat.Likes)
	views_ := int(stat.Views)
	return ctx.JSON(http.StatusOK, api.PostStat{
		PostId: &postId,
		Likes:  &likes_,
		Views:  &views_,
	})
}

// Get top 5 posts by number of likes or views
// (GET /posts/top)
func (s *ServerHandler) GetPostsTop(ctx echo.Context, params api.GetPostsTopParams) error {
	log.Println("Get top5 posts request")

	var sortBy pb.GetTopPostsRequest_SortBy
	switch params.SortBy {
	case api.Likes:
		sortBy = pb.GetTopPostsRequest_LIKES
	case api.Views:
		sortBy = pb.GetTopPostsRequest_VIEWS
	}

	posts, err := s.statService.GetTopPosts(context.Background(), &pb.GetTopPostsRequest{SortBy: sortBy})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if posts == nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Stats not counted"})
	}

	postStats := make([]api.PostStat, 0, len(posts.Posts))
	for _, post := range posts.Posts {
		login, err := s.db.getUserLogin(int(post.PostAuthor))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error resolving user id into login"})
		}

		postId := int64(post.Id)
		likes_ := int(post.Likes)
		views_ := int(post.Views)

		stat := api.PostStat{
			Author: &login,
			PostId: &postId,
			Likes:  &likes_,
			Views:  &views_,
		}

		postStats = append(postStats, stat)
	}
	return ctx.JSON(http.StatusOK, postStats)
}

// Get top 3 users with the highest total likes
// (GET /users/top)
func (s *ServerHandler) GetUsersTop(ctx echo.Context) error {
	log.Println("Get top liked users request")

	users, err := s.statService.GetTopUsers(context.Background(), nil)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if users == nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error counting stats"})
	}

	fmt.Printf("top users: %v\n", users.TopUsers)
	userStats := make([]api.UserStat, 0, len(users.TopUsers))
	for _, user := range users.TopUsers {
		fmt.Printf("user_id: %v\n", user.UserId)
		login, err := s.db.getUserLogin(int(user.UserId))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error resolving user id into login"})
		}

		likes_ := int(user.Likes)

		userStat := api.UserStat{
			Login:      &login,
			TotalLikes: &likes_,
		}

		userStats = append(userStats, userStat)
	}
	return ctx.JSON(http.StatusOK, userStats)
}
