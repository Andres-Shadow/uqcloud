package handlers

import (
	"AppWeb/Models"
	"AppWeb/Utilities"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//var vmtemp Models.VirtualMachineTemp

func ControlMachine(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	log.Println("El email es:", email)
	if email == nil {
		log.Println("Error: Email vacio/invalido")
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Recuperar o inicializar un arreglo de máquinas virtuales en la sesión del usuario
	log.Println("Se ha envidado el email para consultar sus host")
	machines, _ := Utilities.ConsultMachineFromServer(email.(string))

	log.Println("Se ha procesado la solicitud")
	hosts, _ := Utilities.CheckAvaibleHost()

	if sessionMachines, ok := session.Get("machines").([]Models.VirtualMachine); ok {
		log.Println("Se han encontrado las maquinas asociadas al usuario")
		machines = sessionMachines
	} else {
		// Inicializa un nuevo arreglo de máquinas si no existe en la sesión
		log.Println("No se han encontrado las maquinas asociadas al usuario")
		machines = []Models.VirtualMachine{}
	}

	session.Save()
	clientIP := c.ClientIP()

	if len(hosts) <= 0 {
		log.Println("No existen host para realizar esta operación, se deben registrar los host")
		return
	}

	c.HTML(http.StatusOK, "controlMachine.html", gin.H{
		"email":    session.Get("email").(string),
		"rol":      session.Get("rol").(uint8),
		"machines": machines,
		"hosts":    hosts,
		"clientIP": clientIP,
	})
}

func CreateMachinePage(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	log.Println(email)
	if email == nil {
		log.Println("Error: Email vacio/invalido")
		// Si no se ha logueado (usando /admin o usando los botones de creación de vm, entonces redirijalo al index)
		c.Redirect(http.StatusFound, "/")
		return
	}

	hosts, _ := Utilities.CheckAvaibleHost()

	session.Save()
	clientIP := c.ClientIP()

	c.HTML(http.StatusOK, "create-machine.html", gin.H{
		"email":    session.Get("email").(string),
		"rol":      session.Get("rol").(uint8),
		"hosts":    hosts,
		"clientIP": clientIP,
	})
}

// Metodo que se encarga de crear y enviar la maquina virtual para su creación en el servidor
func MainSend(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	log.Println(email)
	if email == nil {
		log.Println("Error: Email vacio/invalido")
		c.Redirect(http.StatusFound, "/")
		return
	}

	maquina_virtual, err := createVirtualMachine(c)
	clientIP := c.ClientIP()

	log.Printf("%+v\n", maquina_virtual)
	log.Printf("IP del cliente: %s", clientIP)

	//Se podría hacer una función para manjejar este tipo de errores ¿?
	if err != nil {
		log.Println("No es posible crear la maquina virtual, debido a alguen parametro invalido", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "No es posible crear maquina virtual" + err.Error()})
		return
	}

	confirmacion, err := Utilities.CreateMachineFromServer(maquina_virtual, clientIP)

	if err != nil {
		log.Println("No es posible crear la maquina virtual, desde el servidor", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se posible crear maquina " + err.Error()})
		return
	}

	if confirmacion {
		log.Println("Maquina virtual creada exitosamente", maquina_virtual)
		c.JSON(http.StatusOK, gin.H{})
	} else {
		log.Println("Error al enviar la soliciutd para crear maquina virtual")
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
}

// Función para crear máquina virtual decodificando sus atributos desde un json
func createVirtualMachine(c *gin.Context) (Models.VirtualMachineTemp, error) {
	var arrivedVM Models.VirtualMachine
	var newVM Models.VirtualMachineTemp

	// Decodificar JSON desde el cuerpo de la solicitud
	if err := json.NewDecoder(c.Request.Body).Decode(&arrivedVM); err != nil {
		log.Println("Error al decoficar el JSON para crear la maquina virtual", err.Error())
		return Models.VirtualMachineTemp{}, err
	}

	// verifiquemos pues si si es como estoy diciendo con el tag del struct de la vm
	log.Println("New Virtual Machine: ", arrivedVM)

	// Validar campos necesarios
	log.Printf("%+v\n", arrivedVM)
	if arrivedVM.Name == "" || arrivedVM.Person_Email == "" || arrivedVM.Ram == 0 || arrivedVM.Cpu == 0 || arrivedVM.Distribucion_SO == "" {
		log.Println("Error: Hay campos vacios")
		return Models.VirtualMachineTemp{}, errors.New("missing required fields")
	}

	// TODO: Crear una funcion que se traiga la CANTIDAD de hosts disponibles, no la info de los hosts
	hosts, _ := Utilities.CheckAvaibleHost()
	if arrivedVM.Host_id == -1 {
		arrivedVM.Hostname = "roundrobin"

	} else if arrivedVM.Host_id == 0 {
		arrivedVM.Hostname = "aleatorio"

	} else {
		// Al servidor se le envia "vm_hostname", el cual lo necesita para usar "GetHostByName"
		for _, host := range hosts {
			if host.Id == arrivedVM.Host_id {
				arrivedVM.Hostname = host.Name
				continue
			}
		}
	}

	// Asignar valores predeterminados si es necesario
	if arrivedVM.Sistema_operativo == "" {
		arrivedVM.Sistema_operativo = "linux" // o cualquiera si sae
	}

	// Mapeado de Maquina a MaquinaTemp
	newVM.Name = arrivedVM.Name
	newVM.Ram = arrivedVM.Ram
	newVM.Cpu = arrivedVM.Cpu
	newVM.Hostname = arrivedVM.Hostname
	newVM.Person_Email = arrivedVM.Person_Email
	newVM.Sistema_operativo = arrivedVM.Sistema_operativo
	newVM.Distribucion_SO = arrivedVM.Distribucion_SO

	return newVM, nil
}

// Funcion encargada de encender la maquinavirtual dado su nombre
func StartMachine(c *gin.Context) {

	var req Models.StateMachineRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request | Se envió mal el nombre de la vm"})
		return
	}
	nombre := req.NombreMaquina
	clientIP := c.ClientIP()

	log.Println("Ip del cliente: ", clientIP)
	log.Println("Nombre de la VM:", nombre)

	confirmacion, err := Utilities.StartMachineFromServer(nombre, clientIP)

	if err != nil {
		log.Println("Error al enceder la maquina virtual")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al enceder la maquina" + err.Error()})
		return
	}
	if confirmacion {
		log.Println("la maquina virtual ha sido encendida con exito ")

		c.JSON(http.StatusOK, gin.H{
			"SuccessMessage": "La máquina " + nombre + " se está encendiendo. Por favor espere",
		})
	} else { // Registro erróneo, muestra un mensaje de error en el HTML
		log.Println("Error al encender la maquina virtual")
		c.JSON(http.StatusInternalServerError, gin.H{
			"ErrorMessage": "Error al encender la maquina " + nombre + "intente de nuevo"})
	}

}

// Funcion encargada de apagar la maquinavirtual dado su nombre
func StopMachine(c *gin.Context) {

	var req Models.StateMachineRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request | Se envió mal el nombre de la vm"})
		return
	}
	nombre := req.NombreMaquina
	clientIP := c.ClientIP()

	log.Println("Nombre de la VM:", nombre)
	log.Println("Ip del cliente: ", clientIP)

	confirmacion, err := Utilities.StopMachineFromServer(nombre, clientIP)

	if err != nil {
		log.Println("Error al apagar la maquina virtual")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al apagar la maquina" + err.Error()})
		return
	}
	if confirmacion {
		log.Println("la maquina virtual ha sido apagada con exito ")
		session := sessions.Default(c)

		c.JSON(http.StatusOK, gin.H{
			"SuccessMessage": "La máquina " + nombre + " se está apagando. Por favor espere",
			"email":          session.Get("email").(string),
			"rol":            session.Get("rol").(uint8),
		})
	} else { // Registro erróneo, muestra un mensaje de error en el HTML
		log.Println("Error al apagar la maquina virtual")
		c.JSON(http.StatusInternalServerError, gin.H{
			"ErrorMessage": "Error al apagar la maquina " + nombre + "intente de nuevo"})
	}

}

// Metodo para eliminar una máquina virutal
func DeleteMachine(c *gin.Context) {
	var req Models.StateMachineRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request | Se envió mal el nombre de la vm"})
		return
	}
	nombre := req.NombreMaquina
	confirmacion, err := Utilities.DeleteMachineFromServer(nombre)

	log.Println("Nombre de la VM a elimianr: ", nombre)
	if err != nil {
		log.Println("Error al eliminar la maquina virtual: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la maquina" + err.Error()})
		return
	}
	if confirmacion {
		// Registro exitoso, muestra un mensaje de éxito en el HTML
		log.Println("La maquina virutal ha sido eliminada con exito")
		c.JSON(http.StatusOK, gin.H{
			"SuccessMessage": "Solicitud para eliminar la màquina enviada con exito."})
	} else {
		// Registro erróneo, muestra un mensaje de error en el HTML
		log.Println("Error al enviar la solicito para eliminar la maquina virtual")
		c.JSON(http.StatusInternalServerError, gin.H{
			"ErrorMessage": "La solicitud para eliminar la màquina no fue exitosa. Intente de nuevo"})
	}
}

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
		c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}

	// Devolver los datos en formato JSON
	c.JSON(http.StatusOK, machines)
}

func ConnectionMachine(c *gin.Context) {

	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")
	machineIP := c.Query("machineIP")
	machineName := c.Query("machineName")

	// Recuperar o inicializar un arreglo de máquinas virtuales en la sesión del usuario
	log.Println("Se ha envidado el email para consultar sus host")
	machines, _ := Utilities.ConsultMachineFromServer(email.(string))

	log.Println("Se ha procesado la solicitud")
	hosts, _ := Utilities.CheckAvaibleHost()

	if sessionMachines, ok := session.Get("machines").([]Models.VirtualMachine); ok {
		log.Println("Se han encontrado las maquinas asociadas al usuario")
		machines = sessionMachines
	} else {
		// Inicializa un nuevo arreglo de máquinas si no existe en la sesión
		log.Println("No se han encontrado las maquinas asociadas al usuario")
		machines = []Models.VirtualMachine{}
	}

	session.Save()
	clientIP := c.ClientIP()

	if len(hosts) <= 0 {
		log.Println("No existen host para realizar esta operación, se deben registrar los host")
		return
	}

	c.HTML(http.StatusOK, "sshMachine.html", gin.H{
		"email":       session.Get("email").(string),
		"rol":         session.Get("rol").(uint8),
		"machines":    machines,
		"hosts":       hosts,
		"clientIP":    clientIP,
		"machineIP":   machineIP,
		"machineName": machineName,
	})
}
