package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

// Enviar contendio
func SendContent(c *gin.Context) {
	var data struct {
		Contenido string `json:"contenido"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": data.Contenido, // Modifica esto según tus necesidades.
	})
}

func UploadJSON(c *gin.Context) {
	// Obtener el archivo JSON del formulario
	file, err := c.FormFile("fileInput")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el archivo"})
		return
	}

	// Abrir el archivo
	jsonFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al abrir el archivo"})
		return
	}
	defer jsonFile.Close()

	// Decodificar el archivo JSON en un mapa
	var jsonData map[string]interface{}
	err = json.NewDecoder(jsonFile).Decode(&jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el archivo JSON"})
		return
	}

	// Enviar los datos JSON al cliente como respuesta
	c.JSON(http.StatusOK, jsonData)
}

// Enviar información del usuario
func SendInfoUserServer(jsonData []byte) (Models.Person, error) {
	serverURL := fmt.Sprintf("http://%s:8081/json/login", Config.ServidorProcesamientoRoute)
	var usuario Models.Person

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return usuario, err
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return usuario, err
	}
	defer resp.Body.Close()

	// Verifica la respuesta del servidor (resp.StatusCode) aquí si es necesario
	if resp.StatusCode != http.StatusOK {
		return usuario, errors.New("Error en la respuesta del servidor")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return usuario, err
	}
	var resultado map[string]interface{}

	if err := json.Unmarshal(responseBody, &resultado); err != nil {
		fmt.Println("Error al deserializar")
		return usuario, err
	}
	specsMap, _ := resultado["usuario"].(map[string]interface{})
	specsJSON, err := json.Marshal(specsMap)
	if err != nil {
		fmt.Println("Error al serializar el usuario:", err)
		return usuario, err
	}

	err = json.Unmarshal(specsJSON, &usuario)
	if err != nil {
		fmt.Println("Error al deserializar el usuario:", err)
		return usuario, err
	}

	return usuario, nil
}
