package handlers

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"AppWeb/Utilities"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateDiskPage(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email").(string)
	rol := session.Get("rol")

	if rol != "Administrador" {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	hosts, err := Utilities.ConsultHostsFromServer(email)

	if err != nil {
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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el JSON:" + err.Error()})
		return
	}

	//Registrar el disk
	// Definir la URL del servidor
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.DISK_VM_URL)

	if err := Utilities.RegisterElements(serverURL, newDisk); err != nil {
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
