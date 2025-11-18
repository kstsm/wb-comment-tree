package apperrors

import "errors"

var (
	ErrParentNotFound  = errors.New("parent not found")
	ErrCommentNotFound = errors.New("comment not found")
)
