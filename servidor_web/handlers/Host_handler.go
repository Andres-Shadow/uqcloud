package handlers

import (
	"AppWeb/Config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func CreateHost(c *gin.Context) {
	// Definir la URL del servidor
	serverURL := fmt.Sprintf("http://%s:8081/json/addHost", Config.ServidorProcesamientoRoute)

	// Obtener los datos del formulario
	nombreHost := c.PostForm("nameHost")
	ipHost := c.PostForm("ipHost")
	macHost := c.PostForm("macHost")
	adapHost := c.PostForm("adapHost")
	soHost := c.PostForm("soHost")
	hostnameHost := c.PostForm("hostnameHost")
	ramHostStr := c.PostForm("ramHost")
	ramHost, _ := strconv.Atoi(ramHostStr)
	cpuHostStr := c.PostForm("cpuHost")
	cpuHost, _ := strconv.Atoi(cpuHostStr)
	almaceHostStr := c.PostForm("almaceHost")
	almaceHost, _ := strconv.Atoi(almaceHostStr)
	sshHost := c.PostForm("sshHost")

	// TODO: hacer este codigo en un packer model
	// Crear un objeto Host con los datos del formulario
	host := Host{
		Nombre:               nombreHost,
		Ip:                   ipHost,
		Mac:                  macHost,
		Adaptador_red:        adapHost,
		Sistema_operativo:    soHost,
		Hostname:             hostnameHost,
		Ram_total:            ramHost,
		Cpu_total:            cpuHost,
		Almacenamiento_total: almaceHost,
		Ruta_llave_ssh_pub:   sshHost,
	}

	// Serializar el objeto host como JSON
	jsonData, err := json.Marshal(host)
	if err != nil {
		// Manejar el error, por ejemplo, responder con un error HTTP
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al serializar el objeto Host"})
		return
	}

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		// Manejar el error, por ejemplo, responder con un error HTTP
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la solicitud HTTP"})
		return
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Manejar el error, por ejemplo, responder con un error HTTP
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar la solicitud HTTP"})
		return
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		// Manejar el error, por ejemplo, responder con un error HTTP
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en la respuesta del servidor"})
		return
	}

	// Mostrar un mensaje de éxito en la página HTML
	c.HTML(http.StatusOK, "createHost.html", gin.H{"message": "Host creado correctamente"})
}

// --------- FUNCIONES NUEVA PARA LA CREACIÓN Y REGISTRO DE HOST -------- //

//Función encargado de crear y registrar un nuevo host
func CreateNewHost(c *gin.Context) {
	// Crear el host a partir de la solicitud
	newHost, err := CreateHostFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar JSON: " + err.Error()})
		return
	}

	// Registrar el host
	if err := RegisterHost(newHost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el host: " + err.Error()})
		return
	}

	// Mostrar un mensaje de éxito
	c.HTML(http.StatusOK, "createHost.html", gin.H{"message": "Host creado correctamente"})
}

//Función que se encarga de registrar un nuevo host, luego de que este sea creado, se encarga de hacer una petición al servidor de procesamiento
func RegisterHost(host Host) error {
	serverURL := fmt.Sprintf("http://%s:8081/json/addHost", Config.ServidorProcesamientoRoute)

	// Crear una solicitud HTTP POST con el objeto Host como cuerpo
	client := &http.Client{}
	jsonData, err := json.Marshal(host)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Realizar la solicitud HTTP
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta del servidor: %s", resp.Status)
	}

	return nil
}

//Funcion que se encarga de descodificar los parametros para crear un nuevo host
func CreateHostFromRequest(c *gin.Context) (Host, error) {
	var newHost Host

	// Decodificar JSON desde el cuerpo de la solicitud
	if err := json.NewDecoder(c.Request.Body).Decode(&newHost); err != nil {
		return Host{}, err
	}

	return newHost, nil
}
