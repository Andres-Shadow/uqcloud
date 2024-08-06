package handlers

import (
	"encoding/json"
	"net/http"
	utilities "servidor_procesamiento/Procesador/Utilities"
)

/*
Clase encargada de contener las funciones que responen hacía los endpoints de las máquinas virtuales
para invitados
*/

// Funcion que responde al endpoint encargado de crear maquinas virtuales para invitados
func CreateGuestVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica que la solicitud sea del método POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	var datos map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&datos); err != nil {
		http.Error(w, "Error al decodificar el JSON ", http.StatusBadRequest)
		return
	}

	clientIP := datos["ip"].(string)
	distribucion_SO := datos["distribucion"].(string)
	email := utilities.CreateTempAccount(clientIP, distribucion_SO)

	// Envía una respuesta al cliente.
	response := map[string]string{"mensaje": email}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}