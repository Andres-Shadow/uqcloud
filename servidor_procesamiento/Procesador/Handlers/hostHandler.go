package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	dto "servidor_procesamiento/Procesador/Models/Dto"
	external "servidor_procesamiento/Procesador/Models/External"
	"servidor_procesamiento/Procesador/Utilities/mapper"
	"strconv"

	utilities "servidor_procesamiento/Procesador/Utilities"

	"github.com/gorilla/mux"
)

/*
Clase encargada de contener los handlers que responden a los endpoints que
atienden las consultas sobre estos
*/

// Funcion que responde al endpoint encargado de consultar las maquinas virtuales
func ConsultHostsHandler(w http.ResponseWriter, r *http.Request) {
	var hosts []map[string]interface{}
	var err error
	hosts, err = database.ConsultHosts() // almacena los hosts que se encuentran en la base de datos

	if err != nil && err.Error() != "no Hosts encontrados" {
		log.Println("Error al consultar los host, no se encontraron máquinas host registradas en la base de datos")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := buildHostList(hosts, "success", "Consulta de hosts exitosa")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func buildHostList(hosts []map[string]interface{}, status, message string) external.HostListResponse {
	var response external.HostListResponse
	var count int

	if hosts != nil {
		count = len(hosts)
	} else {
		count = 0
	}

	response.Status = status
	response.Data = hosts
	response.Message = message
	response.Count = count

	return response

}

// Funcion que responde al endpoint encargado de consultar los host registrados en la base de datos
func ConsultHostHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener la variable "name" de la ruta
	vars := mux.Vars(r)
	email := vars["email"]

	if email == "" {
		log.Println("Error al obtener el email del usuario")
		http.Error(w, "Error al obtener el email del usuario", http.StatusBadRequest)
		return
	}

	_, error := database.GetUser(email)
	if error != nil {
		http.Error(w, "Usuario no encontrado en la base de datos", http.StatusBadRequest)
		return
	}

	hosts, err := database.ConsultHosts()

	if err != nil && err.Error() != "no Hosts encontrados" {
		fmt.Println(err)
		log.Println("Error al consultar los Host")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := buildHostList(hosts, "success", "Consulta de Host exitosa")
	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Funcion que responde al endpoint encargado de agregar un host a la base de datos
func AddHostHandler(w http.ResponseWriter, r *http.Request) {
	var host dto.HostDTO

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&host); err != nil {
		log.Println("Error al decodificar JSON de especificaciones")
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	converted := mapper.ToHostFromDTO(host)

	err := database.AddHost(converted)

	if err != nil {
		fmt.Println(err)
		log.Println("Error al registrar el host en la base de datos")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// recarga la lsita de host contemplados en roundrobin
	config.RoundRobinManager.UpdateHosts(database.GetHosts())

	// Recargar configuración de Prometheus
	config.ReloadPrometheusConfig()

	fmt.Println("Registro del host exitoso")
	response := map[string]bool{"registroCorrecto": true}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func FastRegisterHostsHandler(w http.ResponseWriter, r *http.Request) {
	// Mapa genérico para decodificar el JSON
	var data map[string][]string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Println("Error al decodificar JSON de IPs")
		http.Error(w, "Error al decodificar JSON de IPs", http.StatusBadRequest)
		return
	}

	// Obtenemos el slice de IPs del campo "ip"
	ips, ok := data["ips"]
	if !ok {
		log.Println("El JSON no contiene el campo 'ip'")
		http.Error(w, "El JSON no contiene el campo 'ip'", http.StatusBadRequest)
		return
	}

	// Llamar a utilidades con las IPs
	utilities.FastRegisterHosts(ips)

	// Recargar configuración de Prometheus
	config.ReloadPrometheusConfig()

	// recarga la lsita de host contemplados en roundrobin
	config.RoundRobinManager.UpdateHosts(database.GetHosts())

	confirmation := map[string]bool{"registro rapido correcto": true}
	response := utilities.BuildGenericResponse(confirmation, "success", "Registro rapido de los hosts exitoso")
	fmt.Println("Registro de hosts exitoso")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteHostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostId := vars["id"]

	id, err := strconv.Atoi(hostId)
	if err != nil {
		log.Println("Error al convertir el ID del host a entero")
		http.Error(w, "ID del host inválido", http.StatusBadRequest)
		return
	}

	err = database.DeleteHostById(id)

	if err != nil {
		log.Println("Error al eliminar el host, no se encontró un host con el nombre asociado")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// recarga la lsita de host contemplados en roundrobin
	config.RoundRobinManager.UpdateHosts(database.GetHosts())

	// Recargar configuración de Prometheus
	config.ReloadPrometheusConfig()

	confirmation := map[string]bool{"eliminacion de host": true}
	response := utilities.BuildGenericResponse(confirmation, "success", "Eliminación del host exitosa")
	fmt.Println("Eliminación del host exitosa")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
