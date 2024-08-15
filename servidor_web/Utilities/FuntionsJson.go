package Utilities

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
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
