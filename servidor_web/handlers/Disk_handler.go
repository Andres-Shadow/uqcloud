package handlers

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"AppWeb/Utilities"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreateDiskPage(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)
	rol := session.Get("rol")

	if rol != "Administrador" {
		log.Println("El usuario no es administrador no puede acceder")
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	hosts, err := Utilities.ConsultHostsFromServer(email)

	if err != nil {
		log.Println("Error al obtener el host", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar los host: " + err.Error()})
	}

	c.HTML(http.StatusOK, "createDisk.html", gin.H{
		"email": email,
		"hosts": hosts,
	})
}

// --------- FUNCIONES NUEVA PARA LA CREACIÓN Y REGISTRO DE DISK -------- //

func CreateNewDisk(c *gin.Context) {
	//Crear el Disk a partir de la solicutud
	newDisk, err := CreateDiskFromRequest(c)

	log.Printf("%+v\n", newDisk)
	if err != nil {
		log.Println("errot al decodificar el JSON del disco: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el JSON:" + err.Error()})
		return
	}

	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.DISK_VM_URL)
	log.Println(serverURL)

	if err := Utilities.RegisterElements(serverURL, newDisk); err != nil {
		log.Println("Error al registro el disco: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registro el disk"})
		return
	}
	c.HTML(http.StatusOK, "newDisk.html", gin.H{"message": "Disk Creado correctamente"})
}

/*Funcion que se encarga de decodificar los parametros para crear un nuevo disco
 */
func CreateDiskFromRequest(c *gin.Context) (Models.Disk, error) {
	var newDisk Models.Disk

	//Decodificar JSON DESDE EL CUERPO DE LA SOLICITUD
	if err := json.NewDecoder(c.Request.Body).Decode(&newDisk); err != nil {
		return Models.Disk{}, err
	}

	return newDisk, nil
}
