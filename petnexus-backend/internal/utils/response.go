// Package utils contains small helpers shared across HTTP handlers.
package utils

import "github.com/gin-gonic/gin"

// Success sends the standard PetNexus success response.
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// Error sends the standard PetNexus error response.
func Error(c *gin.Context, statusCode int, message, code string, details interface{}) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"error": gin.H{
			"code":    code,
			"details": details,
		},
	})
}
