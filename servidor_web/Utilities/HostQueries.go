package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/goccy/go-json"
)

// Funcion encargada de consultar la cantidad de host asociados al email de una persona
func ConsultHostsFromServer(email string) ([]Models.Host, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.HOST_URL)
	log.Println(serverURL)

	log.Println(email)
	jsonData, err := json.Marshal(email)

	if err != nil {
		log.Println("error al convertir la estructura persona a JSON", err.Error())
		return nil, fmt.Errorf("error al convertir persona a JSON: %w", err)
	}

	req, err := http.NewRequest("GET", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("error al crear la solicitud HTTP", err.Error())
		return nil, fmt.Errorf("error al crear solicitud HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Realizar la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error al realizar la solictud HTTP", err.Error())
		return nil, fmt.Errorf("error al realizar solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: La solicutd no fue exitosa", resp.StatusCode)
		return nil, fmt.Errorf("la solicitud al servidor no fue exitosa: %s", resp.Status)
	}

	// Leer el cuerpo de la respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error al leer el cuerpo la solicitud HTTP", err.Error())
		return nil, fmt.Errorf("error al leer cuerpo de la respuesta: %w", err)
	}

	// Decodificar los datos de respuesta en la variable hosts
	var hosts []Models.Host
	if err := json.Unmarshal(responseBody, &hosts); err != nil {
		log.Println("error al decodificar el JSON de la respuesta", err.Error())
		return nil, fmt.Errorf("error al decodificar JSON de respuesta: %w", err)
	}

	return hosts, nil
}

// Consultar lo host disponibles
func CheckAvaibleHost() ([]Models.Host, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.HOSTS_URL)
	log.Println(serverURL)

	persona := Models.Person{Email: "123"}
	jsonData, err := json.Marshal(persona)
	if err != nil {
		log.Println("error al convertir la estructura persona a JSON", err.Error())
		return nil, err
	}

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("GET", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("error al crear la solicitud HTTP", err.Error())
		return nil, err
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error al realizar la solicitud HTTP", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// Verifica la respuesta del servidor (resp.StatusCode) aquí si es necesario
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: La solicitud no fue exitosa", resp.StatusCode)
		return nil, errors.New("La solicitud al servidor no fue exitosa")
	}

	// Lee la respuesta del cuerpo de la respuesta HTTP
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error al leer el cuerpo la solicitud HTTP", err.Error())
		return nil, err
	}

	var hosts []Models.Host

	// Decodifica los datos de respuesta en la variable machines.
	if err := json.Unmarshal(responseBody, &hosts); err != nil {
		log.Println("error al decodificar la JSON de la respuesta", err.Error())
		return nil, fmt.Errorf("error al decodificar JSON de respuesta: %w", err)
	}

	return hosts, nil
}

func GetHostsFromServer() ([]Models.HostConsult, error) {

	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.HOSTS_URL)
	log.Println(serverURL)

	resp, err := http.Get(serverURL)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error en la solicitud: estado HTTP no es 200 OK")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer el cuerpo de la respuesta: %w", err)
	}

	var hosts []Models.HostConsult
	if err := json.Unmarshal(body, &hosts); err != nil {
		return nil, fmt.Errorf("error al decodificar el JSON: %w", err)
	}

	return hosts, nil
}
