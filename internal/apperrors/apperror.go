package apperrors

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIError struct {
	Status       int    `json:"status"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Error        error  `json:"-"`
}

var (
	ErrUnauthorized = &APIError{
		Status:       http.StatusUnauthorized,
		ErrorCode:    "UNAUTHORIZED",
		ErrorMessage: "authorization token required",
	}

	ErrInvalidToken = &APIError{
		Status:       http.StatusUnauthorized,
		ErrorCode:    "INVALID_TOKEN",
		ErrorMessage: "invalid or expired token",
	}

	ErrForbidden = &APIError{
		Status:       http.StatusForbidden,
		ErrorCode:    "FORBIDDEN",
		ErrorMessage: "you are not allowed to access this resource",
	}
)

func (APIError) NewError(err error) *APIError {
	return APIError{}.New(0, "", err)
}

func (APIError) New(httpStatus int, errorCode string, err error) *APIError {
	var apiErr = APIError{
		Status:    httpStatus,
		ErrorCode: errorCode,
		Error:     err,
	}
	logger.Warn("API Error: %+v", apiErr)
	return &apiErr
}

func GetAPIError(ctx *gin.Context, data interface{}, statusCode int, err *APIError) {
	if err != nil {
		httpStatus := err.Status
		if httpStatus == 0 {
			httpStatus = http.StatusInternalServerError
		}

		ctx.JSON(httpStatus, gin.H{
			"status":       err.Status,
			"errorCode":    err.ErrorCode,
			"errorMessage": err.ErrorMessage,
		})
		return
	}

	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	if data == nil {
		ctx.Status(statusCode)
		return
	}

	ctx.JSON(statusCode, data)
}
