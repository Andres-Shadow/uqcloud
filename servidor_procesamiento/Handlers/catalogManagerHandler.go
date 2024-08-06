package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	utilities "servidor_procesamiento/Procesador/Utilities"
)

/*
Clase encargada de contener los handlers que responden a los eventos de la gestión de los catálogos
*/

// Función que responde a la solicitud de consulta de catálogos
// realiza un llamado a su respectiva función para retornar el catalogo de host
func ConsultCatalogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Se requiere una solicitud Get", http.StatusMethodNotAllowed)
		return
	}

	catalogo, err := utilities.ConsultCatalog()
	if err != nil {
		log.Printf("Error al consultar el catálogo: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(catalogo); err != nil {
		log.Printf("Error al codificar la respuesta JSON: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

}