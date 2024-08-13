package Utilities

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
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
