package service

import (
	"context"
	"github.com/kstsm/wb-comment-tree/internal/apperrors"
	"github.com/kstsm/wb-comment-tree/internal/dto"
	"time"
)

func (s *Service) CreateComment(ctx context.Context, req dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	comment, err := s.repo.CreateComment(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.CommentResponse{
		ID:        comment.ID,
		ParentID:  comment.ParentID,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) DeleteComment(ctx context.Context, id int64) error {
	return s.repo.DeleteComment(ctx, id)
}

func (s *Service) GetComments(ctx context.Context, req dto.GetCommentsRequest) (interface{}, error) {
	if req.Search != "" {
		comments, total, err := s.repo.GetComments(ctx, req)
		if err != nil {
			return nil, err
		}

		responses := make([]dto.CommentResponse, 0, len(comments))
		for _, comment := range comments {
			responses = append(responses, dto.CommentResponse{
				ID:        comment.ID,
				ParentID:  comment.ParentID,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt.Format(time.RFC3339),
			})
		}

		return dto.PaginatedResponse{
			Comments: responses,
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		}, nil
	}

	if req.ParentID != nil {
		_, err := s.repo.GetCommentByID(ctx, *req.ParentID)
		if err != nil {
			return nil, apperrors.ErrCommentNotFound
		}

		allComments, err := s.repo.GetAllComments(ctx)
		if err != nil {
			return nil, err
		}

		tree := s.buildTreeFromParent(allComments, *req.ParentID)
		if tree == nil {
			return nil, apperrors.ErrCommentNotFound
		}

		return []dto.CommentResponse{*tree}, nil
	}

	allComments, err := s.repo.GetAllComments(ctx)
	if err != nil {
		return nil, err
	}

	trees := s.buildTree(allComments)
	trees = s.sortRoots(trees, req.SortBy, req.Order)
	total := len(trees)

	if req.Page > 0 && req.PageSize > 0 {
		start := (req.Page - 1) * req.PageSize
		end := start + req.PageSize
		if start < len(trees) {
			if end > len(trees) {
				end = len(trees)
			}
			trees = trees[start:end]
		} else {
			trees = []dto.CommentResponse{}
		}
	}

	return map[string]interface{}{
		"comments":  trees,
		"page":      req.Page,
		"page_size": req.PageSize,
		"total":     total,
	}, nil
}
