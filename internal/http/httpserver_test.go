package http

import (
	"GoAPIBackEnd/internal/api"
	"GoAPIBackEnd/internal/database"
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: Add route tests
func TestRoutes(t *testing.T) {
	t.Run("valid login for a simple user", func(t *testing.T) {
		// Define valid login for a simple user, err is nil by default
		conn, err := pgx.Connect(context.Background(), "host=localhost user=postgres password=postgres port=5432 sslmode=disable")
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		defer conn.Close(context.Background())
		database.Conn = conn

		signupUser := map[string]string{
			"name":     "testUser",
			"password": "12345678",
			"email":    "testUser@gmail.com",
			"role":     "athlete",
		}

		r := gin.Default()
		r.POST("/api/auth/signup", api.Signup)
		r.POST("/api/auth/login", api.Login)

		signupJson, _ := json.Marshal(signupUser)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(signupJson))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(signupJson))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		//fmt.Println("Login response:", w.Body.String())

		// Login
		loginData := map[string]string{
			"name":     "testUser",
			"password": "12345678",
		}

		loginJson, _ := json.Marshal(loginData)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginJson))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		t.Logf("Login response: %d - %s", w.Code, w.Body.String())

		// Check for 200 status
		assert.Equal(t, 200, w.Code, "Login should return 200")
	})

	t.Run("valid login for a admin user", func(t *testing.T) {
		// Define valid login for a simple admin, err is nil by default
		conn, err := pgx.Connect(context.Background(), "host=localhost user=postgres password=postgres port=5432 sslmode=disable")
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		defer conn.Close(context.Background())
		database.Conn = conn

		signupUser := map[string]string{
			"name":     "testAdmin",
			"password": "12345679",
			"email":    "testAdmin@gmail.com",
			"role":     "admin",
		}

		r := gin.Default()
		r.POST("/api/auth/signup", api.Signup)
		r.POST("/api/auth/login", api.Login)

		jsonAdminUser, _ := json.Marshal(signupUser)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(jsonAdminUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		t.Logf("Admin signup response: %d - %s", w.Code, w.Body.String())
		assert.Equal(t, 200, w.Code)

		// Login
		loginUser := map[string]string{
			"name":     "testAdmin",
			"password": "12345679",
		}

		loginJson, _ := json.Marshal(loginUser)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginJson))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		t.Logf("Admin login response: %d - %s", w.Code, w.Body.String())
		assert.Equal(t, 200, w.Code)
	})
}
