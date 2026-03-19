package comment

import "time"

type CreateCommentRequest struct {
	ItemID  string `json:"item_id" binding:"required,uuid"`
	Content string `json:"content" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type ListCommentsRequest struct {
	Limit  *int32 `form:"limit"`
	Offset *int32 `form:"offset"`
}

type CommentResponse struct {
	ID        string    `json:"id"`
	ItemID    string    `json:"item_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListCommentsResponse struct {
	Total   int64             `json:"total"`
	Results []CommentResponse `json:"results"`
}

func toResponse(c *Comment) CommentResponse {
	return CommentResponse{
		ID:        c.ID.String(),
		ItemID:    c.ItemID.String(),
		AuthorID:  c.AuthorID.String(),
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
