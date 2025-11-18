package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func parsePagination(r *http.Request) (page, pageSize int) {
	page = 1
	pageSize = 10

	if v := r.URL.Query().Get("page"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			page = i
		}
	}
	if v := r.URL.Query().Get("page_size"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			pageSize = i
		}
	}

	return
}

func parseSorting(r *http.Request, defaultSort string) (sortBy, order string) {
	sortBy = r.URL.Query().Get("sort_by")
	if sortBy != "created_at" && sortBy != "updated_at" {
		sortBy = defaultSort
	}

	order = strings.ToLower(r.URL.Query().Get("order"))
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return
}
