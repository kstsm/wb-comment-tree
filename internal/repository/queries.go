package repository

const (
	CreateCommentQuery = `
		INSERT INTO comments (
		                      parent_id, 
		                      comment, 
		                      created_at
		                      )
		VALUES ($1, $2, NOW())
		RETURNING id, parent_id, comment, created_at
	`

	DeleteCommentQuery = `
		WITH RECURSIVE comment_tree AS (
    		SELECT id
    		FROM comments
    		WHERE id = $1
    		UNION ALL
    		SELECT c.id
    		FROM comments c
    		INNER JOIN comment_tree ct ON c.parent_id = ct.id
		)
		DELETE FROM comments
		WHERE id IN (SELECT id FROM comment_tree)
	`

	GetCommentsQuery = `
		SELECT id, parent_id, comment, created_at
		FROM comments
		WHERE ($1::bigint IS NULL OR parent_id = $1)
		AND ($2 = '' OR to_tsvector('simple', comment) @@ plainto_tsquery('simple', $2))
		ORDER BY 
			CASE WHEN $3 = 'created_at' AND $4 = 'asc' THEN created_at END,
			CASE WHEN $3 = 'created_at' AND $4 = 'desc' THEN created_at END DESC,
			created_at DESC
		LIMIT $5 OFFSET $6
			`

	GetCommentsCountQuery = `
		SELECT COUNT(*)
		FROM comments
		WHERE ($1::bigint IS NULL OR parent_id = $1)
		AND ($2 = '' OR to_tsvector('simple', comment) @@ plainto_tsquery('simple', $2))
	`

	GetAllCommentsQuery = `
		SELECT id, parent_id, comment, created_at
		FROM comments
		ORDER BY created_at DESC;
		`

	GetCommentByIDQuery = `
		SELECT id, parent_id, comment, created_at
		FROM comments
		WHERE id = $1
	`
)
