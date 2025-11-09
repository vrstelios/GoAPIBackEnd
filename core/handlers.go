package core

import (
	"GoAPIBackEnd/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"strings"
	"time"
)

// @Summary Create Task.
// @Description Notes:<br>Leave the "relations" sub-entity empty.
// @Tags Task
// @Produce json
// @Param request body model.Task true "Task"
// @Success 200 {object} model.APIError
// @Failure 400 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /tasks [post]
func Post(ctx *gin.Context) {
	task := model.Task{}
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(task.Id) == 0 {
		uid, err := uuid.NewUUID()
		if err != nil {
			return
		}
		task.Id = uid.String()
	}

	for _, t := range model.LibTasks {
		if t.Id == task.Id {
			ctx.JSON(http.StatusAlreadyReported, gin.H{"error": "Same Id"})
			return
		}
	}

	model.LibTasks[task.Id] = &task
	ctx.JSON(http.StatusCreated, model.LibTasks[task.Id])
}

// @Summary Get Task.
// @Tags Task
// @Produce json
// @Param id path string true "Task Id"
// @Success 200 {object} model.Task
// @Failure 404 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /tasks/{id} [get]
func Get(ctx *gin.Context) {
	id := ctx.Param("id")

	task, ok := model.LibTasks[id]
	if !ok {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// @Summary	Query Task.
// @Tags		Task
// @Produce	json
// @Param		id		query		string	false	"Task Id"
// @Success	200		{array}		model.Task
// @Failure	500		{object}	model.APIError
// @Router		/tasks [get]
func Query(ctx *gin.Context) {
	libTasks := model.LibTasks
	if len(libTasks) == 0 {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	var tasks []model.Task
	for _, task := range libTasks {
		tasks = append(tasks, *task)
	}

	ctx.JSON(http.StatusOK, tasks)
}

// @Summary Update Task.
// @Tags Task
// @Produce json
// @Param id path string false "Task"
// @Param request body model.Task true "Task"
// @Success 204
// @Failure 400 {object} model.APIError
// @Failure 404
// @Failure 500 {object} model.APIError
// @Router /tasks/{id} [put]
func Put(ctx *gin.Context) {
	task := model.Task{}
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := ctx.Param("id")
	if task.Id != "" && task.Id != id {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id in body must match URL"})
		return
	}

	t, ok := model.LibTasks[id]
	if !ok {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	//t.Id = task.Id
	t.Title = task.Title
	t.Done = task.Done

	ctx.JSON(http.StatusOK, t)
}

// @Summary Delete Task.
// @Tags Task
// @Produce json
// @Param id path string false "Task"
// @Success 204
// @Failure 404
// @Failure 500 {object} model.APIError
// @Router /tasks/{id} [delete]
func Del(ctx *gin.Context) {
	id := ctx.Param("id")

	_, ok := model.LibTasks[id]
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	delete(model.LibTasks, id)
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// @Summary	Download images
// @Tags		Task
// @Produce	json
// @Param		request	body		model.Urls	true	"Urls"
// @Success	200		{object}	model.DownloadImages
// @Failure	400		{object}	model.APIError
// @Failure	500		{object}	model.APIError
// @Router		/download/images [post]
func DownloadUrls(ctx *gin.Context) {
	urls := model.Urls{}
	err := ctx.ShouldBindJSON(&urls)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resultCh := make(chan model.DownloadImages, len(urls.Urls))
	for _, url := range urls.Urls {
		go downLoadImage(url, resultCh)
	}

	timeout := time.After(10 * time.Second)

	results := make([]model.DownloadImages, 0, len(urls.Urls))
	for i := 0; i < len(urls.Urls); i++ {
		select {
		case result := <-resultCh:
			results = append(results, result)
		case <-timeout:
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Timeout waiting for downloads",
				"results": results,
			})
			return
		}

	}

	ctx.JSON(http.StatusOK, gin.H{"results": results})
}

// @Summary	Query Task (server side filtering & paging)
// @Tags		Task
// @Produce	json
// @Param		request	body		model.QueryTasksRequest	true	"QueryTasksRequest"
// @Success	200		{object}	model.QueryTasksResponse
// @Failure	400		{object}	model.APIError
// @Failure	500		{object}	model.APIError
// @Router		/tasks/query [post]
func QueryTasksV2(ctx *gin.Context) {
	queryRequest := model.QueryTasksRequest{}
	err := ctx.ShouldBindJSON(&queryRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var filtered []model.Task
	for _, task := range model.LibTasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(queryRequest.Search)) {
			filtered = append(filtered, *task)
		}
	}

	totalRows := len(filtered)
	if totalRows == 0 {
		ctx.JSON(http.StatusOK, model.QueryTasksResponse{TotalResults: totalRows, TotalPages: 0, Tasks: []model.Task{}})
		return
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(queryRequest.Paging.PageSize)))

	startPage := (queryRequest.Paging.Page - 1) * queryRequest.Paging.PageSize
	endPage := startPage + queryRequest.Paging.PageSize

	switch {
	case startPage > totalRows:
		startPage = totalRows
	case endPage > totalRows:
		endPage = totalRows
	}

	pagedTasks := filtered[startPage:endPage]

	ctx.JSON(http.StatusOK, model.QueryTasksResponse{
		TotalResults: totalRows,
		TotalPages:   totalPages,
		Tasks:        pagedTasks,
	})
}

// @Summary Register a new user.
// @Description Creates a new user account with username and password
// @Tags auth
// @Produce json
// @Param user body model.User true "User registration data"
// @Success 200 {object} model.APIError
// @Failure 400 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /auth/register [post]
func Register(ctx *gin.Context) {
	user := model.User{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(user.Username) < 0 || len(user.Password) < 0 {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid username/password"})
		return
	}

	if _, ok := model.Users[user.Username]; ok {
		ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, _ := hashPassword(user.Password)
	model.Users[user.Username] = model.Login{
		HashedPassword: hashedPassword,
	}

	ctx.JSON(http.StatusTemporaryRedirect, "User registered successfully!")
}

// @Summary Login user
// @Description Authenticates user and returns session cookies
// @Tags auth
// @Produce json
// @Param user body model.User true "User registration data"
// @Success 200 {object} model.APIError
// @Failure 400 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /auth/login [post]
func Login(ctx *gin.Context) {
	userPayload := model.User{}
	err := ctx.ShouldBindJSON(&userPayload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, ok := model.Users[userPayload.Username]
	if !ok || !checkPasswordHash(userPayload.Password, user.HashedPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	sessionToken := generateToken(32)
	csrToken := generateToken(32)

	// Set session cookie
	ctx.SetCookie("session_token", sessionToken, 3600, "/", "", false, true)

	// Set CSRF token in a cookie
	ctx.SetCookie("csr_token", csrToken, 3600, "/", "", false, false)

	// Store tokens in the database
	user.SessionToken = sessionToken
	user.CSRFToken = csrToken
	model.Users[userPayload.Username] = user

	ctx.JSON(http.StatusTemporaryRedirect, "Login successfully!")
}

func Protected(ctx *gin.Context) {
	if err := Authorize(ctx); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username := ctx.GetString("username")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "CSRF validation successful! Welcome, " + username,
	})

}

// @Summary	Logout user.
// @Tags		auth
// @Produce	json
// @Success	200		{array}		model.Task
// @Failure	500		{object}	model.APIError
// @Router		/auth/logout [get]
func Logout(ctx *gin.Context) {
	if err := Authorize(ctx); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username := ctx.GetString("username")

	// Clear cookie
	ctx.SetCookie("session_token", "", -1, "/", "", false, true)
	ctx.SetCookie("csr_token", "", -1, "/", "", false, false)

	// Clear the tokens from the database
	user := model.Users[username]
	user.SessionToken = ""
	user.CSRFToken = ""
	model.Users[username] = user

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successfully!"})
}

func downLoadImage(url string, ch chan<- model.DownloadImages) {
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		ch <- model.DownloadImages{
			URL:     url,
			Success: false,
			Error:   fmt.Errorf("HTTP error: %v", err),
		}
		return
	}
	defer resp.Body.Close()

	ch <- model.DownloadImages{
		URL:     url,
		Success: true,
		Error:   nil,
	}
}
