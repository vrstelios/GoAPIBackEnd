package api

import (
	"GoAPIBackEnd/internal/type/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *APIClient) QueryWorkouts(userId string) ([]models.Workouts, error) {
	baseURL := strings.TrimSuffix(c.BaseURL, "/")
	url := fmt.Sprintf("%s/api/workouts", baseURL)

	if userId != "" {
		url = fmt.Sprintf("%s?userId=%s", url, userId)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return []models.Workouts{}, nil
		}

		var apiErr models.APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var workouts []models.Workouts
	if err := json.Unmarshal(body, &workouts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return workouts, nil
}

func (c *APIClient) CreateWorkoutLogs(logs []models.WorkoutLog) ([]models.WorkoutLog, error) {
	baseURL := strings.TrimSuffix(c.BaseURL, "/")
	url := fmt.Sprintf("%s/api/workoutLogs", baseURL)

	if len(logs) == 0 {
		fmt.Println("  WARNING: No logs to create")
		return []models.WorkoutLog{}, nil
	}

	jsonData, err := json.Marshal(logs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr models.APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var createdLogs []models.WorkoutLog
	if err := json.Unmarshal(body, &createdLogs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return createdLogs, nil
}
