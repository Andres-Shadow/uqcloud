package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cargar la pagina principal
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
