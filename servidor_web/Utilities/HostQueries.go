package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

// Funcion encargada de consultar la cantidad de host asociados al email de una persona
func ConsulHostsFromServer(email string) ([]Models.Host, error) {
	URL := fmt.Sprintf("http://%s:8081/json/consultHost", Config.ServidorProcesamientoRoute)
	jsonData, err := json.Marshal(email)
	if err != nil {
		return nil, fmt.Errorf("error al convertir persona a JSON: %w", err)
	}

	// Crear una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
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
