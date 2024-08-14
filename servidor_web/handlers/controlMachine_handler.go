package handlers

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"AppWeb/Utilities"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var vmtemp Models.VirtualMachineTemp

func ControlMachine(c *gin.Context) {
	// Acceder a la sesión
	session := sessions.Default(c)
	email := session.Get("email")

	//TODO: Se debe adaptar para las sesiones de usuarios temporales
	/*ToDo: Se debería hacer un metodo aparte para ajustador lo de la verificación del usuario
	* Esto se utiliza en muchas parte pero primero debemos definir como se va a tratar este tipo de usuarios
	 */
	if email == nil {
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

	if email == nil {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	maquina_virtual, err := createVirualMachine(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se posible crear maquinva virutal" + err.Error()})
	}

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

	if Utilities.SendJSONMachineToServer(jsonData) {
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

// Función para crear máquina virtual decodificando sus atributos desde un json
func createVirualMachine(c *gin.Context) (Models.VirtualMachine, error) {
	var newVM Models.VirtualMachine

	// Decodificar JSON desde el cuerpo de la solicitud
	if err := json.NewDecoder(c.Request.Body).Decode(&newVM); err != nil {
		return Models.VirtualMachine{}, err
	}

	// Validar campos necesarios
	if newVM.Name == "" || newVM.Person_Email == "" || newVM.Ram == 0 || newVM.Cpu == 0 || newVM.Host_id == 0 {
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

	fmt.Println(nombre)

	err, respuesta, num := Utilities.PowerMachineServer(nombre, clientIP)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enceder la maquina" + err.Error()})
	}
	if num == 1 {
		textMessege := "¡" + respuesta + nombre + ". Por favor espere."
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"SuccessMessage": textMessege,
		})
	}
	if num == 2 {
		// Registro erróneo, muestra un mensaje de error en el HTML
		textMessege := " Error al Encender " + nombre + ". Intente de nuevo."
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"ErrorMessage": textMessege,
		})
	}

}

// ToDo: Aqui voy jajajaa
func DeleteMachine(c *gin.Context) {
	serverURL := fmt.Sprintf("http://%s:8081/json/deleteVM", Config.ServidorProcesamientoRoute)
	nombre := c.PostForm("vmnameDelete")

	payload := map[string]interface{}{
		"tipo_solicitud": "delete",
		"nombreVM":       nombre,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Registro exitoso, muestra un mensaje de éxito en el HTML
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"SuccessMessage": "Solicitud para eliminar la màquina enviada con èxito.",
		})
	} else {
		// Registro erróneo, muestra un mensaje de error en el HTML
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"ErrorMessage": "La solicitud para eliminar la màquina no fue exitosa. Intente de nuevo",
		})
	}
}

func ConfigMachine(c *gin.Context) {
	serverURL := fmt.Sprintf("http://%s:8081/json/modifyVM", Config.ServidorProcesamientoRoute)

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
	vmname := c.PostForm("vmnameConfig")
	memory, _ := strconv.Atoi(c.PostForm("memoryConfig"))
	cpu, _ := strconv.Atoi(c.PostForm("cpuConfig"))

	// Crear una estructura Maquina_virtual y convertirla a JSON
	Specifications := Maquina_virtual{Nombre: vmname, Ram: memory, Cpu: cpu, Persona_email: userEmail}

	fmt.Println(Specifications)

	payload := map[string]interface{}{
		"tipo_solicitud": "modify",
		"specifications": Specifications,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Crear una solicitud HTTP POST con el JSON como cuerpo
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	// Establecer el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realizar la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Registro exitoso, muestra un mensaje de éxito en el HTML
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"SuccessMessage": "Solicitud para modificar la màquina virtual enviada con èxito",
		})
	} else {
		// Registro erróneo, muestra un mensaje de error en el HTML
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{
			"ErrorMessage": "La solicitud para modificar la màquina virtual no tuvo èxito. Intente de nuevo",
		})
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
	machines, err := consultarMaquinas(userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver los datos en formato JSON
	c.JSON(http.StatusOK, machines)
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

func EnviarContenido(c *gin.Context) {
	var data struct {
		Contenido string `json:"contenido"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": data.Contenido, // Modifica esto según tus necesidades.
	})
}

// SEGUNDA ITERACION DEKTOP CLOUD
// TODO: Revisar metodo
func Checkhost(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email")

	if email == nil {
		// Si el usuario no está autenticado, redirige a la página de inicio de sesión
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userEmail := email.(string)
	idHostStr := c.PostForm("host")
	idHost, err := strconv.Atoi(idHostStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	fmt.Print(vmtemp)
	fmt.Print(idHost)

	// Obtener el idHost del formulario
	memoryint, _ := strconv.Atoi(vmtemp.Memory)
	cpuint, _ := strconv.Atoi(vmtemp.CPU)
	// Obtener la dirección IP del cliente
	clientIP := c.ClientIP()
	maquina_virtual := Maquina_virtual{

		Nombre:                         vmtemp.VMName,
		Sistema_operativo:              "linux",
		Distribucion_sistema_operativo: vmtemp.OS,
		Ram:                            memoryint,
		Cpu:                            cpuint,
		Persona_email:                  userEmail,
		Host_id:                        idHost}

	// Crear el objeto JSON con los datos del cliente
	payload := map[string]interface{}{
		"clientIP":       clientIP,
		"ubicacion":      idHost,
		"specifications": maquina_virtual,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	serverURL := fmt.Sprintf("http://%s:8081/json/checkhost", Config.ServidorProcesamientoRoute)

	// Realizar una solicitud POST al servidor remoto con los datos en formato JSON
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	// Establecer el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realizar la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode == http.StatusOK {
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{"SuccessMessage": "Solicitud para chequear maquina virtual enviada con éxito."})
	} else {
		c.HTML(http.StatusOK, "controlMachine.html", gin.H{"ErrorMessage": "Esta maquina virtual tiene problemas :(  selecciona otra por favor "})
	}
}

func Mvtemp(c *gin.Context) {

	// Deserializa el JSON recibido
	if err := c.ShouldBindJSON(&vmtemp); err != nil {
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos JSON inválidos",
		})
		return
	}
}
