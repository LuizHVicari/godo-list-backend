package project

import "time"

type CreateProjectRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type ListProjectsRequest struct {
	Name      *string `form:"name"`
	Sort      *string `form:"sort"`
	Direction *string `form:"direction"`
	Limit     *int32  `form:"limit"`
	Offset    *int32  `form:"offset"`
}

type ListProjectsResponse struct {
	Total   int64             `json:"total"`
	Results []ProjectResponse `json:"results"`
}

type ProjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toResponse(p *Project) ProjectResponse {
	return ProjectResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID.String(),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
