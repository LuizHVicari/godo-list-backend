package item

import "time"

type CreateItemRequest struct {
	StepID      string       `json:"step_id" binding:"required,uuid"`
	Name        string       `json:"name" binding:"required"`
	Description *string      `json:"description"`
	Priority    ItemPriority `json:"priority"`
	Position    *int32       `json:"position"`
}

type UpdateItemRequest struct {
	StepID      *string      `json:"step_id"`
	Name        string       `json:"name" binding:"required"`
	Description *string      `json:"description"`
	Priority    ItemPriority `json:"priority"`
	Position    int32        `json:"position" binding:"required"`
	IsCompleted bool         `json:"is_completed"`
}

type ItemRepositionItem struct {
	ID       string `json:"id" binding:"required,uuid"`
	Position int32  `json:"position" binding:"required"`
}

type RepositionItemsRequest struct {
	StepID string               `json:"step_id" binding:"required,uuid"`
	Items  []ItemRepositionItem `json:"items" binding:"required,min=1"`
}

type ListItemsRequest struct {
	Name      *string `form:"name"`
	Sort      *string `form:"sort"`
	Direction *string `form:"direction"`
	Limit     *int32  `form:"limit"`
	Offset    *int32  `form:"offset"`
}

type ItemResponse struct {
	ID          string       `json:"id"`
	StepID      string       `json:"step_id"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	Priority    ItemPriority `json:"priority"`
	Position    int32        `json:"position"`
	IsCompleted bool         `json:"is_completed"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type ListItemsResponse struct {
	Total   int64          `json:"total"`
	Results []ItemResponse `json:"results"`
}

func toResponse(i *Item) ItemResponse {
	return ItemResponse{
		ID:          i.ID.String(),
		StepID:      i.StepID.String(),
		Name:        i.Name,
		Description: i.Description,
		Priority:    i.Priority,
		Position:    i.Position,
		IsCompleted: i.IsCompleted,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}
