package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/wb-comment-tree/internal/dto"
	"github.com/kstsm/wb-comment-tree/internal/models"
)

type RepositoryI interface {
	CreateComment(ctx context.Context, req dto.CreateCommentRequest) (*models.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	GetComments(ctx context.Context, req dto.GetCommentsRequest) ([]models.Comment, int, error)
	GetAllComments(ctx context.Context) ([]models.Comment, error)
	GetCommentByID(ctx context.Context, id int64) (*models.Comment, error)
}

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) RepositoryI {
	return &Repository{
		conn: conn,
	}
}
