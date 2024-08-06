package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	config "servidor_procesamiento/Procesador/Config"
	utilities "servidor_procesamiento/Procesador/Utilities"
)

/*
Clase encargada de contener los manejadores que responden a las imagenes de docker
*/

// Funcion que responde al endpoint de crear un contenedor a partir de una imagen de dockerhub
func CreateImageDockerHubHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método HTTP sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	// Ingresar Datos de la Imagen Docker

	imagen := payload["imagen"].(string)
	version := payload["version"].(string)

	ip := payload["ip"].(string)
	hostname := payload["hostname"].(string)

	mensaje := utilities.CreateImageDockerHub(imagen, version, ip, hostname)

	fmt.Println(mensaje)

	// Respondemos con la lista de máquinas virtuales en formato JSON
	response := map[string]string{"mensaje": mensaje}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Funcion que responde al endpoint de crear un contenedor a partir de un archivo .tar
func CreateImageDockerTarHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método HTTP sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	nombreArchivo := payload["archivo"].(string)

	ip := payload["ip"].(string)
	hostname := payload["hostname"].(string)

	mensaje := utilities.CreateImageTarFile(nombreArchivo, ip, hostname)

	fmt.Println(mensaje)

	// Respondemos con la lista de máquinas virtuales en formato JSON
	response := map[string]string{"mensaje": mensaje}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Funcion que responde al endpoint encargado de crear imagenes docker a partir de un archivo Dockerfiel
func CreateImageDockerfileHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método HTTP sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	nombreArchivo := payload["archivo"].(string)
	nombreImagen := payload["nombreImagen"].(string)

	ip := payload["ip"].(string)
	hostname := payload["hostname"].(string)

	mensaje := utilities.CreateImageDockerFile(nombreArchivo, nombreImagen, ip, hostname)

	fmt.Println(mensaje)

	// Respondemos con la lista de máquinas virtuales en formato JSON
	response := map[string]string{"mensaje": mensaje}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Funcion que responde al endpoint encargado de eliminar una imagen de docker
func DeleteDockerImageHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método HTTP sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	config.GetMu().Lock()
	config.GetDocker_imagesQueue().Queue.PushBack(payload)
	config.GetMu().Unlock()

	// Respondemos con la lista de máquinas virtuales en formato JSON
	response := map[string]string{"mensaje": "Se elimino la Imagen"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Funcion que responde al endpoint encargado de eliminar un contenedor de docker
func CheckVirtualMachineDockerImagesHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método HTTP sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	ip := payload["ip"].(string)
	hostname := payload["hostname"].(string)

	imagenes, err := utilities.ListImages(ip, hostname)

	if err != nil && err.Error() != "Fallo en la ejecucion" {
		fmt.Println(err)
		log.Println("Error al enviar datos")
		return
	}

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(imagenes)
}

// Funcion que responde al endpoint encargado de crear una maquina virtual que contenga contenedores
func CreateDockerHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método HTTP sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	imagen := payload["imagen"].(string)
	comando := payload["comando"].(string)

	ip := payload["ip"].(string)
	hostname := payload["hostname"].(string)

	mensaje := utilities.CreateContainer(imagen, comando, ip, hostname)

	// Respondemos con la lista de máquinas virtuales en formato JSON
	response := map[string]string{"mensaje": mensaje}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Funcion que responde al endpoint encargada de administrar las imagenes de docker en el sistema
func ManageDockerImagesHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método HTTP sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	config.GetMu().Lock()
	config.GetDocker_contenedorQueue().Queue.PushBack(payload)
	config.GetMu().Unlock()

	// Respondemos con la lista de máquinas virtuales en formato JSON
	response := map[string]string{"mensaje": "Comando Exitoso"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Funcion que responde al endpoint encargado de verificar los contenedores de una maquina virtual
func CheckContainersHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método HTTP sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el cuerpo JSON de la solicitud
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	ip := payload["ip"].(string)
	hostname := payload["hostname"].(string)

	contenedor, err := utilities.ListContainers(ip, hostname)

	if err != nil && err.Error() != "Fallo en la ejecucion" {
		fmt.Println(err)
		log.Println("Error al enviar datos")
		return
	}

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contenedor)

}
