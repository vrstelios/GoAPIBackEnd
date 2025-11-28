package server

import (
	"GoAPIBackEnd/internal/apperrors"
	"GoAPIBackEnd/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// @Summary Create Task.
// @Description Notes:<br>Leave the "relations" sub-entity empty.
// @Tags Task
// @Produce json
// @Param request body model.Task true "Task"
// @Success 200 {object} model.APIError
// @Failure 208 {object} model.APIError
// @Failure 400 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /tasks [post]
func Post(ctx *gin.Context) {
	task := model.Task{}
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
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
			apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusAlreadyReported, "DUPLICATE_ID", fmt.Errorf("Task with id %s already exists", task.Id)))
			return
		}
	}

	model.LibTasks[task.Id] = &task
	apperrors.GetAPIError(ctx, model.LibTasks[task.Id], http.StatusOK, nil)
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
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "TASK_NOT_FOUND", fmt.Errorf("Task with id %s not found", id)))
		return
	}

	apperrors.GetAPIError(ctx, task, http.StatusOK, nil)
}

// @Summary	Query Task.
// @Tags		Task
// @Produce	json
// @Param		id		query		string	false	"Task Id"
// @Success	200		{array}		model.Task
// @Failure 404     {object}    model.APIError
// @Failure	500		{object}	model.APIError
// @Router		/tasks [get]
func Query(ctx *gin.Context) {
	libTasks := model.LibTasks
	if len(libTasks) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "TASK_NOT_FOUND", fmt.Errorf("Task not found")))
		return
	}

	var tasks []model.Task
	for _, task := range libTasks {
		tasks = append(tasks, *task)
	}

	apperrors.GetAPIError(ctx, tasks, http.StatusOK, nil)
}

// @Summary Update Task.
// @Tags Task
// @Produce json
// @Param id path string false "Task"
// @Param request body model.Task true "Task"
// @Success 200 {array}	 model.Task
// @Failure 400 {object} model.APIError
// @Failure 404 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /tasks/{id} [put]
func Put(ctx *gin.Context) {
	task := model.Task{}
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	id := ctx.Param("id")
	if task.Id != "" && task.Id != id {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "ID_MISMATCH", fmt.Errorf("Id in body must match URL")))
		return
	}

	t, ok := model.LibTasks[id]
	if !ok {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "TASK_NOT_FOUND", fmt.Errorf("Task not found")))
		return
	}

	t.Title = task.Title
	t.Done = task.Done

	apperrors.GetAPIError(ctx, t, http.StatusOK, nil)
}

// @Summary Delete Task.
// @Tags Task
// @Produce json
// @Param id path string false "Task"
// @Success 200 {object} model.APIError
// @Failure 404 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /tasks/{id} [delete]
func Del(ctx *gin.Context) {
	id := ctx.Param("id")

	_, ok := model.LibTasks[id]
	if !ok {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotFound, "TASK_NOT_FOUND", fmt.Errorf("Task not found")))
		return
	}

	delete(model.LibTasks, id)
	apperrors.GetAPIError(ctx, nil, http.StatusNoContent, nil)
}

// @Summary	Download images
// @Tags		Task
// @Produce	json
// @Param		request	body		model.Urls	true	"Urls"
// @Success	200		{object}	model.DownloadImages
// @Failure	400		{object}	model.APIError
// @Failure	408     {object}	model.APIError
// @Failure	500		{object}	model.APIError
// @Router		/download/images [post]
func DownloadUrls(ctx *gin.Context) {
	urls := model.Urls{}
	err := ctx.ShouldBindJSON(&urls)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
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

	apperrors.GetAPIError(ctx, gin.H{"results": results}, http.StatusOK, nil)
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
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
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
		apperrors.GetAPIError(ctx, model.QueryTasksResponse{TotalResults: totalRows, TotalPages: 0, Tasks: []model.Task{}}, http.StatusOK, nil)
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

	response := model.QueryTasksResponse{
		TotalResults: totalRows,
		TotalPages:   totalPages,
		Tasks:        pagedTasks,
	}
	apperrors.GetAPIError(ctx, response, http.StatusOK, nil)
}

// @Summary Register a new user.
// @Description Creates a new user account with username and password
// @Tags auth
// @Produce json
// @Param user body model.User true "User registration data"
// @Success 200 {object} model.APIError
// @Failure 400 {object} model.APIError
// @Failure 406 {object} model.APIError
// @Failure 409 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /auth/register [post]
func Register(ctx *gin.Context) {
	user := model.User{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	if len(user.Username) < 0 || len(user.Password) < 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusNotAcceptable, "INVALID_CREDENTIALS", fmt.Errorf("Invalid username/password")))
		return
	}

	if _, ok := model.Users[user.Username]; ok {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusConflict, "USER_ALREADY_EXISTS", fmt.Errorf("User already exists")))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "HASH_ERROR", err))
		return
	}

	if user.Username == "admin" {
		model.Users["admin"] = model.Login{
			HashedPassword: string(hashedPassword),
			Role:           "admin",
		}
	} else {
		model.Users[user.Username] = model.Login{
			HashedPassword: string(hashedPassword),
			Role:           "user",
		}
	}

	// TODO Create user in Database
	/*user := model.User{}
	result := database.DB.Create(&user)
	if result.Error != nil {
	   return result.Error
	}*/

	apperrors.GetAPIError(ctx, gin.H{"message": "User registered successfully!"}, http.StatusOK, nil)
}

// @Summary Login user
// @Description Authenticates user and returns session cookies
// @Tags auth
// @Produce json
// @Param user body model.User true "User registration data"
// @Success 200 {object} model.APIError
// @Failure 400 {object} model.APIError
// @Failure 401 {object} model.APIError
// @Failure 500 {object} model.APIError
// @Router /auth/login [post]
func Login(ctx *gin.Context) {
	userPayload := model.User{}
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	user, ok := model.Users[userPayload.Username]
	if !ok {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", fmt.Errorf("Invalid username or password")))
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userPayload.Password))
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", fmt.Errorf("Invalid username or password")))
		return
	}

	// Generate a jwt token
	expHours, err := strconv.Atoi(os.Getenv("EXPIRATION_HOURS"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userPayload.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Duration(expHours) * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "TOKEN_GENERATION_ERROR", fmt.Errorf("Failed to generate token: %v", err)))
		return
	}

	// Send it back
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	apperrors.GetAPIError(ctx, gin.H{}, http.StatusOK, nil)
}

// @Summary	Logout user.
// @Tags		auth
// @Produce	json
// @Success	200		{array}		model.APIError
// @Failure	500		{object}	model.APIError
// @Router		/auth/logout [get]
func Logout(ctx *gin.Context) {
	// Delete the JWT cookie
	ctx.SetCookie("Authorization", "", -1, "", "", false, true)

	apperrors.GetAPIError(ctx, gin.H{"message": "Logout successfully!"}, http.StatusOK, nil)
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
