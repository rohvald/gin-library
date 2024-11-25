package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func SendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": data})
}

func SendBadRequestError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": message})
}

func SendNotFoundError(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"success": false, "error": message})
}

func SendInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "internal server error"})
}
