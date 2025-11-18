package models

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentTree struct {
	Comment
	Children []CommentTree `json:"children,omitempty"`
}
