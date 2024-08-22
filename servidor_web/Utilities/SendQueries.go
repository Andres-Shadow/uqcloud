package Utilities

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
)

// Fucnión encargada de registrar cualquier tipo de elemento
func RegisterElements[T any](URL string, element T) error {
	//Crear una solicitud HTTP Post con el elemento como cuerpo
	client := &http.Client{}
	jsonData, err := json.Marshal(element)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Verificar respuesta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("register element failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// Funcion generica para enviar cualquier tipo de peticiones
func SendRequest(method, url string, payload interface{}) (bool, error) {
	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("error marshaling payload: %w", err)
	}

	// Create a new HTTP request with the given method, URL, and JSON payload
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("error creating request: %w", err)
	}

	// Set the content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the status code indicates success
	if resp.StatusCode != http.StatusOK {
		// Read the response body for error details
		body, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	return true, nil
}
