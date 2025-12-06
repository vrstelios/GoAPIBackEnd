package server

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: Add route tests
func TestRoutes(t *testing.T) {
	type testUserCase struct {
		Username string
		Password string
	}

	type testTaskCase struct {
		Id    string
		Title string
		Done  bool
	}
	var resp map[string]string
	var token string

	t.Run("valid login for a simple user", func(t *testing.T) {
		// Define valid login for a simple user, err is nil by default
		var user testUserCase
		user.Username = "vrstelios"
		user.Password = "12345678"

		r := gin.Default()
		r.POST("/api/auth/register", Register)
		r.POST("/api/auth/login", Login)

		jsonUser, _ := json.Marshal(user)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		//fmt.Println("Login response:", w.Body.String())

		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		tokenString := resp["token"]
		if tokenString == "" {
			t.Fatalf("Login did not return a token. Response: %s", w.Body.String())
		}
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}
		claims := token.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		if role != "user" {
			t.Fatalf("expected role user, got %s", role)
		}
	})

	t.Run("valid login for a admin user", func(t *testing.T) {
		// Define valid login for a simple admin, err is nil by default
		var userTest testUserCase
		userTest.Username = "admin"
		userTest.Password = "12345678"

		r := gin.Default()
		r.POST("/api/auth/register", Register)
		r.POST("/api/auth/login", Login)

		jsonUser, _ := json.Marshal(userTest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		token = resp["token"]
		if token == "" {
			t.Fatalf("Token missing: %s", w.Body.String())
		}
		token, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}
		claims := token.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		if role != "admin" {
			t.Fatalf("expected role admin, got %s", role)
		}
	})

	t.Run("Valid server-side filtering and paging", func(t *testing.T) {
		// Define Valid server-side filtering and paging, err is nil by default
		if models.LibTasks == nil {
			models.LibTasks = make(map[string]*models.Task)
		}

		var userTest testUserCase
		userTest.Username = "admin"
		userTest.Password = "12345678"

		var taskTest testTaskCase
		taskTest.Id = ""
		taskTest.Title = "ba"
		taskTest.Done = true

		r := gin.Default()
		//r.POST("/api/auth/register", Register)
		//r.POST("/api/auth/login", Login)
		r.POST("/api/tasks", AuthJWTMiddleware(), RoleMiddleware("admin"), Post)
		r.POST("/api/tasks/query", QueryTasksV2)

		/*jsonUser, _ := json.Marshal(userTest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)*/

		jsonTask, _ := json.Marshal(taskTest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonTask))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var response testTaskCase
		json.Unmarshal(w.Body.Bytes(), &response)

		query := models.QueryTasksRequest{
			Paging: models.QueryPagingRequest{PageSize: 2, Page: 1},
			Search: response.Title,
		}
		jsonQuery, _ := json.Marshal(query)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/tasks/query", bytes.NewBuffer(jsonQuery))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})
}
