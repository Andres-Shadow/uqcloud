package handlers

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {

	// Acceder a la sesión
	session := sessions.Default(c)
	rol := session.Get("rol")
	if rol != "Administrador" {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Calcula los datos para el catálogo (esto es solo un ejemplo, debes obtener estos datos de tu lógica)
	datosDashboard, _ := CheckMetrics()

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"email":          "email",
		"machines":       nil,
		"machinesChange": nil,
		"datosDashboard": datosDashboard,
	})
}

// Funcion que consulta las metricas
func CheckMetrics() (Models.DashboardData, error) {
	var metricas Models.DashboardData
	serverURL := fmt.Sprintf("http://%s:8081/json/consultMetrics", Config.ServidorProcesamientoRoute)

	resp, err := http.Get(serverURL)
	if err != nil {
		return metricas, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return metricas, err
	}

	err = json.NewDecoder(resp.Body).Decode(&metricas)
	if err != nil {
		return metricas, err
	}

	return metricas, nil
}
