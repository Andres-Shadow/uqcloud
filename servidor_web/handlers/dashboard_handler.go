package handlers

import (
	"AppWeb/Utilities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {

	// TODO: DESCOMENTAR PARA QUE ENTREN SOLO LOS ADMIN
	// Acceder a la sesión
	// session := sessions.Default(c)
	// rol := session.Get("rol")
	// if rol != "Administrador" {
	// 	// Si el usuario no está autenticado, redirige a la página de inicio de sesión
	// 	c.Redirect(http.StatusFound, "/login")
	// 	return
	// }

	// Calcula los datos para el catálogo (esto es solo un ejemplo, debes obtener estos datos de tu lógica)
	datosDashboard, _ := Utilities.CheckMetrics()

	// TODO: CAMBIAR AL DASHBOARD.HTML
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"email":          "email",
		"machines":       nil,
		"machinesChange": nil,
		"datosDashboard": datosDashboard,
	})
}
