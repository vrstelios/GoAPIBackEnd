package api

import (
	"GoAPIBackEnd/internal/apperrors"
	"GoAPIBackEnd/internal/config"
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/type/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// @Summary Create Workout Log.
// @Tags WorkoutLog
// @Produce json
// @Param request body models.WorkoutLog true "Workout Log"
// @Success 200 {object} models.APIError
// @Failure 208 {object} models.APIError
// @Failure 400 {object} models.APIError
// @Failure 500 {object} models.APIError
// @Router /workoutLogs [post]
func PostWorkoutLogs(ctx *gin.Context) {
	var dbConn = config.Conn
	var wls []models.WorkoutLog
	err := ctx.ShouldBindJSON(&wls)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	if len(wls) == 0 {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "EMPTY_REQUEST", fmt.Errorf("no logs provided")))
		return
	}

	for i := range wls {
		if len(wls[i].Id) == 0 {
			uid, err := uuid.NewUUID()
			if err != nil {
				apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "UUID_ERROR", err))
				return
			}
			wls[i].Id = uid.String()
		}
	}

	err = database.PostWorkoutLog(dbConn, wls)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusInternalServerError, "DB_ERROR", err))
		return
	}

	apperrors.GetAPIError(ctx, wls, http.StatusOK, nil)
}
