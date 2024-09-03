package handlers

import (
	"AppWeb/Utilities"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {

	// Calcula los datos para el catálogo (esto es solo un ejemplo, debes obtener estos datos de tu lógica)
	datosDashboard, _ := Utilities.CheckMetrics()

	session := sessions.Default(c)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"email":          session.Get("email").(string),
		"nombre":         session.Get("nombre").(string),
		"apellido":       session.Get("apellido").(string),
		"rol":            session.Get("rol").(uint8),
		"machines":       nil,
		"machinesChange": nil,
		"datosDashboard": datosDashboard,
	})
}
