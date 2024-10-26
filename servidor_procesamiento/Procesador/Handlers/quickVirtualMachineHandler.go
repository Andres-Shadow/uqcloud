package handlers

import (
	"encoding/json"
	"net/http"
	virtualmachineutilities "servidor_procesamiento/Procesador/Utilities/VirtualMachineUtilities"
)

func CreateQuickVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {

	var datos map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&datos); err != nil {
		http.Error(w, "Error al decodificar el JSON ", http.StatusBadRequest)
		return
	}

	//ip del localhost ->
	clientIP := datos["ip"].(string)
	email := virtualmachineutilities.CreateQuickVirtualMachine(clientIP)

	response := map[string]string{"mensaje": email}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
