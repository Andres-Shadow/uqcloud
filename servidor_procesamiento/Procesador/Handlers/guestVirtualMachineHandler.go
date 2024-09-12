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
	// retorna el correo temporal
	email := utilities.CreateTempAccount()

	// Envía una respuesta al cliente.
	response := map[string]string{"mensaje": email}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
