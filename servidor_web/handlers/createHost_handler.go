package handlers

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"AppWeb/Utilities"
	"log"

	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func CreateHostPage(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	rol := session.Get("rol")

	if rol != "Administrador" {
		log.Println("El usuario no es administrador no puede acceder")
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "createHost.html", nil)
}

// --------- FUNCIONES NUEVA PARA LA CREACIÓN Y REGISTRO DE HOST -------- //

// Función encargado de crear y registrar un nuevo
func CreateNewHost(c *gin.Context) {
	newHost, err := CreateHostFromRequest(c)

	log.Printf("%+v\n", newHost)

	if err != nil {
		log.Println("Error al decodificar el JSON del host", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar JSON: " + err.Error()})
		return
	}

	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.HOST_URL)
	log.Println(serverURL)

	if err := Utilities.RegisterElements(serverURL, newHost); err != nil {
		log.Println("Error al registrar el host", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el host: " + err.Error()})
		return
	}

	// Mostrar un mensaje de éxito
	c.HTML(http.StatusOK, "createHost.html", gin.H{"message": "Host creado correctamente"})
}

// Funcion que se encarga de descodificar los parametros para crear un nuevo host
func CreateHostFromRequest(c *gin.Context) (Models.Host, error) {
	var newHost Models.Host

	// Decodificar JSON desde el cuerpo de la solicitud
	if err := json.NewDecoder(c.Request.Body).Decode(&newHost); err != nil {
		return Models.Host{}, err
	}

	return newHost, nil
}
