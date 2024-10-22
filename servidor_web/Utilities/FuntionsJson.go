package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/goccy/go-json"
)

// Enviar información del usuario
func SendInfoUserServer(jsonData []byte) (Models.Person, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.LOGIN_URL)

	log.Println("sever URL:", serverURL)

	var usuario Models.Person

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error al crear la solicitud HTTP: ", err)
		return usuario, err
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error al realizar la solicitud", err)
		return usuario, err
	}
	defer resp.Body.Close()

	// Verifica la respuesta del servidor (resp.StatusCode) aquí si es necesario
	if resp.StatusCode != http.StatusOK {
		log.Println("Error al obtener la respuesta HTTP", resp.StatusCode)
		return usuario, errors.New("Error en la respuesta del servidor")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error al leer la respuesta del servidor: ", err)
		return usuario, err
	}
	var resultado map[string]interface{}

	if err := json.Unmarshal(responseBody, &resultado); err != nil {
		log.Println("Error al serializar la respuesta: ", err)
		return usuario, err
	}
	specsMap, _ := resultado["usuario"].(map[string]interface{})
	specsJSON, err := json.Marshal(specsMap)
	if err != nil {
		log.Println("Error al serializar el usuario ", err)
		return usuario, err
	}

	err = json.Unmarshal(specsJSON, &usuario)
	if err != nil {
		log.Println("Error al deserializar el usuario:", err)
		return usuario, err
	}

	return usuario, nil
}
