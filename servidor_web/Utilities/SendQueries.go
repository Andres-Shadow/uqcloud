package Utilities

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"log"
	"net/http"
)

// Fucnión encargada de registrar cualquier tipo de elemento
func RegisterElements[T any](URL string, element T) error {
	//Crear una solicitud HTTP Post con el elemento como cuerpo
	client := &http.Client{}
	jsonData, err := json.Marshal(element)

	if err != nil {
		log.Println("Error al decodificar la estructura a JSON", err.Error())
		return err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error al enviar la solicitud HTTP", err.Error())
		return err
	}
	defer resp.Body.Close()

	//Verificar respuesta
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: La solicitud no fue exitosa")
		return fmt.Errorf("register element failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// Funcion generica para enviar cualquier tipo de peticiones
func SendRequest(method, url string, payload interface{}) (bool, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error al decodificar la estructura a JSON", err.Error())
		return false, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error al enviar la solicitud HTTP", err.Error())
		return false, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error: La solicitud no fue exitosa")
		body, _ := ioutil.ReadAll(resp.Body)
		return false, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	return true, nil
}
