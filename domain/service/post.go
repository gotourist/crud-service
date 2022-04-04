package service

import (
	"context"
	"fmt"
	"github.com/iman_task/crud-service/domain/entities"
	pb "github.com/iman_task/crud-service/genproto/post"
	loggerPkg "github.com/iman_task/crud-service/pkg/logger"
)

func (s *PostService) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	var pbPosts []*pb.Post

	posts, count, err := s.storage.Post().ListPosts(
		&entities.ListPostsRequest{
			Offset: (req.Page - 1) * req.Limit,
			Limit:  req.Limit,
		},
	)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get posts list from db"), loggerPkg.Error(err))
		return nil, err
	}

	for _, post := range posts {
		var pbPost pb.Post
		pbPost.Id = post.Id
		pbPost.Title = post.Title
		pbPost.Body = post.Body

		pbPosts = append(pbPosts, &pbPost)
	}

	return &pb.ListPostsResponse{
		Count:   count,
		Results: pbPosts,
	}, nil
}

func (s *PostService) DetailPost(ctx context.Context, req *pb.DetailPostRequest) (*pb.DetailPostResponse, error) {

	post, err := s.storage.Post().DetailPost(req.Id)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get post detail from db"), loggerPkg.Error(err))
		return nil, err
	}

	return &pb.DetailPostResponse{
		Result: &pb.Post{
			Id:    post.Id,
			Title: post.Title,
			Body:  post.Body,
		},
	}, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {

	tx, err := s.storage.Post().UpdatePost(
		&entities.UpdatePostRequest{
			Id:    req.Id,
			Title: req.Title,
			Body:  req.Body,
		},
	)
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		s.logger.Error(fmt.Sprintf("failed to update post in db"), loggerPkg.Error(err))
		return nil, err
	}

	post, err := s.storage.Post().DetailPost(req.Id)
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return nil, err
	}

	err = s.publishUpdatePostMessage(post)
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return nil, err
	}

	if tx != nil {
		tx.Commit()
	}

	return &pb.UpdatePostResponse{
		Result: &pb.Post{
			Id:    post.Id,
			Title: post.Title,
			Body:  post.Body,
		},
	}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {

	tx, err := s.storage.Post().DeletePost(req.Id)
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		s.logger.Error(fmt.Sprintf("failed to get posts list from db"), loggerPkg.Error(err))
		return nil, err
	}

	post, err := s.storage.Post().DetailPost(req.Id)
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return nil, err
	}
	post.IsDeleted = true

	err = s.publishUpdatePostMessage(post)
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return nil, err
	}

	if tx != nil {
		tx.Commit()
	}

	return &pb.DeletePostResponse{
		Code:   0,
		Errors: nil,
	}, nil
}
