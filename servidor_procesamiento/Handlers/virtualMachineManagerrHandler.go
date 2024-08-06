package handlers

/*
Clase encargada de contener las funciones asociadas los handlers para cada Endpoint
al cual la API responda
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
	"strconv"
)

  
// Funcion que responde al endpoint encargado de crear una maquina virtual
func CreateVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica que la solicitud sea del método POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica el JSON recibido en la solicitud en una estructura Specifications.
	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Error al decodificar JSON de la solicitud", http.StatusBadRequest)
		return
	}

	// Encola las especificaciones.
	config.GetMu().Lock()
	config.GetMaquina_virtualQueue().Queue.PushBack(payload)
	config.GetMu().Unlock()

	fmt.Println("Cantidad de Solicitudes de Especificaciones en la Cola: " + strconv.Itoa(config.GetMaquina_virtualQueue().Queue.Len()))

	// Envía una respuesta al cliente.
	response := map[string]string{"mensaje": "Mensaje JSON de crear MV recibido correctamente"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Funcion que responde al endpoint encargado de consultar el estado de las maquinas virtuales en tiempo real
func ConsultVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	var persona models.Persona
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&persona); err != nil { //Solo llega el email
		http.Error(w, "Error al decodificar JSON de inicio de sesión", http.StatusBadRequest)
		return
	}

	persona, error := database.GetUser(persona.Email)
	if error != nil {
		return
	}

	machines, err := database.ConsultMachines(persona)
	if err != nil && err.Error() != "no Machines Found" {
		fmt.Println(err)
		log.Println("Error al consultar las màquinas del usuario")
		return
	}

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(machines)

}

// Funcion que responde al endpoint encargado de modificar una maquina virtual (en caliente o apagada)
func ModifyVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica que la solicitud sea del método POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica el JSON recibido en la solicitud en un mapa genérico.
	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Error al decodificar JSON de la solicitud", http.StatusBadRequest)
		return
	}

	// Verifica que el campo "tipo_solicitud" esté presente y sea "modify".
	tipoSolicitud, isPresent := payload["tipo_solicitud"].(string)
	if !isPresent || tipoSolicitud != "modify" {
		http.Error(w, "El campo 'tipo_solicitud' debe ser 'modify'", http.StatusBadRequest)
		return
	}

	// Extrae el objeto "specifications" del JSON.
	specificationsData, isPresent := payload["specifications"].(map[string]interface{})
	if !isPresent || specificationsData == nil {
		http.Error(w, "El campo 'specifications' es inválido", http.StatusBadRequest)
		return
	}

	// Encola las peticiones.
	config.GetMu().Lock()
	config.GetManagementQueue().Queue.PushBack(payload)
	config.GetMu().Unlock()

	// Envía una respuesta al cliente.
	response := map[string]string{"mensaje": "Mensaje JSON de especificaciones para modificar MV recibido correctamente"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Funcion que responde al endpoint encargado de eliminar una maquina virtual
func DeleteVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica que la solicitud sea del método POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	var datos map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&datos); err != nil {
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	// Verificar si el nombre de la máquina virtual, la IP del host y el tipo de solicitud están presentes y no son nulos
	nombre, nombrePresente := datos["nombreVM"].(string)
	tipoSolicitud, tipoPresente := datos["tipo_solicitud"].(string)

	if !tipoPresente || tipoSolicitud != "delete" {
		http.Error(w, "El campo 'tipo_solicitud' debe ser 'delete'", http.StatusBadRequest)
		return
	}

	if !nombrePresente || !tipoPresente || nombre == "" || tipoSolicitud == "" {
		http.Error(w, "El tipo de solicitud y el nombre de la máquina virtual son obligatorios", http.StatusBadRequest)
		return
	}

	// Encola las peticiones.
	config.GetMu().Lock()
	config.GetManagementQueue().Queue.PushBack(datos)
	config.GetMu().Unlock()

	// Envía una respuesta al cliente.
	response := map[string]string{"mensaje": "Mensaje JSON para eliminar MV recibido correctamente"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Funcion que responde al endpoint encargado de encender una maquina virtual
func StartVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica que la solicitud sea del método POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	var datos map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&datos); err != nil {
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	// Verificar si el nombre de la máquina virtual, la IP del host y el tipo de solicitud están presentes y no son nulos
	nombreVM, nombrePresente := datos["nombreVM"].(string)

	tipoSolicitud, tipoPresente := datos["tipo_solicitud"].(string)

	if !tipoPresente || tipoSolicitud != "start" {
		http.Error(w, "El campo 'tipo_solicitud' debe ser 'start'", http.StatusBadRequest)
		return
	}

	if !nombrePresente || !tipoPresente || nombreVM == "" || tipoSolicitud == "" {
		http.Error(w, "El tipo de solicitud y nombre de la máquina virtual son obligatorios", http.StatusBadRequest)
		return
	}

	// Encola las peticiones.
	config.GetMu().Lock()
	config.GetManagementQueue().Queue.PushBack(datos)
	config.GetMu().Unlock()

	query := "select estado from maquina_virtual where nombre = ?"
	var estado string

	//Registra el usuario en la base de datos
	database.DB.QueryRow(query, nombreVM).Scan(&estado)

	mensaje := "Apagando "
	if estado == "Apagado" {
		mensaje = "Encendiendo "
	}

	// Envía una respuesta al cliente.
	response := map[string]string{"mensaje": mensaje}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Funcion que responde al endpoint encargado de apagar una maquina virtual
func StopVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica que la solicitud sea del método POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	var datos map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&datos); err != nil {
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	// Verificar si el nombre de la máquina virtual, la IP del host y el tipo de solicitud están presentes y no son nulos
	nombreVM, nombrePresente := datos["nombreVM"].(string)

	tipoSolicitud, tipoPresente := datos["tipo_solicitud"].(string)

	if !tipoPresente || tipoSolicitud != "stop" {
		http.Error(w, "El campo 'tipo_solicitud' debe ser 'stop'", http.StatusBadRequest)
		return
	}

	if !nombrePresente || !tipoPresente || nombreVM == "" || tipoSolicitud == "" {
		http.Error(w, "El tipo de solicitud y nombre de la máquina virtual son obligatorios", http.StatusBadRequest)
		return
	}

	// Encola las peticiones.
	config.GetMu().Lock()
	config.GetManagementQueue().Queue.PushBack(datos)
	config.GetMu().Unlock()

	// Envía una respuesta al cliente.
	response := map[string]string{"mensaje": "Mensaje JSON para apagar MV recibido correctamente"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}