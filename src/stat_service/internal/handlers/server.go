package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "stat_service/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	dbWrapper
	pb.UnimplementedStatServiceServer
}

func NewServer() *Server {
	db, err := newDBWrapper()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	log.Println("Db instance successfully created")
	return &Server{dbWrapper: db}
}

func (s *Server) ListenAndServe(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStatServiceServer(grpcServer, s)
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) GetViewsLikes(ctx context.Context, in *pb.GetViewsLikesRequest) (*pb.GetViewsLikesResponse, error) {
	select {
	case <-ctx.Done():
		return &pb.GetViewsLikesResponse{}, ctx.Err()
	default:
		likes, views, err := s.dbWrapper.countLikesViews(in.Id)
		if err != nil {
			return &pb.GetViewsLikesResponse{}, err
		}
		return &pb.GetViewsLikesResponse{Views: views, Likes: likes}, nil
	}
}

func (s *Server) GetTopPosts(ctx context.Context, in *pb.GetTopPostsRequest) (*pb.GetTopPostsResponse, error) {
	select {
	case <-ctx.Done():
		return &pb.GetTopPostsResponse{}, ctx.Err()
	default:
		var sortBy eventType

		switch in.SortBy {
		case pb.GetTopPostsRequest_VIEWS:
			sortBy = views
		case pb.GetTopPostsRequest_LIKES:
			sortBy = likes
		}
		posts, err := s.dbWrapper.getTopPosts(sortBy)

		if err != nil {
			return &pb.GetTopPostsResponse{}, err
		}

		pbPosts := make([]*pb.GetTopPostsResponse_Post, 0, len(posts))
		for _, post := range posts {
			pbPosts = append(pbPosts, &pb.GetTopPostsResponse_Post{
				Id:         post.id,
				PostAuthor: post.post_author,
				Views:      post.views,
				Likes:      post.likes,
			})
		}
		return &pb.GetTopPostsResponse{Posts: pbPosts}, nil
	}
}

func (s *Server) GetTopUsers(ctx context.Context, _ *emptypb.Empty) (*pb.GetTopUsersResponse, error) {
	select {
	case <-ctx.Done():
		return &pb.GetTopUsersResponse{}, ctx.Err()
	default:
		userInfo, err := s.dbWrapper.getTopUsers()
		if err != nil {
			return &pb.GetTopUsersResponse{}, err
		}

		pbUserStat := make([]*pb.GetTopUsersResponse_UserStat, 0, len(userInfo))
		for _, user := range userInfo {
			pbUserStat = append(pbUserStat, &pb.GetTopUsersResponse_UserStat{
				UserId: user.postAuthor,
				Likes:  user.likes,
			})
		}
		return &pb.GetTopUsersResponse{TopUsers: pbUserStat}, nil
	}
}
