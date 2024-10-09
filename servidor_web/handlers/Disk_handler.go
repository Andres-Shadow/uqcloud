package handlers

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"AppWeb/Utilities"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreateDiskPage(c *gin.Context) {
	session := sessions.Default(c)

	c.HTML(http.StatusOK, "createDisk.html", gin.H{
		"email":    session.Get("email").(string),
		"nombre":   session.Get("nombre").(string),
		"apellido": session.Get("apellido").(string),
		"rol":      session.Get("rol").(uint8),
		// "hosts": hosts,
	})
}

// --------- FUNCIONES NUEVA PARA LA CREACIÓN Y REGISTRO DE DISK -------- //

func CreateNewDisk(c *gin.Context) {
	//Crear el Disk a partir de la solicutud
	newDisk, err := CreateDiskFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el JSON: " + err.Error()})
		return
	}

	//Registrar el disk
	// Definir la URL del servidor
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.DISK_VM_URL)

	if err := Utilities.RegisterElements(serverURL, newDisk); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registro el disco"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Disco creado correctamente"})
}

/*Funcion que se encarga de decodificar los parametros para crear un nuevo disco
 */
func CreateDiskFromRequest(c *gin.Context) (Models.Disk, error) {
	var newDisk Models.Disk

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	log.Printf("Request Body: %s", string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &newDisk); err != nil {
		log.Println("Error al decodificar el JSON---: " + err.Error())
		return Models.Disk{}, err
	}

	if newDisk.Name == "" || newDisk.Ruta_Ubicacion == "" || newDisk.Sistema_Operativo == "" ||
		newDisk.Arquitectura < 0 || newDisk.Host_id <= 0 {
		log.Println("Error existen campos vacios")
		return Models.Disk{}, errors.New("error existen campos vacios")
	}

	return newDisk, nil
}
