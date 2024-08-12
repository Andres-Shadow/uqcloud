package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// TODO: Modificarlo para solo usuarios admin
func SigninPage(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	if email != nil {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/mainPage")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// TODO: Eliminar la opción de registro
func Signin(c *gin.Context) {
	// Obtener los datos del formulario
	nombre := c.PostForm("nombre")
	apellido := c.PostForm("apellido")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Crear una estructura Account y convertirla a JSON
	persona := Persona{Nombre: nombre, Apellido: apellido, Email: email, Contrasenia: password}
	jsonData, err := json.Marshal(persona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sendRegisterJSONToServer(jsonData) {
		// Registro exitoso, redirige a la página de login con un indicador de éxito
		c.Redirect(http.StatusFound, "/login?success=true")
	} else {
		// Registro erróneo, muestra un mensaje de error en el HTML
		c.HTML(http.StatusOK, "login.html", gin.H{
			"ErrorMessage": "El registro ha fallado. Inténtalo de nuevo.",
		})
	}
}

func sendRegisterJSONToServer(jsonData []byte) bool {
	serverURL := fmt.Sprintf("http://%s:8081/json/signin", ServidorProcesamientoRoute)

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Verifica la respuesta del servidor (resp.StatusCode) aquí si es necesario
	if resp.StatusCode != http.StatusOK {
		return false
	} else {
		return true
	}
}
