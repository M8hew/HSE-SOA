package server

import (
	pb "content_service/proto"
	"context"
	"fmt"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type Server struct {
	dbWrapper
	pb.UnimplementedContentServiceServer
}

func NewServer() *Server {
	db, err := NewDBWrapper()
	if err != nil {
		return nil
	}
	return &Server{dbWrapper: db}
}

func (s *Server) ListenAndServe(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterContentServiceServer(grpcServer, s)
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) CreatePost(ctx context.Context, in *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	select {
	case <-ctx.Done():
		return &pb.CreatePostResponse{}, ctx.Err()
	default:
		post := PostIntRep{Author: in.Author, Content: in.Content}
		err := s.dbWrapper.CreatePost(&post)
		if err != nil {
			return &pb.CreatePostResponse{}, fmt.Errorf("error creating post: %v", err)
		}
		return &pb.CreatePostResponse{Id: uint64(post.ID)}, nil
	}
}

func (s *Server) UpdatePost(ctx context.Context, in *pb.UpdatePostRequest) (*empty.Empty, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		post := PostIntRep{Author: in.Author, Content: in.Content}
		err := s.dbWrapper.UpdatePost(&post, uint(in.Id))
		return nil, err
	}
}

func (s *Server) DeletePost(ctx context.Context, in *pb.DeletePostRequest) (*empty.Empty, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, s.dbWrapper.DeletePost(&in.Author, uint(in.Id))
	}
}

func (s *Server) GetPost(ctx context.Context, in *pb.GetPostRequest) (*pb.Post, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		obj, err := s.dbWrapper.GetPostObj(uint(in.Id), in.Author)
		if err != nil {
			return nil, err
		}
		return &pb.Post{
			Id:      uint64(obj.ID),
			Author:  obj.Author,
			Content: obj.Content,
		}, nil
	}
}

func (s *Server) ListPosts(ctx context.Context, in *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		objList, err := s.dbWrapper.GetPosts(int(in.Offset), int(in.MaxCnt), in.Author)
		if err != nil {
			return nil, err
		}

		var result []*pb.Post
		for _, obj := range objList {
			result = append(result, &pb.Post{
				Id:      uint64(obj.ID),
				Author:  obj.Author,
				Content: obj.Content,
			})
		}
		return &pb.ListPostsResponse{
			Posts:      result,
			LastOffset: in.Offset + int32(len(result)),
		}, nil
	}
}