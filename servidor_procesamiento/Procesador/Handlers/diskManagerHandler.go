package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	database "servidor_procesamiento/Procesador/Database"
	dto "servidor_procesamiento/Procesador/Models/Dto"
	utilities "servidor_procesamiento/Procesador/Utilities"
	"servidor_procesamiento/Procesador/Utilities/mapper"
	"strings"

	"github.com/gorilla/mux"
)

/*
Clase encargada de contener los handlers que responden a las acciones sobre los discos
*/

// Función que se encarga de registrar un disco en la base de datos
// realiza un llamado a su respectiva función en la base de datos
// para registrar en la base de datos un nuevo disco para maquina virtual
func AddDiskHandler(w http.ResponseWriter, r *http.Request) {
	var disco dto.DiscoDTO
	var rutaDisco string

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&disco); err != nil {
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}
	rutaDisco = strings.ReplaceAll(disco.Ruta_ubicacion, "/", "\\")

	disco.Ruta_ubicacion = rutaDisco

	convertedDisk := mapper.ToDiscoFromDTO(disco)

	err := database.CreateDisck(convertedDisk)

	if err != nil {
		fmt.Println("Error al registrar el disco en la base de datos: " + err.Error())
		http.Error(w, "Error al registrar el disco en la base de datos: "+err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("Registro del disco exitoso")
	confirmation := map[string]bool{"registro_correcto": true}

	response := utilities.BuildGenericResponse(confirmation, "Success", "Registro del disco exitoso")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetDisksHandler(w http.ResponseWriter, r *http.Request) {
	disks, err := database.ListUniquesDisks()

	if err != nil {
		http.Error(w, "Error al obtener los discos de la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mapped := mapper.ToDTOFromDiskDistroList(disks)
	response := utilities.BuildGenericResponse(mapped, "Success", "Discos obtenidos correctamente")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetHostsWithDiskHandler(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	diskSoDistro := pathVars["name"]

	hosts, err := database.ListHostWhereDiskExists(diskSoDistro)

	if err != nil {
		http.Error(w, "Error al obtener los hosts de la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mapped := mapper.ToDTOFromHostWithDiskList(hosts)
	response := utilities.BuildGenericResponse(mapped, "Success", "Hosts obtenidos correctamente")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteDiskHandler(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	diskSoDistro := pathVars["name"]

	queryVars := r.URL.Query()
	hostId := queryVars.Get("host_id")

	err := database.DeleteDiskFromHost(hostId, diskSoDistro)

	if err != nil {
		log.Println("Error al eliminar el disco de la base de datos: " + err.Error())
		http.Error(w, "Error al eliminar el disco de la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	confirmation := map[string]bool{"eliminacion_correcta": true}
	response := utilities.BuildGenericResponse(confirmation, "Success", "Disco eliminado correctamente")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
