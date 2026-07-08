package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

func parseAppointmentID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid appointment ID", "INVALID_APPOINTMENT_ID", "Appointment ID must be a valid UUID")
		return uuid.Nil, false
	}
	return id, true
}

func respondInvalidAppointmentJSON(c *gin.Context) {
	utils.Error(c, http.StatusBadRequest, "Invalid request", "INVALID_REQUEST", "Request body must be valid JSON")
}

func respondAppointmentError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus >= http.StatusInternalServerError {
			log.Printf("appointment request failed: %v", err)
		}
		utils.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	log.Printf("unexpected appointment error: %v", err)
	utils.Error(c, http.StatusInternalServerError, "Something went wrong", "INTERNAL_SERVER_ERROR", "An internal server error occurred")
}
