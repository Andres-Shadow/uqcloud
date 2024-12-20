package handlers

/*
Clase encargada de contener las funciones asociadas los handlers para cada Endpoint
al cual la API responda
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	entities "servidor_procesamiento/Procesador/Models/Entities"
	apiutilities "servidor_procesamiento/Procesador/Utilities/ApiUtilities"
	systemutilities "servidor_procesamiento/Procesador/Utilities/SystemUtilities"
	"strconv"

	"github.com/gorilla/mux"
)

// Funcion que responde al endpoint encargado de crear una maquina virtual
func CreateVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {

	// Decodifica el JSON recibido en la solicitud en una estructura Specifications.
	// var virtualMachineDTO dto.CreateVMRequestDTO
	var virtualMachineDTO map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&virtualMachineDTO); err != nil {
		http.Error(w, "Error al decodificar JSON de la solicitud", http.StatusBadRequest)
		log.Println("Error al decodificar JSON de la solicitud")
		return
	}

	// // Verifica si el JSON recibido en la solicitud no es un JSON vacío
	if virtualMachineDTO == nil {
		http.Error(w, "El JSON de la solicitud está vacío", http.StatusBadRequest)
		log.Println("El JSON de la solicitud está vacío")
		return
	}

	// Encola las especificaciones.
	config.GetMu().Lock()
	config.GetMaquina_virtualQueue().Queue.PushBack(virtualMachineDTO)
	config.GetMu().Unlock()

	fmt.Println("Cantidad de Solicitudes de Especificaciones en la Cola: " + strconv.Itoa(config.GetMaquina_virtualQueue().Queue.Len()))

	// Envía una respuesta al cliente.
	confirmation := map[string]string{"mensaje": "Mensaje JSON de crear MV recibido correctamente"}

	response := apiutilities.BuildGenericResponse(confirmation, "Success", "Mensaje JSON de crear MV recibido correctamente")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Funcion que responde al endpoint encargado de consultar el estado de las maquinas virtuales en tiempo real
func ConsultVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {

	// extrae el path param con el correo
	vars := mux.Vars(r)
	email := vars["email"]

	persona, error := database.GetUser(email)
	if error != nil {
		http.Error(w, "Error al consultar el usuario en la base de datos", http.StatusBadRequest)
		return
	}

	machines, err := database.ConsultMachines(persona)

	if err != nil && err.Error() != "no Machines Found" {
		log.Println("Error al consultar las màquinas del usuario")
		http.Error(w, "Error al consultar las màquinas del usuario", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Println("No se encontraron màquinas virtuales para el usuario")
		http.Error(w, "No se encontraron màquinas virtuales para el usuario", http.StatusNoContent)
		return
	}

	response := apiutilities.BuildGenericResponse(machines, "Success", "Máquinas virtuales consultadas correctamente")

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Funcion que responde al endpoint encargado de eliminar una maquina virtual
func DeleteVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {

	// Obtener el nombre de la máquina virtual a partir del path param.
	vars := mux.Vars(r)
	virtualMachineName := vars["name"]

	// verificar que el path param no esté vacío
	if virtualMachineName == "" {
		http.Error(w, "El nombre de la máquina virtual es obligatorio", http.StatusBadRequest)
		return
	}

	payload := make(map[string]interface{})
	payload["tipo_solicitud"] = "delete"
	payload["nombreVM"] = virtualMachineName

	// Encola las peticiones.
	config.GetMu().Lock()
	config.GetManagementQueue().Queue.PushBack(payload)
	config.GetMu().Unlock()

	// Envía una respuesta al cliente.
	confirmation := map[string]string{"mensaje": "Mensaje JSON para eliminar MV recibido correctamente"}

	response := apiutilities.BuildGenericResponse(confirmation, "Success", "Mensaje JSON para eliminar MV recibido correctamente")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Funcion que responde al endpoint encargado de encender una maquina virtual
func StartVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	var datos map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&datos); err != nil {
		log.Println("Error al decodificar JSON de especificaciones")
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	// Verifica si el JSON recibido en la solicitud no es un JSON vacío
	if datos == nil {
		http.Error(w, "El JSON de la solicitud está vacío", http.StatusBadRequest)
		log.Println("El JSON de la solicitud está vacío")
		return
	}

	// Verificar si el nombre de la máquina virtual, la IP del host y el tipo de solicitud están presentes y no son nulos
	nombreVM, nombrePresente := datos["nombreVM"].(string)

	if !nombrePresente || nombreVM == "" {
		log.Println("El tipo de solicitud y nombre de la máquina virtual son obligatorios")
		http.Error(w, "El tipo de solicitud y nombre de la máquina virtual son obligatorios", http.StatusBadRequest)
		return
	}

	datos["tipo_solicitud"] = "start"
	// Encola las peticiones.
	config.GetMu().Lock()
	config.GetManagementQueue().Queue.PushBack(datos)
	config.GetMu().Unlock()

	estado, err := database.GetStateFromVirtualMachineName(nombreVM)

	if err != nil {
		http.Error(w, "Error al obtener el estado de la máquina virtual", http.StatusBadRequest)
		return
	}

	mensaje := "Apagando "
	if estado == "Apagado" {
		mensaje = "Encendiendo "
	}

	// Envía una respuesta al cliente.
	confirmation := map[string]string{"mensaje": mensaje}

	response := apiutilities.BuildGenericResponse(confirmation, "Success", mensaje)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Funcion que responde al endpoint encargado de apagar una maquina virtual
func StopVirtualMachineHandler(w http.ResponseWriter, r *http.Request) {
	var datos map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&datos); err != nil {
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	// Verifica si el JSON recibido en la solicitud no es un JSON vacío
	if datos == nil {
		http.Error(w, "El JSON de la solicitud está vacío", http.StatusBadRequest)
		log.Println("El JSON de la solicitud está vacío")
		return
	}

	// Verificar si el nombre de la máquina virtual, la IP del host y el tipo de solicitud están presentes y no son nulos
	nombreVM, nombrePresente := datos["nombreVM"].(string)

	if !nombrePresente || nombreVM == "" {
		log.Println("El tipo de solicitud y nombre de la máquina virtual son obligatorios")
		http.Error(w, "El tipo de solicitud y nombre de la máquina virtual son obligatorios", http.StatusBadRequest)
		return
	}

	datos["tipo_solicitud"] = "stop"
	// Encola las peticiones.
	config.GetMu().Lock()
	config.GetManagementQueue().Queue.PushBack(datos)
	config.GetMu().Unlock()

	// Envía una respuesta al cliente.
	confirmation := map[string]string{"mensaje": "Mensaje JSON para apagar MV recibido correctamente"}

	response := apiutilities.BuildGenericResponse(confirmation, "Success", "Mensaje JSON para apagar MV recibido correctamente")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func GetShhPrivateKeyHandler(w http.ResponseWriter, r *http.Request) {
	// extrae el path param con el correo
	vars := mux.Vars(r)
	nombre := vars["name"]
	var err error
	var virtualMachine entities.Maquina_virtual
	var keyPath string
	var filePath string
	var fileName string

	// verificar que el path param no esté vacío
	if nombre == "" {
		log.Println("No se especificó el nombre de la maquina virtual")
		http.Error(w, "El nombre de la máquina virtual es obligatorio", http.StatusBadRequest)
	}

	log.Println("Obteniendo la máquina virtual con nombre: ", nombre)

	virtualMachine, err = database.GetVM(nombre)

	if err != nil {
		http.Error(w, "No se encontró una maquina virtual con el nombre especificado", http.StatusBadRequest)
		return
	}

	if virtualMachine.Estado == "Apagado" {
		log.Println("La máquina virtual está apagada")
		http.Error(w, "La máquina virtual está apagada", http.StatusBadRequest)
		return
	}

	keyPath = "./keys/" + nombre

	err = systemutilities.ProcessSshPublicKeyConfiguration(keyPath, virtualMachine.Ip)

	if err != nil {
		log.Println("Error al generar la llave privada")
		http.Error(w, "Error al generar la llave privada", http.StatusBadRequest)
		return
	}

	filePath = keyPath + "/id_rsa_test"

	// Abrir el archivo
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "No se pudo abrir el archivo", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Obtener el nombre del archivo para la cabecera
	fileName = "id_rsa_test"

	// Establecer encabezados para forzar la descarga
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")

	// Escribir el archivo en la respuesta
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "No se pudo enviar el archivo", http.StatusInternalServerError)
		return
	}

	// Asegurarse de cerrar el archivo antes de eliminar la carpeta
	file.Close()

	// Eliminar la carpeta después de enviar el archivo
	if err := os.RemoveAll(keyPath); err != nil {
		log.Printf("Error al eliminar la carpeta %s: %v", keyPath, err)
	}

}
