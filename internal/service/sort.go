package service

import (
	"github.com/kstsm/wb-comment-tree/internal/dto"
	"sort"
)

func (s *Service) sortTree(trees []dto.CommentResponse) {
	for i := range trees {
		if len(trees[i].Children) > 0 {
			s.sortTreeChildren(&trees[i])
		}
	}
}

func (s *Service) sortTreeChildren(node *dto.CommentResponse) {
	if len(node.Children) == 0 {
		return
	}

	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].CreatedAt < node.Children[j].CreatedAt
	})

	for i := range node.Children {
		s.sortTreeChildren(&node.Children[i])
	}
}

func (s *Service) sortRoots(roots []dto.CommentResponse, sortBy, order string) []dto.CommentResponse {
	if sortBy == "" {
		sortBy = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	sorted := make([]dto.CommentResponse, len(roots))
	copy(sorted, roots)

	sort.Slice(sorted, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "created_at", "id":
			less = sorted[i].CreatedAt < sorted[j].CreatedAt
		default:
			less = sorted[i].CreatedAt < sorted[j].CreatedAt
		}
		if order == "desc" {
			return !less
		}
		return less
	})

	return sorted
}
