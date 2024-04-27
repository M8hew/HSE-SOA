package handlers

import (
	"context"
	"log"
	"net/http"
	"user_service/api"

	pb "user_service/api/proto"

	"github.com/labstack/echo/v4"
)

// Create new post
// (POST /posts)
func (s *ServerHandler) PostPosts(ctx echo.Context) error {
	userID, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}
	log.Println("UserAuthorized")

	reqBody := api.PostPostsJSONRequestBody{}
	if err = ctx.Bind(&reqBody); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := s.contentService.CreatePost(context.Background(), &pb.CreatePostRequest{
		Author:  int32(userID),
		Content: *reqBody.Content,
	})

	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Cannot create task"})
	}
	respID := int64(resp.Id)
	return ctx.JSON(http.StatusCreated, struct{ Id *int64 }{&respID})
}

// Delete post
// (DELETE /posts/{post_id})
func (s *ServerHandler) DeletePostsPostId(ctx echo.Context, postId int64) error {
	userID, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}
	log.Println("UserAuthorized")

	_, err = s.contentService.DeletePost(context.Background(), &pb.DeletePostRequest{Id: uint64(postId), Author: int32(userID)})
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusNoContent, map[string]string{"message": "Post deleted successfully"})
}

// Update post
// (PUT /posts/{post_id})
func (s *ServerHandler) PutPostsPostId(ctx echo.Context, postId int64) error {
	userID, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}
	log.Println("UserAuthorized")

	reqBody := api.PutPostsPostIdJSONRequestBody{}
	if err = ctx.Bind(&reqBody); err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	_, err = s.contentService.UpdatePost(context.Background(), &pb.UpdatePostRequest{
		Id:      uint64(postId),
		Author:  int32(userID),
		Content: *reqBody.Content,
	})
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Post updated successfully"})
}

// Retrieve a post by ID
// (GET /posts/{post_id})
func (s *ServerHandler) GetPostsPostId(ctx echo.Context, postId int64) error {
	userID, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}
	log.Println("UserAuthorized")

	post, err := s.contentService.GetPost(context.Background(), &pb.GetPostRequest{
		Id:     uint64(postId),
		Author: int32(userID),
	})
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	if post == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
	}

	userLogin, err := s.db.getUserLogin(userID)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Cannot resolve userID into userLogin"})
	}
	return ctx.JSON(http.StatusOK, api.Post{
		Id:      &postId,
		Author:  &userLogin,
		Content: &post.Content,
	})

}

// Retrieve list of posts with pagination
// (GET /posts)
func (s *ServerHandler) GetPosts(ctx echo.Context, params api.GetPostsParams) error {
	userID, err := checkAuth(ctx, s.keys)
	if err != nil {
		return err
	}
	log.Println("UserAuthorized")

	posts, err := s.contentService.ListPosts(context.Background(), &pb.ListPostsRequest{
		Offset: int32(params.FirstId),
		MaxCnt: int32(params.MaxPosts),
		Author: int32(userID),
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	userLogin, err := s.db.getUserLogin(userID)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Cannot resolve userID into userLogin"})
	}

	postsList := make([]api.Post, 0, len(posts.Posts))
	for _, post := range posts.Posts {
		postID := int64(post.Id)
		postsList = append(postsList, api.Post{
			Id:      &postID,
			Author:  &userLogin,
			Content: &post.Content,
		})
	}

	return ctx.JSON(http.StatusOK, postsList)
}
