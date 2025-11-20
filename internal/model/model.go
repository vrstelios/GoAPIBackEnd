package model

type Login struct {
	HashedPassword string
	SessionToken   string
	//CSRFToken      string
	Role string
}

var Users = map[string]Login{} // Key is the username

var LibTasks map[string]*Task

type User = struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Task struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Urls struct {
	Urls []string `json:"urls"`
}

type DownloadImages struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
	Error   error  `json:"error"`
}

type QueryTasksRequest struct {
	Search string             `json:"search"`
	Paging QueryPagingRequest `json:"paging"`
	//Ordering int                `json:"ordering"`
}

type QueryTasksResponse struct {
	TotalResults int    `json:"totalResults"`
	TotalPages   int    `json:"totalPages"`
	Tasks        []Task `json:"tasks"`
}

type QueryPagingRequest struct {
	PageSize int `json:"pageSize" binding:"gt=0"`
	Page     int `json:"page" binding:"gte=0"`
}

type APIError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request"`
}
