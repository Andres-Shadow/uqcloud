package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
)

/*
Clase encargada de contener los handlers que responden a las acciones sobre los discos
*/

// Función que se encarga de registrar un disco en la base de datos
// realiza un llamado a su respectiva función en la base de datos
// para registrar en la base de datos un nuevo disco para maquina virtual 
func AddDiskHandler(w http.ResponseWriter, r *http.Request) {
	var disco models.Disco

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&disco); err != nil {
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	query := "insert into disco (nombre, ruta_ubicacion, sistema_operativo, distribucion_sistema_operativo, arquitectura, host_id) values (?, ?, ?, ?, ?, ?);"

	_, err := database.DB.Exec(query, disco.Nombre, disco.Ruta_ubicacion, disco.Sistema_operativo, disco.Distribucion_sistema_operativo, disco.Arquitectura, disco.Host_id)
	if err != nil {
		log.Println("Error al registrar el disco.")
		return

	} else if err != nil {
		panic(err.Error())
	}
	fmt.Println("Registro del disco exitoso")
	response := map[string]bool{"registroCorrecto": true}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}