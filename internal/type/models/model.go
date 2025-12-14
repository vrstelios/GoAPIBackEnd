package models

/*var Users = map[string]Login{} // Key is the username

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

*/

type APIError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request"`
}
