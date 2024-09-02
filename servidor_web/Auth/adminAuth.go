package authentication

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acceso no autorizado"})
		return
	}
	c.Next()
}
