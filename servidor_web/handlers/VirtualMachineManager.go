package handlers

import (
	"AppWeb/Models"
	"AppWeb/Utilities"
	"encoding/json"
	"errors"
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
	log.Println("El email es:", email)
	if email == nil {
		log.Println("Error: Email vacio/invalido")
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
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

	// Agregar la variable booleana `machinesChange` a la sesión y establecerla en true
	session.Set("machinesChange", true)
	session.Save()

	machinesChange := session.Get("machinesChange")
	clientIP := c.ClientIP()
	showNewButton := false

	if len(hosts) <= 0 {
		log.Println("No existen host para realizar esta operación, se deben registrar los host")
		return
	}
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

func CreateMachinePage(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	//TODO: Se debe adaptar para las sesiones de usuarios temporales
	/*TODO: Se debería hacer un metodo aparte para ajustador lo de la verificación del usuario
	* Esto se utiliza en muchas parte pero primero debemos definir como se va a tratar este tipo de usuarios
	 */
	log.Println(email)
	if email == nil {
		log.Println("Error: Email vacio/invalido")
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	hosts, _ := Utilities.CheckAvaibleHost()

	// for _, host := range hosts {
	// 	log.Println("---hostName: ", host.Name)
	// 	log.Println("---host: ", host)
	// }

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
	c.HTML(http.StatusOK, "create-machine.html", gin.H{
		"email":          email,
		"machinesChange": machinesChange, // TODO: Esta si será necesaria?
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
		// Registro exitoso, muestra un mensaje de éxito en el HTML
		log.Println("Maquina virtual creada exitosamente", maquina_virtual)
		c.HTML(http.StatusOK, "create-machine.html", gin.H{
			"SuccessMessage": "Solicitud para crear màquina virtual enviada con èxito."})
	} else {
		log.Println("Error al enviar la soliciutd para crear maquina virtual")
		// Registro erróneo, muestra un mensaje de error en el HTML
		c.HTML(http.StatusInternalServerError, "create-machine.html", gin.H{
			"ErrorMessage": "Error al enviar la solicitud para crear màquina virtual. Intente de nuevo"})
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

	// TODO: IMPLEMENTAR COMO DISTRIBUIR LO DEL ALEATORIO (ASIGNACION DE RECURSOS miamor)
	// TODO: Crear una funcion que se traiga la CANTIDAD de hosts disponibles, no la info de los hosts
	hosts, _ := Utilities.CheckAvaibleHost()
	if arrivedVM.Host_id == 0 {
		// rand.Seed(time.Now().UnixNano())
		// randomHost := rand.Intn(len(hosts))
		// log.Println("RANDOM host: ", hosts[randomHost])

		// newVM.Host_id = hosts[randomHost].Id
		// newVM.Hostname = hosts[randomHost].Name
		// log.Println("RANDOM host name: ", hosts[randomHost].Name)

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

		log.Printf("%+v\n", specifications)

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
