package service

import (
	"github.com/kstsm/wb-comment-tree/internal/dto"
	"github.com/kstsm/wb-comment-tree/internal/models"
	"time"
)

func (s *Service) buildTree(comments []models.Comment) []dto.CommentResponse {
	commentMap := make(map[int64]*dto.CommentResponse)
	var rootIDs []int64

	for _, comment := range comments {
		commentMap[comment.ID] = &dto.CommentResponse{
			ID:        comment.ID,
			ParentID:  comment.ParentID,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
			Children:  []dto.CommentResponse{},
		}
	}

	for _, comment := range comments {
		if comment.ParentID == nil {
			rootIDs = append(rootIDs, comment.ID)
		} else {

			if parent, exists := commentMap[*comment.ParentID]; exists {
				node := commentMap[comment.ID]

				parent.Children = append(parent.Children, *node)
			}
		}
	}

	var copyNode func(int64) dto.CommentResponse
	copyNode = func(nodeID int64) dto.CommentResponse {
		node := commentMap[nodeID]
		copied := dto.CommentResponse{
			ID:        node.ID,
			ParentID:  node.ParentID,
			Comment:   node.Comment,
			CreatedAt: node.CreatedAt,
			Children:  make([]dto.CommentResponse, 0, len(node.Children)),
		}
		for _, child := range node.Children {
			copied.Children = append(copied.Children, copyNode(child.ID))
		}
		return copied
	}

	roots := make([]dto.CommentResponse, 0, len(rootIDs))
	for _, id := range rootIDs {
		roots = append(roots, copyNode(id))
	}

	s.sortTree(roots)

	return roots
}

func (s *Service) buildTreeFromParent(comments []models.Comment, parentID int64) *dto.CommentResponse {
	commentMap := make(map[int64]*dto.CommentResponse)

	descendantIDs := s.getAllDescendants(comments, parentID)
	descendantIDs[parentID] = true

	for _, comment := range comments {
		if descendantIDs[comment.ID] {
			commentMap[comment.ID] = &dto.CommentResponse{
				ID:        comment.ID,
				ParentID:  comment.ParentID,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt.Format(time.RFC3339),
				Children:  []dto.CommentResponse{},
			}
		}
	}

	parentNode, exists := commentMap[parentID]
	if !exists {
		return nil
	}

	childrenMap := make(map[int64][]int64)
	for _, comment := range comments {
		if !descendantIDs[comment.ID] {
			continue
		}
		if comment.ParentID != nil {
			parentID := *comment.ParentID
			if descendantIDs[parentID] {
				childrenMap[parentID] = append(childrenMap[parentID], comment.ID)
			}
		}
	}

	var buildNode func(int64) dto.CommentResponse
	buildNode = func(nodeID int64) dto.CommentResponse {
		node := commentMap[nodeID]
		result := dto.CommentResponse{
			ID:        node.ID,
			ParentID:  node.ParentID,
			Comment:   node.Comment,
			CreatedAt: node.CreatedAt,
			Children:  make([]dto.CommentResponse, 0),
		}
		if childIDs, exists := childrenMap[nodeID]; exists {
			for _, childID := range childIDs {
				result.Children = append(result.Children, buildNode(childID))
			}
		}
		return result
	}

	if parentNode != nil {
		result := buildNode(parentID)
		s.sortTreeChildren(&result)
		return &result
	}

	return nil
}

func (s *Service) getAllDescendants(comments []models.Comment, parentID int64) map[int64]bool {
	descendants := make(map[int64]bool)
	commentMap := make(map[int64]*models.Comment)

	for i := range comments {
		commentMap[comments[i].ID] = &comments[i]
	}

	var findDescendants func(id int64)
	findDescendants = func(id int64) {
		for _, comment := range comments {
			if comment.ParentID != nil && *comment.ParentID == id {
				descendants[comment.ID] = true
				findDescendants(comment.ID)
			}
		}
	}

	findDescendants(parentID)
	return descendants
}
