package misc

import "GoAPIBackEnd/internal/type/models"

type QueryWorkoutsRequest struct {
	Search   string             `json:"search"`
	Paging   QueryPagingRequest `json:"paging"`
	Ordering int                `json:"ordering"`
}

type QueryWorkoutsResponse struct {
	TotalResults int               `json:"totalResults"`
	TotalPages   int               `json:"totalPages"`
	Workouts     []models.Workouts `json:"workouts"`
}

type QueryPagingRequest struct {
	PageSize int `json:"pageSize" binding:"gt=0"`
	Page     int `json:"page" binding:"gte=0"`
}
