package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

// Funcion que consulta las metricas
func CheckMetrics() (Models.DashboardData, error) {
	var metricas Models.DashboardData
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.METRICS_URL)

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
