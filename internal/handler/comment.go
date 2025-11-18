package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/kstsm/wb-comment-tree/internal/apperrors"
	"github.com/kstsm/wb-comment-tree/internal/dto"
	"github.com/kstsm/wb-comment-tree/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Comment) == "" {
		utils.WriteError(w, http.StatusUnprocessableEntity, "content cannot be empty")
		return
	}

	comment, err := h.service.CreateComment(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrParentNotFound):
			utils.WriteError(w, http.StatusNotFound, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendJSON(w, http.StatusCreated, comment)
}

func (h *Handler) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid comment id")
		return
	}

	if err := h.service.DeleteComment(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrCommentNotFound):
			utils.WriteError(w, http.StatusNotFound, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "comment deleted successfully"})
}

func (h *Handler) getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetCommentsRequest

	query := r.URL.Query()

	if parentStr := query.Get("parent"); parentStr != "" {
		if parentID, err := strconv.ParseInt(parentStr, 10, 64); err == nil {
			req.ParentID = &parentID
		}
	}

	req.Page, req.PageSize = parsePagination(r)
	req.SortBy, req.Order = parseSorting(r, "created_at")

	req.Search = strings.TrimSpace(query.Get("search"))

	result, err := h.service.GetComments(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrCommentNotFound):
			utils.WriteError(w, http.StatusNotFound, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendJSON(w, http.StatusOK, result)
}
