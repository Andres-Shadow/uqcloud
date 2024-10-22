package handlers

import (
	"AppWeb/Config"
	"AppWeb/Utilities"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginAdminPage(c *gin.Context) {
	session := sessions.Default(c)

	errorMessage := session.Get("loginError")
	session.Delete("loginError")
	session.Save()

	c.HTML(http.StatusOK, "login.html", gin.H{
		"ErrorMessage": errorMessage,
	})
}

func AdminLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	persona := map[string]string{
		"usr_email":    email,
		"usr_password": password,
	}

	jsonData, err := json.Marshal(persona)
	if err != nil {
		log.Println("Error al decodificar: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	usuario, er := Utilities.SendInfoUserServer(jsonData)
	log.Println("Usuario: ", usuario)
	if er == nil {
		session := sessions.Default(c)
		session.Set("authenticated", true)

		session.Set("email", email)
		session.Set("nombre", usuario.Nombre)
		session.Set("apellido", usuario.Apellido)
		session.Set("rol", usuario.Rol)

		log.Println("Usuario inicia sesion con exito")
		log.Printf("%+v\n", usuario)

		session.Save()

		c.Redirect(http.StatusFound, "/auth-admin/dashboard")
	} else {
		log.Println("Credenciales invalidas/Usuario no encontrado: ", err)
		session := sessions.Default(c)
		session.Set("loginError", ErrorMessage())
		session.Save()
		c.Redirect(http.StatusFound, "/admin")
	}
}

func ErrorMessage() string {
	return "Credenciales incorrectas. Inténtalo de nuevo."
}

// Funcion para el Login Temporal de un usuario
func LoginTemp(c *gin.Context) {
	session := sessions.Default(c)
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.TEMP_USER_ACCOUNT)

	// Crea una solicitud HTTP con el cuerpo JSON
	req, err := http.NewRequest("POST", serverURL, nil)
	if err != nil {
		// Maneja el error si la creación de la solicitud falla
		log.Println("Error al crear la solicitud HTTP:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error al realizar la solicitud HTTP:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Lee el cuerpo de la respuesta
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error al leer el cuerpo de la respuesta:", err)
			return
		}

		// Convierte el cuerpo de la respuesta en un mapa
		var data map[string]string
		if err := json.Unmarshal(responseBody, &data); err != nil {
			fmt.Println("Error al decodificar el JSON:", err)
			return
		}

		// Accede a los datos del mapa
		mensaje := data["mensaje"]
		log.Println("usuario creado con exito", mensaje)

		session.Set("email", mensaje)
		session.Set("nombre", "Usuario")
		session.Set("apellido", "Temporal")
		session.Set("rol", uint8(0))
		session.Save()

		c.Redirect(http.StatusSeeOther, "/mainpage/control-machine/create-machine")
	} else {
		log.Println("No fue posible crear el usuario")
		c.Redirect(http.StatusNotFound, "/")
	}

}

func QuickMachine(c *gin.Context) {
	session := sessions.Default(c)
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.QUICK_VM)

	//obtener la ip del cliente
	clienteIp := c.ClientIP()

	//hacer el mapa con la ip del usuario
	var userSpecs = make(map[string]interface{})

	//adaptar el mapa con la ip del usuario
	userSpecs["ip"] = clienteIp

	// Crea una solicitud HTTP con el cuerpo JSON
	jsonData, err := json.Marshal(userSpecs)
	if err != nil {
		log.Println("Error al convertir userSpecs a JSON:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	req, err := http.NewRequest("POST", serverURL, bytes.NewReader(jsonData))
	if err != nil {
		// Maneja el error si la creación de la solicitud falla
		log.Println("Error al crear la solicitud HTTP:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error al realizar la solicitud HTTP:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Lee el cuerpo de la respuesta
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error al leer el cuerpo de la respuesta:", err)
			return
		}

		// Convierte el cuerpo de la respuesta en un mapa
		var data map[string]string
		if err := json.Unmarshal(responseBody, &data); err != nil {
			fmt.Println("Error al decodificar el JSON:", err)
			return
		}

		// Accede a los datos del mapa
		mensaje := data["mensaje"]
		log.Println("usuario creado con exito", mensaje)

		session.Set("email", mensaje)
		session.Set("nombre", "Usuario")
		session.Set("apellido", "Temporal")
		session.Set("rol", uint8(0))
		session.Save()

		confirmacion, _ := Utilities.VerifyMachineCreated("QuickGuest", mensaje)

		if !confirmacion {
			log.Println("No fue posible crear la maquina")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor, no se pudo crear la máquina rápida"})
			return
		}

		// c.Redirect(http.StatusSeeOther, "/mainpage/control-machine")
		c.JSON(http.StatusOK, gin.H{"mensaje": mensaje})

	} else {
		log.Println("No fue posible crear el usuario")
		c.Redirect(http.StatusNotFound, "/")
	}
}

func Logout(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)

	// Eliminar toda la información de la sesión
	session.Clear()
	session.Save()

	// Redirigir al usuario a la página de inicio de sesión u otra página
	c.Redirect(http.StatusFound, "/")
}
