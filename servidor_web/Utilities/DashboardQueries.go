package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"fmt"
	"log"
	"net/http"

	"github.com/goccy/go-json"
)

// Funcion que consulta las metricas
func CheckMetrics() (Models.DashboardData, error) {
	var metricas Models.DashboardData
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.METRICS_URL)

	log.Println(serverURL)

	resp, err := http.Get(serverURL)
	if err != nil {
		log.Println("Error getting metrics")
		return metricas, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error al enviar la peticion", resp.Status)
		return metricas, err
	}

	err = json.NewDecoder(resp.Body).Decode(&metricas)
	if err != nil {
		log.Println("error al decodificar el cuerpo", err.Error())
		return metricas, err
	}

	return metricas, nil
}
