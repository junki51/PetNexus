// Package handlers translates HTTP requests into application responses.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

// Health reports whether the backend process is running.
func Health(c *gin.Context) {
	utils.Success(c, http.StatusOK, "PetNexus backend is running", gin.H{
		"status":  "ok",
		"service": "petnexus-backend",
	})
}
