package step

import "time"

type CreateStepRequest struct {
	Name     string `json:"name" binding:"required"`
	Position *int32 `json:"position"`
}

type UpdateStepRequest struct {
	Name        string `json:"name" binding:"required"`
	Position    int32  `json:"position" binding:"required"`
	IsCompleted bool   `json:"is_completed"`
}

type StepRepositionItem struct {
	ID       string `json:"id" binding:"required,uuid"`
	Position int32  `json:"position" binding:"required"`
}

type RepositionStepsRequest struct {
	Steps []StepRepositionItem `json:"steps" binding:"required,min=1"`
}

type ListStepsRequest struct {
	Name      *string `form:"name"`
	Sort      *string `form:"sort"`
	Direction *string `form:"direction"`
	Limit     *int32  `form:"limit"`
	Offset    *int32  `form:"offset"`
}

type StepResponse struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Name        string    `json:"name"`
	Position    int32     `json:"position"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListStepsResponse struct {
	Total   int64          `json:"total"`
	Results []StepResponse `json:"results"`
}

func toResponse(s *Step) StepResponse {
	return StepResponse{
		ID:          s.ID.String(),
		ProjectID:   s.ProjectID.String(),
		Name:        s.Name,
		Position:    s.Position,
		IsCompleted: s.IsCompleted,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
