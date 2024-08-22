package handlers

import (
	"AppWeb/Models"
	"AppWeb/Utilities"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//var vmtemp Models.VirtualMachineTemp

func ControlMachine(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	//TODO: Se debe adaptar para las sesiones de usuarios temporales
	/*ToDo: Se debería hacer un metodo aparte para ajustador lo de la verificación del usuario
	* Esto se utiliza en muchas parte pero primero debemos definir como se va a tratar este tipo de usuarios
	 */
	log.Println(email)
	if email == nil {
		log.Println("Error: Email vacio/invalido")
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Recuperar o inicializar un arreglo de máquinas virtuales en la sesión del usuario
	machines, _ := Utilities.ConsultMachineFromServer(email.(string))

	hosts, _ := Utilities.CheckAvaibleHost()

	if sessionMachines, ok := session.Get("machines").([]Models.VirtualMachine); ok {
		machines = sessionMachines
	} else {
		// Inicializa un nuevo arreglo de máquinas si no existe en la sesión
		machines = []Models.VirtualMachine{}
	}

	// Agregar la variable booleana `machinesChange` a la sesión y establecerla en true
	session.Set("machinesChange", true)
	session.Save()

	machinesChange := session.Get("machinesChange")
	clientIP := c.ClientIP()
	showNewButton := false
	for _, host := range hosts {
		// Depuración
		if host.Ip == clientIP {
			showNewButton = true
			break
		}
	}
	c.HTML(http.StatusOK, "controlMachine.html", gin.H{
		"email":          email,
		"machines":       machines,
		"machinesChange": machinesChange,
		"hosts":          hosts,
		"showNewButton":  showNewButton,
		"clientIP":       clientIP,
	})
}

// Metodo que se encarga de crear y enviar la maquina virtual para su creación en el servidor
func MainSend(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	//Se podría hacer una función solo para verificar lo del Email, pero antes se debería definir como manejar los usuarios temporales
	log.Println(email)
	if email == nil {
		log.Println("Error: Email vacio/invalido")
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	maquina_virtual, err := createVirualMachine(c)
	clientIP := c.ClientIP()

	log.Printf("%+v\n", maquina_virtual)
	log.Printf("IP del cliente: %s", clientIP)

	//Se podría hacer una función para manjejar este tipo de errores ¿?
	if err != nil {
		log.Println("No es posible crear la maquina virutal, debido a alguen parametro invalido", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se posible crear maquinva virutal" + err.Error()})
		return
	}

	confirmacion, err := Utilities.CreateMachineFromServer(maquina_virtual, clientIP)

	if err != nil {
		log.Println("No es posible crear la maquina virutal, desde el servidor", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se posible crear maquina " + err.Error()})
		return
	}

	if confirmacion {
		// Registro exitoso, muestra un mensaje de éxito en el HTML
		log.Println("Maquina virutal creada exitosamente", maquina_virtual)
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"SuccessMessage": "Solicitud para crear màquina virtual enviada con èxito."})
	} else {
		log.Println("Error al enviar la soliciutd para crear maquina virtual")
		// Registro erróneo, muestra un mensaje de error en el HTML
		c.HTML(http.StatusInternalServerError, "controlMachine.html", gin.H{
			"ErrorMessage": "Error al enviar la solicitud para crear màquina virtual. Intente de nuevo"})
	}
}

// Función para crear máquina virtual decodificando sus atributos desde un json
func createVirualMachine(c *gin.Context) (Models.VirtualMachine, error) {
	var newVM Models.VirtualMachine

	// Decodificar JSON desde el cuerpo de la solicitud
	if err := json.NewDecoder(c.Request.Body).Decode(&newVM); err != nil {
		log.Println("Erro al decoficar el JSON para crear la maquina virutal", err.Error())
		return Models.VirtualMachine{}, err
	}

	// Validar campos necesarios
	log.Printf("%+v\n", newVM)
	if newVM.Name == "" || newVM.Person_Email == "" || newVM.Ram == 0 || newVM.Cpu == 0 || newVM.Host_id == 0 {
		log.Println("Error: Hay campos vacios")
		return Models.VirtualMachine{}, errors.New("missing required fields")
	}

	// Asignar valores predeterminados si es necesario
	if newVM.Sistema_operativo == "" {
		newVM.Sistema_operativo = "Linux"
	}

	return newVM, nil
}

// Funcion encargada de encender la maquinavirtual dado su nombre
func PowerMachine(c *gin.Context) {

	nombre := c.PostForm("nombreMaquina")
	clientIP := c.ClientIP()

	log.Println("Ip del cliente: ", clientIP)
	log.Println("Nombre de l VM:", nombre)

	confirmacion, err := Utilities.PowerMachineFromServer(nombre, clientIP)

	if err != nil {
		log.Println("Error al enceder la maquina virtual")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al enceder la maquina" + err.Error()})
		return
	}
	if confirmacion {
		log.Println("la maquina virtual ha sido encendida con exito ")
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{"SuccessMessage": "La máquina " + nombre + "Se esta encendiendo. Por favor espere"})
	} else { // Registro erróneo, muestra un mensaje de error en el HTML
		log.Println("Error al encender la maquina virtual")
		c.HTML(http.StatusInternalServerError, "controlMachine.html", gin.H{
			"ErrorMessage": "Error al encender la maquina " + nombre + "intente de nuevo"})
	}

}

// Metodo para eliminar una máquina virutal
func DeleteMachine(c *gin.Context) {
	nombre := c.PostForm("vmnameDelete")
	confirmacion, err := Utilities.DeleteMachineFromServer(nombre)

	log.Println("Nombre de la VM a elimianr: ", nombre)
	if err != nil {
		log.Println("Error al eliminar la maquina virtual", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la maquina" + err.Error()})
		return
	}
	if confirmacion {
		// Registro exitoso, muestra un mensaje de éxito en el HTML
		log.Println("La maquina virutal ha sido eliminada con exito")
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"SuccessMessage": "Solicitud para eliminar la màquina enviada con exito."})
	} else {
		// Registro erróneo, muestra un mensaje de error en el HTML
		log.Println("Error al enviar la solicito para eliminar la maquina virtual")
		c.HTML(http.StatusInternalServerError, "controlMachine.html", gin.H{
			"ErrorMessage": "La solicitud para eliminar la màquina no fue exitosa. Intente de nuevo"})
	}
}

// Metodo para configuarar una maqquina virtual
func ConfigMachine(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email").(string)

	log.Println(email)
	if email == " " {
		log.Println("Error: Email vacio")
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}
	var specifications Models.VirtualMachineTemp

	// Obtener los datos del formulario
	if err := c.BindJSON(&specifications); err != nil {
		// Manejar el error si el JSON no es válido o no coincide con la estructura
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		mv := fmt.Sprintf("Nombre: %s, RAM: %d, CPU: %d, Email: %s\n", specifications.VMName, specifications.Ram, specifications.CPU, specifications.Email)
		log.Println("Informacion de configuracion de la VM", mv)
		confirmacion, err := Utilities.ConfigMachienFromServer(specifications)

		if err != nil {
			log.Println("Error al configuracion de la VM", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al intentar configurar la maquina " + err.Error()})
		}
		if confirmacion {
			log.Println("Maquina virtual modificada con exito")
			// Registro exitoso, muestra un mensaje de éxito en el HTML
			c.HTML(http.StatusOK, "controlMachine.html", gin.H{
				"SuccessMessage": "Solicitud para modificar la màquina virtual enviada con èxito"})
		} else {
			// Registro erróneo, muestra un mensaje de error en el HTML
			log.Println("Error al configuracion de la VM")
			c.HTML(http.StatusInternalServerError, "controlMachine.html", gin.H{
				"ErrorMessage": "La solicitud para modificar la màquina virtual no tuvo èxito. Intente de nuevo"})
		}
	}
}

// No se utiliza jajaja
func GetMachines(c *gin.Context) {
	// Acceder a la sesión para obtener el email del usuario
	session := sessions.Default(c)
	email := session.Get("email")

	if email == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userEmail := email.(string)

	// Obtener los datos de las máquinas utilizando el email del usuario
	machines, err := Utilities.ConsultMachineFromServer(userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver los datos en formato JSON
	c.JSON(http.StatusOK, machines)
}

// SEGUNDA ITERACION DEKTOP CLOUD
func Checkhost(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email").(string)

	log.Println(email)
	if email == " " {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	idHostStr := c.PostForm("host")
	_, err := strconv.Atoi(idHostStr)
	if err != nil {
		log.Println("Error: formato de host invalido", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	var specifications Models.VirtualMachine

	log.Printf("%+v\n", specifications)
	// Obtener los datos del formulario
	if err := c.BindJSON(&specifications); err != nil {
		// Manejar el error si el JSON no es válido o no coincide con la estructura
		log.Println("Formado inadeciado", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inadecuado" + err.Error()})
	}

	// Obtener la dirección IP del cliente
	clienteIP := c.ClientIP()
	confirmacion, err := Utilities.CheckStatusMachineFromServer(specifications, clienteIP)

	log.Println("Ip del cliente", clienteIP)
	if err != nil {
		log.Println("Error al intentar configurar la VM", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al intentar configurar la maquina " + err.Error()})
	}

	// Verificar el código de estado de la respuesta
	if confirmacion {
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{"SuccessMessage": "Solicitud para chequear maquina virtual enviada con éxito."})
	} else {
		log.Println("Error al enviar la VM")
		c.HTML(http.StatusInternalServerError, "controlMachine.html", gin.H{"ErrorMessage": "Esta maquina virtual tiene problemas :(  selecciona otra por favor " + err.Error()})
	}
}

/* ToDo: Mira para que se utiliza esta función
func Mvtemp(c *gin.Context) {

	// Deserializa el JSON recibido
	if err := c.ShouldBindJSON(&vmtemp); err != nil {
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos JSON inválidos",
		})
		return
	}
}*/
