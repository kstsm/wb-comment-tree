package dto

type CreateCommentRequest struct {
	ParentID *int64 `json:"parent_id,omitempty"`
	Comment  string `json:"comment"`
}

type GetCommentsRequest struct {
	ParentID *int64 `json:"parent_id,omitempty" query:"parent"`
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	Order    string `json:"order" query:"order"`
	Search   string `json:"search" query:"search"`
}

type CommentResponse struct {
	ID        int64             `json:"id"`
	ParentID  *int64            `json:"parent_id,omitempty"`
	Comment   string            `json:"comment"`
	CreatedAt string            `json:"created_at"`
	Children  []CommentResponse `json:"children,omitempty"`
}

type PaginatedResponse struct {
	Comments []CommentResponse `json:"comments"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int               `json:"total"`
}
