package handlers

import (
	"AppWeb/Utilities"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func MainPage(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")
	rol := session.Get("rol")

	//TODO: Si no existe el usuario temporal se le deberia crearle uno
	if email == nil {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Recuperar o inicializar un arreglo de máquinas virtuales en la sesión del usuario
	machines, _ := Utilities.ConsultMachineFromServer(email.(string))

	c.HTML(http.StatusOK, "controlMachine.html", gin.H{
		"email":    email,
		"machines": machines,
		"rol":      rol,
	})
}

// TODO: Esto se hace con un nav en html no debe ser una ruta
// func Scrollmenu(c *gin.Context) {

// 	// Acceder a la sesión
// 	session := sessions.Default(c)
// 	email := session.Get("email")
// 	rol := session.Get("rol")

// 	// Recuperar o inicializar un arreglo de máquinas virtuales en la sesión del usuario
// 	machines, _ := Utilities.ConsultMachineFromServer(email.(string))

// 	c.HTML(http.StatusOK, "scrollmenu.html", gin.H{
// 		"email":    email,
// 		"machines": machines,
// 		"rol":      rol,
// 	})
// }

// TODO: Cambiar nombre de la funcion a ingles
// TODO: Moverlo la función a otra clase
func ActualizacionesMaquinas(c *gin.Context) {

	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	// Obtén las máquinas actualizadas (por ejemplo, desde una base de datos)
	machines, err := Utilities.ConsultMachineFromServer(email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener actualizaciones de máquinas"})
		return
	}

	// Devuelve las máquinas en formato JSON
	c.JSON(http.StatusOK, machines)
}
