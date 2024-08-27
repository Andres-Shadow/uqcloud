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
	// email := session.Get("email")

	// TODO: POR QUE NO SIRVE BIEN??????
	// log.Println(email)
	// if email != nil {
	// 	log.Println("Email invalido")
	// 	c.Redirect(http.StatusFound, "/mainPage")
	// 	return
	// }

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
		"email":       email,
		"contrasenia": password,
	}

	jsonData, err := json.Marshal(persona)
	if err != nil {
		log.Println("Error al decodificar: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	usuario, er := Utilities.SendInfoUserServer(jsonData)
	if er == nil {
		session := sessions.Default(c)
		session.Set("email", email)
		session.Set("nombre", usuario.Name)
		session.Set("apellido", usuario.LastName)
		session.Set("rol", usuario.Role)
		session.Save()

		log.Println("Usuario inicia sesion con exito")
		log.Printf("%+v\n", usuario)
		c.Redirect(http.StatusFound, "/mainPage")
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
// ToDo: Esta función puede cambiar en el futuro
func LoginTemp(c *gin.Context) {
	session := sessions.Default(c)
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.CREATE_GUEST_VM_URL)

	clientIP := c.ClientIP()
	distribucion := c.PostForm("osCreate")

	//Crea un mapa con la dirección IP del cliente
	data := map[string]string{
		"ip":           clientIP,
		"distribucion": distribucion,
	}

	// Convierte el mapa a JSON
	jsonBody, err := json.Marshal(data)
	if err != nil {
		// Maneja el error si la conversión falla
		log.Println("Error al convertir a JSON:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	// Crea una solicitud HTTP con el cuerpo JSON
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonBody))
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
		session.Set("rol", "Invitado")
		session.Save()

		c.Redirect(http.StatusSeeOther, "/mainpage/control-machine")
	} else {
		log.Println("No fue posible crear el usuario")
		c.Redirect(http.StatusNotFound, "/login")
	}

}

func Logout(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	// Eliminar la información de la sesión, incluyendo el email
	session.Delete("email")
	session.Save()

	// Redirigir al usuario a la página de inicio de sesión u otra página
	c.Redirect(http.StatusFound, "/")
}

/*TODO: Funcion no se utiliza revisar si se puede eliminar o simplemente comentar
func GuestLoginSend(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	if email == nil {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userEmail := email.(string)

	// Obtener los datos del formulario
	vmname := c.PostForm("vmnameCreate")
	if vmname == "" {
		// Si el nombre de la máquina virtual está vacío, mostrar un mensaje de error en el HTML
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"ErrorMessage": "El nombre de la máquina virtual no puede estar vacío.",
		})
		return
	}
	ditOs := c.PostForm("osCreate")
	memoryStr := c.PostForm("memoryCreate")
	memory, err := strconv.Atoi(memoryStr)
	cpuStr := c.PostForm("cpuCreate")
	cpu, _ := strconv.Atoi(cpuStr)
	os := "Linux"

	// Crear una estructura Account y convertirla a JSON
	maquina_virtual := Models.VirtualMachine{Name: vmname, Sistema_operativo: os, Distrubucion_SO: ditOs, Ram: memory, Cpu: cpu, Person_Email: userEmail}
	clientIP := c.ClientIP()

	payload := map[string]interface{}{
		"specifications": maquina_virtual,
		"clientIP":       clientIP,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sendJSONMachineToServer(jsonData) {
		// Registro exitoso, muestra un mensaje de éxito en el HTML
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"SuccessMessage": "Solicitud para crear màquina virtual enviada con èxito.",
		})
	} else {
		// Registro erróneo, muestra un mensaje de error en el HTML
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"ErrorMessage": "Error al enviar la solicitud para crear màquina virtual. Intente de nuevo",
		})
	}
}
*/
