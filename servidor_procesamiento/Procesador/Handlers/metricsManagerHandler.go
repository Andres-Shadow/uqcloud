package handlers

import (
	"encoding/json"
	"net/http"
	database "servidor_procesamiento/Procesador/Database"

)

/*
Clase encargada de contener los manejadores que responden a peticiones sobre
las metricas del sistema
*/

// Funcion que responde al endpoint encargado de consultar las metricas del sistema
func ConsultMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricas, err := database.GetMetrics()

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricas)
}