package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"bytes"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"log"
	"net/http"
)

// Consultat maquinas virtuales asociadas a un email
func ConsultMachineFromServer(email string) ([]Models.VirtualMachine, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.VIRTUAL_MACHINE_URL)

	persona := Models.Person{Email: email}
	jsonData, err := json.Marshal(persona)
	if err != nil {
		log.Println("Error al deccodificar la estrcutura peronsa como json", err.Error())
		return nil, err
	}

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("GET", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return nil, err
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error al realizar la solicutad HTTP", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// Verifica la respuesta del servidor (resp.StatusCode) aquí si es necesario
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: La solicitud no fue exitosa", err.Error())
		return nil, errors.New("La solicitud al servidor no fue exitosa")
	}

	// Lee la respuesta del cuerpo de la respuesta HTTP
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error al leer el cuerpo de la respuesta", err.Error())
		return nil, err
	}

	var machines []Models.VirtualMachine

	// Decodifica los datos de respuesta en la variable machines.
	if err := json.Unmarshal(responseBody, &machines); err != nil {
		// Maneja el error de decodificación aquí
		log.Println("error al decodifcar el JSON", err.Error())
		return nil, fmt.Errorf("error al decodificar JSON de respuesta: %w", err)
	}

	return machines, nil
}

// Enviar creación de la VM al servidor
func CreateMachineFromServer(VM Models.VirtualMachine, clienteIp string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.VIRTUAL_MACHINE_URL)

	payload := map[string]interface{}{
		"specifications": VM,
		"clientIP":       clienteIp,
	}

	confirmacion, err := SendRequest("POST", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil
}

// Encender Maquina virtual
func PowerMachineFromServer(nombre string, clientIP string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.START_VM_URL)
	payload := map[string]interface{}{
		"tipo_solicitud": "start",
		"nombreVM":       nombre,
		"clientIP":       clientIP,
	}

	confirmacion, err := SendRequest("POST", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil

}

// Eliminar una Maquina virtual
func DeleteMachineFromServer(nombre string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.VIRTUAL_MACHINE_URL)

	payload := map[string]interface{}{
		"tipo_solicitud": "delete",
		"nombreVM":       nombre,
	}

	confirmacion, err := SendRequest("DELETE", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil
}

// Modifica las maquinas virtuales
func ConfigMachienFromServer(specifications Models.VirtualMachineTemp) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.VIRTUAL_MACHINE_URL)

	payload := map[string]interface{}{
		"tipo_solicitud": "modify",
		"specifications": specifications,
	}

	confirmacion, err := SendRequest("PUT", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil
}

// Consultar estado de la maquina virtual
func CheckStatusMachineFromServer(VM Models.VirtualMachine, clienteIp string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.CHECK_HOST_URL)

	// Crear el objeto JSON con los datos del cliente
	payload := map[string]interface{}{
		"clientIP":       clienteIp,
		"ubicacion":      VM.Host_id,
		"specifications": VM,
	}

	confirmacion, err := SendRequest("GET", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil

}
