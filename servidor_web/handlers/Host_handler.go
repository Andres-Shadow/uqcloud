package handlers

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"AppWeb/Utilities"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"strconv"

	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func CreateHostPage(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "createHost.html", gin.H{
		"email":    session.Get("email").(string),
		"nombre":   session.Get("nombre").(string),
		"apellido": session.Get("apellido").(string),
		"rol":      session.Get("rol").(uint8),
	})
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
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.HOST_URL)
	if err := Utilities.RegisterElements(serverURL, newHost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el host: " + err.Error()})
		return
	}

	// Mostrar un mensaje de éxito
	c.JSON(http.StatusOK, gin.H{"message": "Host creado correctamente"})
}

// Funcion que se encarga de descodificar los parametros para crear un nuevo host
func CreateHostFromRequest(c *gin.Context) (Models.Host, error) {
	var newHost Models.Host

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	log.Printf("Request Body: %s", string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &newHost); err != nil {
		log.Println("Error al decodificar el JSON---: " + err.Error())
		return Models.Host{}, err
	}

	if newHost.Name == "" || newHost.Ip == "" || newHost.Hostname == "" || newHost.Ram_total <= 0 || newHost.Cpu_total <= 0 ||
		newHost.Almacenamiento_total <= 0 || newHost.Adaptador_red == "" || newHost.Sistema_operativo == "" {
		log.Println("Error existen campos vacios")
		return Models.Host{}, errors.New("Error existen campos vacios")

	}

	return newHost, nil
}

func GetHosts(c *gin.Context) {
	hostDTO, err := Utilities.GetHostsFromServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los hosts"})
		return
	}

	for _, host := range hostDTO.Data {
		log.Println("HOST GETHOSTS: ", host)
	}

	c.JSON(http.StatusOK, hostDTO.Data)
}

func DeleteHost(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de ID inválido"})
		return
	}

	var deletedHosts []int
	var errors []string

	// Iterar sobre cada ID y hacer la solicitud DELETE al servidor externo
	for _, id := range ids {
		serverURL := fmt.Sprintf("http://%s:%s%s"+strconv.Itoa(id), Config.ServidorProcesamientoRoute, Config.PUERTO, Config.HOST_URL)

		req, err := http.NewRequest("DELETE", serverURL, bytes.NewBuffer(nil))
		if err != nil {
			errors = append(errors, "Error al crear la solicitud para el ID "+strconv.Itoa(id))
			continue
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errors = append(errors, "Error al enviar la solicitud para el ID "+strconv.Itoa(id))
			continue
		}
		defer resp.Body.Close()

		// Verificar el código de estado de la respuesta
		if resp.StatusCode == http.StatusOK {
			deletedHosts = append(deletedHosts, id)
		} else {
			errors = append(errors, "No se pudo eliminar el host con ID "+strconv.Itoa(id))
		}
	}

	// Devolver los resultados
	if len(errors) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"deleted": deletedHosts,
			"errors":  errors,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"deleted": deletedHosts, "message": "Todos los hosts eliminados exitosamente"})
	}
}
