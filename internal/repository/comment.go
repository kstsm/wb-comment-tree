package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kstsm/wb-comment-tree/internal/apperrors"
	"github.com/kstsm/wb-comment-tree/internal/dto"
	"github.com/kstsm/wb-comment-tree/internal/models"
)

func (r *Repository) CreateComment(ctx context.Context, req dto.CreateCommentRequest) (*models.Comment, error) {
	var comment models.Comment

	err := r.conn.QueryRow(ctx, CreateCommentQuery, req.ParentID, req.Comment).Scan(
		&comment.ID,
		&comment.ParentID,
		&comment.Comment,
		&comment.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return nil, apperrors.ErrParentNotFound
			}
		}

		return nil, fmt.Errorf("QueryRow-CreateComment: %w", err)
	}

	return &comment, nil
}

func (r *Repository) DeleteComment(ctx context.Context, id int64) error {
	result, err := r.conn.Exec(ctx, DeleteCommentQuery, id)
	if err != nil {
		return fmt.Errorf("Exec-DeleteComment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrCommentNotFound
	}

	return nil
}

func (r *Repository) GetComments(ctx context.Context, req dto.GetCommentsRequest) ([]models.Comment, int, error) {
	var comments []models.Comment
	var total int

	if req.Search != "" {
		req.ParentID = nil
	}

	var parentID *int64 = req.ParentID
	var search *string
	if req.Search != "" {
		search = &req.Search
	}

	offset := (req.Page - 1) * req.PageSize

	err := r.conn.QueryRow(ctx, GetCommentsCountQuery, parentID, search).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("QueryRow-GetCommentsCount: %w", err)
	}

	rows, err := r.conn.Query(ctx, GetCommentsQuery, parentID, search, req.SortBy, req.Order, req.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("Query-GetComments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Comment,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("Query-GetComments: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("Rows-GetComments: %w", err)
	}

	return comments, total, nil
}

func (r *Repository) GetCommentByID(ctx context.Context, id int64) (*models.Comment, error) {
	var comment models.Comment

	err := r.conn.QueryRow(ctx, GetCommentByIDQuery, id).Scan(
		&comment.ID,
		&comment.ParentID,
		&comment.Comment,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("QueryRow-GetCommentByID: %w", err)
	}

	return &comment, nil
}

func (r *Repository) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	var comments []models.Comment

	rows, err := r.conn.Query(ctx, GetAllCommentsQuery)
	if err != nil {
		return nil, fmt.Errorf("Query-GetAllComments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Comment,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Query-GetAllComments: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Rows-GetAllComments: %w", err)
	}

	return comments, nil
}
