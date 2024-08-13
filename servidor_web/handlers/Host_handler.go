package handlers

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"AppWeb/Utilities"

	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
)

func CreateHostPage(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	rol := session.Get("rol")

	//TODO: Revisar si los roles pueden ser enum en vez de string (Revisar BASE DE DATOS)
	if rol != "Administrador" {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "createHost.html", nil)
}

// --------- FUNCIONES NUEVA PARA LA CREACIÓN Y REGISTRO DE HOST -------- //

// Función encargado de crear y registrar un nuevo
func CreateNewHost(c *gin.Context) {
	// Crear el host a partir de la solicitud
	newHost, err := CreateHostFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar JSON: " + err.Error()})
		return
	}

	// Registrar el host
	serverURL := fmt.Sprintf("http://%s:8081/json/addHost", Config.ServidorProcesamientoRoute)
	if err := Utilities.RegisterElements(serverURL, newHost); err != nil {
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
