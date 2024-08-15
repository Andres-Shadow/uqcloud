package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
	utilities "servidor_procesamiento/Procesador/Utilities"
	"strconv"

	"github.com/gorilla/mux"
)

/*
Clase encargada de contener los handlers que responden a los endpoints que
atienden las consultas sobre estos
*/

// Funcion que responde al endpoint encargado de consultar las maquinas virtuales
func ConsultHostsHandler(w http.ResponseWriter, r *http.Request) {
	hosts, err := database.ConsultHosts()
	if err != nil && err.Error() != "no Hosts encontrados" {
		fmt.Println(err)
		log.Println("Error al consultar los host, no se encontraron máquinas host registradas en la base de datos")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hosts)

}

// Funcion que responde al endpoint encargado de verificar los diferentes host registrados en la base de datos
func CheckHostHandler(w http.ResponseWriter, r *http.Request) {

	// Decodifica el JSON recibido en la solicitud en una estructura Specifications.
	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Error al decodificar JSON de la solicitud", http.StatusBadRequest)
		return
	}

	//Se capturan los datos de la maquina virtual a crear
	mv := payload["specifications"].(map[string]interface{})

	/*En la base de datos los indices de los host empiezan desde el indice 1
	si el valor es cero se utiliza para disparar el algoritmo aleatorio
	*/
	id := int(mv["Host_id"].(float64))
	switch {
	case id == 0:
		//Se encola la maquina virtual a crear
		config.GetMu().Lock()
		config.GetMaquina_virtualQueue().Queue.PushBack(payload)
		config.GetMu().Unlock()

		//Se imprime el estado actual de la cola
		fmt.Println("Cantidad de Solicitudes de Especificaciones en la Cola: " + strconv.Itoa(config.GetMaquina_virtualQueue().Queue.Len()))

		// Envía una respuesta al servidor web
		response := map[string]string{"mensaje": "Mensaje JSON de crear MV recibido correctamente", "centinela": "true"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	case id > 0:

		mihost, _ := database.GetHost(int(mv["Host_id"].(float64)))
		estadossh := utilities.Pacemaker(config.GetPrivateKeyPath(), mihost.Hostname, mihost.Ip)
		if estadossh {
			//Se encola la maquina virtual a crear
			config.GetMu().Lock()
			config.GetMaquina_virtualQueue().Queue.PushBack(payload)
			config.GetMu().Unlock()

			//Se imprime el estado actual de la cola
			fmt.Println("Cantidad de Solicitudes de Especificaciones en la Cola: " + strconv.Itoa(config.GetMaquina_virtualQueue().Queue.Len()))

			// Envía una respuesta al servidor web
			response := map[string]string{"mensaje": "Mensaje JSON de crear MV recibido correctamente", "centinela": "true"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			// Envía una respuesta al servidor web
			response := map[string]string{"mensaje": "Esta maquina tiene problemas :(", "centinela": "false", "hostmalo:": strconv.Itoa(id)}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}

	case id < 0:
		// Envía una respuesta al cliente.
		response := map[string]string{"mensaje": "Error de seleccion de maquina", "centinela": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

// Funcion que responde al endpoint encargado de consultar los host registrados en la base de datos
func ConsultHostHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener la variable "name" de la ruta
    vars := mux.Vars(r)
    email := vars["email"]

	if email == "" {
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

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hosts)

}

// Funcion que responde al endpoint encargado de agregar un host a la base de datos
func AddHostHandler(w http.ResponseWriter, r *http.Request) {
	var host models.Host

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&host); err != nil {
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	err := database.AddHost(host)

	if err != nil {
		fmt.Println(err)
		log.Println("Error al registrar el host en la base de datos")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Println("Registro del host exitoso")
	response := map[string]bool{"registroCorrecto": true}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
