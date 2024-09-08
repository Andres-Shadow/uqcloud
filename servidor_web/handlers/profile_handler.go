package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ProfilePage(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)

	if session.Get("email") == nil {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"email":    session.Get("email").(string),
		"nombre":   session.Get("nombre").(string),
		"apellido": session.Get("apellido").(string),
		"rol":      session.Get("rol").(uint8),
	})
}
