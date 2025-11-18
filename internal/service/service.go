package service

import (
	"context"
	"github.com/kstsm/wb-comment-tree/internal/dto"
	"github.com/kstsm/wb-comment-tree/internal/repository"
)

type ServiceI interface {
	CreateComment(ctx context.Context, req dto.CreateCommentRequest) (*dto.CommentResponse, error)
	DeleteComment(ctx context.Context, id int64) error
	GetComments(ctx context.Context, req dto.GetCommentsRequest) (interface{}, error)
}

type Service struct {
	repo repository.RepositoryI
}

func NewService(repo repository.RepositoryI) *Service {
	return &Service{
		repo: repo,
	}
}
