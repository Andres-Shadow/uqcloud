package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"bytes"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"io/ioutil"
	"net/http"
)

// Funcion encargada de consultar la cantidad de host asociados al email de una persona
func ConsultHostsFromServer(email string) ([]Models.Host, error) {
	URL := fmt.Sprintf("http://%s:8081/api/v1/host", Config.ServidorProcesamientoRoute)
	jsonData, err := json.Marshal(email)
	if err != nil {
		return nil, fmt.Errorf("error al convertir persona a JSON: %w", err)
	}

	// Crear una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("GET", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error al crear solicitud HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Realizar la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("la solicitud al servidor no fue exitosa: %s", resp.Status)
	}

	// Leer el cuerpo de la respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer cuerpo de la respuesta: %w", err)
	}

	// Decodificar los datos de respuesta en la variable hosts
	var hosts []Models.Host
	if err := json.Unmarshal(responseBody, &hosts); err != nil {
		return nil, fmt.Errorf("error al decodificar JSON de respuesta: %w", err)
	}

	return hosts, nil
}

// Consultar lo host disponibles
func CheckAvaibleHost() ([]Models.Host, error) {
	serverURL := fmt.Sprintf("http://%s:8081/api/v1/hosts", Config.ServidorProcesamientoRoute)

	persona := Models.Person{Email: "123"}
	jsonData, err := json.Marshal(persona)
	if err != nil {
		return nil, err
	}

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Verifica la respuesta del servidor (resp.StatusCode) aquí si es necesario
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("La solicitud al servidor no fue exitosa")
	}

	// Lee la respuesta del cuerpo de la respuesta HTTP
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hosts []Models.Host

	// Decodifica los datos de respuesta en la variable machines.
	if err := json.Unmarshal(responseBody, &hosts); err != nil {
		// Maneja el error de decodificación aquí
	}

	return hosts, nil
}
