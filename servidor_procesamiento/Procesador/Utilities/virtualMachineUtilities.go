package utilities

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models/Entities"
	"strings"

	"time"

	"golang.org/x/crypto/ssh"
)

/*
Clase encargada de contener las funciones relacionadas con la gestion de maquinas virtuales
*/

func CreateVirtualMachineFromSpecifications(specs map[string]interface{}) {
	jsonData, _ := json.Marshal(specs) //Se codifica en formato JSON

	var decodedPayload map[string]interface{}
	err := json.Unmarshal(jsonData, &decodedPayload) //Se decodifica para meterlo en la cola
	if err != nil {
		fmt.Println("Error al decodificar el JSON:", err)
		// Manejar el error según tus necesidades
		return
	}
	// Encola la peticiòn
	config.GetMu().Lock()
	config.GetMaquina_virtualQueue().Queue.PushBack(decodedPayload)
	config.GetMu().Unlock()
}

/*
Funciòn que verifica el tiempo de creaciòn de las màquinas de los usuarios invitados con el fin de determinar si se ha pasado o no del tiempo lìmite (2.5horas)
En caso de que una màquina se haya pasado del tiempo, se procederà a eliminarla.
*/

func CheckMachineTime(privateKey string) {

	// Obtiene todas las máquinas virtuales de la base de datos
	maquinas, err := database.GetGuestMachines()
	if err != nil {
		log.Println("Error al obtener las máquinas virtuales:", err)
		return
	}

	// Obtiene la hora actual
	horaActual := time.Now().UTC()

	for _, maquina := range maquinas {
		// Calcula la diferencia de tiempo entre la hora actual y la fecha de creación de la máquina
		diferencia := horaActual.Sub(maquina.Fecha_creacion)

		// Verifica si la máquina ha excedido su tiempo de duración, en este caso: 2horas 20minutos
		if diferencia > (2*time.Hour + 20*time.Minute) {

			//Obtiene el host en el cual està alojada la MV
			host, err1 := database.GetHost(maquina.Host_id)
			if err1 != nil {
				log.Println("Error al obtener el host:", err)
				return
			}
			//Configura la conexiòn SSH con el host
			config, err2 := ConfigureSSH(host.Hostname, config.GetPrivateKeyPath())
			if err2 != nil {
				log.Println("Error al configurar SSH:", err2)
				return
			}

			//Variable que contiene el estado de la MV (Encendida o apagada)
			running, err3 := IsRunning(maquina.Nombre, host.Ip, config)
			if err3 != nil {
				log.Println("Error al obtener el estado de la MV:", err3)
				return
			}
			if running {
				TurnOffVM(maquina.Nombre, "")
			}
			DeleteVM(maquina.Nombre)
		}
	}
}

/*
Funciòn que se encarga de realizar la configuraciòn SSH con el host
@user Paràmetro que contiene el nombre del usuario al cual se va a conectar
@privateKeyPath Paràmetro que contiene la ruta de la llave privada SSH
*/
func ConfigureSSH(user string, privateKeyPath string) (*ssh.ClientConfig, error) {
	authMethod, err := privateKeyFile(privateKeyPath)
	fmt.Println("authMethod", authMethod)
	fmt.Println(privateKeyFile(privateKeyPath))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config, nil
}

/*
	Esta funciòn verifica si una màquina virtual està encendida

@nameVM Paràmetro que contiene el nombre de la màquina virtual a verificar
@hostIP Paràmetro que contiene la direcciòn Ip del host en el cual està la MV
@return Retorna true si la màquina està encendida o false en caso contrario
*/
func IsRunning(nameVM string, hostIP string, config *ssh.ClientConfig) (bool, error) {

	//Comando para saber el estado de una màquina virtual
	command := "VBoxManage showvminfo " + "\"" + nameVM + "\"" + " | findstr /C:\"State:\""
	running := false

	salida, err := SendSSHCommand(hostIP, command, config)
	if err != nil {
		log.Println("Error al ejecutar el comando para obtener el estado de la màquina:", err)
		return running, err
	}

	// Expresión regular para buscar el estado (running)
	regex := regexp.MustCompile(`State:\s+(running|powered off)`)
	matches := regex.FindStringSubmatch(salida)

	// matches[1] contendrá "running" o "powered off" dependiendo del estado
	if len(matches) > 1 {
		estado := matches[1]
		if estado == "running" {
			running = true
		}
	}
	return running, nil
}

/* Funciòn que permite enviar los comandos necesarios para eliminar una màquina virtual
@nameVM Paràmetro que contiene el nombre de la màquina virtual a eliminar
*/

func DeleteVM(virtualMachineName string) string {

	//Obtiene el objeto "maquina_virtual"
	maquinaVirtual, err := database.GetVM(virtualMachineName)
	if err != nil {
		log.Println("Error al obtener la MV:", err)
		return "Error al obtener  la MV"
	}
	//Obtiene el host en el cual està alojada la MV
	host, err1 := database.GetHost(maquinaVirtual.Host_id)
	if err1 != nil {
		log.Println("Error al obtener el host:", err)
		return "Error al obtener el host"
	}
	//Configura la conexiòn SSH con el host
	config, err2 := ConfigureSSH(host.Hostname, config.GetPrivateKeyPath())
	if err2 != nil {
		log.Println("Error al configurar SSH:", err2)
		return "Error al configurar SSH"
	}

	//Comando para desconectar el disco de la MV
	disconnectCommand := "VBoxManage storageattach " + "\"" + virtualMachineName + "\"" + " --storagectl hardisk --port 0 --device 0 --medium none"

	//Comando para eliminar la MV
	deleteCommand := "VBoxManage unregistervm " + "\"" + virtualMachineName + "\"" + " --delete"

	//Variable que contiene el estado de la MV (Encendida o apagada)
	running, err3 := IsRunning(virtualMachineName, host.Ip, config)
	if err3 != nil {
		log.Println("Error al obtener el estado de la MV:", err3)
		return "Error al obtener el estado de la MV"
	}
	if running {
		fmt.Println("Debe apagar la màquina para eliminarla")
		return "Debe apagar la màquina para eliminarla"

	} else {
		//Envìa el comando para desconectar el disco de la MV
		_, err4 := SendSSHCommand(host.Ip, disconnectCommand, config)
		if err4 != nil {
			log.Println("Error al desconectar el disco de la MV:", err4)
			return "Error al desconectar el disco de la MV"
		}
		//Envìa el comando para eliminar la MV del host
		_, err5 := SendSSHCommand(host.Ip, deleteCommand, config)
		if err5 != nil {
			log.Println("Error al eliminar la MV:", err5)
			return "Error al eliminar la MV"
		}
		//Elimina la màquina virtual de la base de datos

		err6 := database.DeleteVirtualMachine(virtualMachineName)
		if err6 == nil {
			log.Println("Error al eliminar el registro de la base de datos: ", err6)
			return "Error al eliminar el registro de la base de datos"
		}
		//Calcula los recursos usados del host, descontando los recursos liberados por la MV eliminada
		ram_host_usada := host.Ram_usada - maquinaVirtual.Ram
		cpu_host_usada := host.Cpu_usada - maquinaVirtual.Cpu
		//Actualiza los recursos usados del host en la base de datos
		err7 := database.UpdateHostRamAndCPU(host.Id, ram_host_usada, cpu_host_usada)
		if err7 == nil {
			log.Println("Error al actualizar los recursos usados del host en la base de datos: ", err7)
			return "Error al actualizar los recursos usados del host en la base de datos"
		}
	}
	fmt.Println("Màquina eliminada correctamente")
	return "Màquina eliminada correctamente"
}

/*
Funciòn que permite conocer si ya existe o no una màquina virtual en la base de datos con el nombre proporcionado.
@nameVM Paràmetro que representa el nombre de la màquina virtual a buscar
@Return Retorna true si ya existe una MV con ese nombre, o false en caso contrario
*/
func ExistVM(virtualMachineName string) (bool, error) {

	//err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM maquina_virtual WHERE nombre = ?)", nameVM).Scan(&existe)
	existe, err := database.ExistVirtualMachine(virtualMachineName)
	if err != nil && existe {
		log.Println("Error al realizar la consulta: ", err)
		return existe, err
	}
	return existe, nil
}

// Función que imprime las especificaciones de una máquina virtual.
func PrintVirtualMachine(specs models.Maquina_virtual, isCreateVM bool) {

	// Imprime las especificaciones recibidas.
	fmt.Printf("-------------------------\n")
	fmt.Printf("Nombre de la Máquina: %s\n", specs.Nombre)
	fmt.Printf("Sistema Operativo: %s\n", specs.Sistema_operativo)
	fmt.Printf("Distribuciòn SO: %s\n", specs.Distribucion_sistema_operativo)
	fmt.Printf("Memoria Requerida: %d Mb\n", specs.Ram)
	fmt.Printf("CPU Requerida: %d núcleos\n", specs.Cpu)

}

// Funcion para actualizar la base de datos con las maquinas virtuales que estan realmente disponibles en los hosts
func UpdateVirtualMachinesActualStatus() {
	// Se obtienen todas las maquinas virtuales de la base de datos
	maquinas, err := database.GetAllVirtualMachines()
	if err != nil {
		log.Println("Error al obtener las maquinas virtuales:", err)
		return
	}

	log.Println("Maquinas virtuales registradas la base de datos: ", maquinas)

	/* Se comparan con las maquinas virtuales que estan realmente en los hosts utilizando VBoxManage remotamente
	y se eliminan de la base de datos las maquinas virtuales que no estan realmente en los hosts*/
	for _, maquina := range maquinas {
		// Se verifica si la maquina virtual esta realmente en el host
		if !ExistVirtualMachineInHost(maquina.Nombre) {
			log.Println("La maquina virtual ", maquina.Nombre, " no esta realmente en el host, se eliminara de la base de datos")
			// Se elimina la maquina virtual de la base de datos
			database.DeleteVirtualMachine(maquina.Nombre)
		}
	}

}

// Función que permite verificar si una máquina virtual se encuentra realmente en alguno de los hosts registrados
func ExistVirtualMachineInHost(nombreVM string) bool {
	// Obtiene la máquina virtual
	maquinaVirtual, err := database.GetVM(nombreVM)
	if err != nil {
		log.Println("Error al obtener la MV:", err)
		return false
	}
	// Obtiene el host en el cual está alojada la MV
	host, err1 := database.GetHost(maquinaVirtual.Host_id)
	if err1 != nil {
		log.Println("Error al obtener el host:", err1)
		return false
	}
	// Configura la conexión SSH con el host
	config, err2 := ConfigureSSH(host.Hostname, config.GetPrivateKeyPath())
	if err2 != nil {
		log.Println("Error al configurar SSH:", err2)
		return false
	}

	command := "VBoxManage list vms"
	// Obtiene la lista de máquinas virtuales en el host
	salida, err := SendSSHCommand(host.Ip, command, config)
	if err != nil {
		log.Println("Error al obtener la lista de MVs:", err)
		return false
	}

	// Si la salida contiene el nombre de la máquina virtual, entonces la máquina virtual existe en el host
	return strings.Contains(salida, nombreVM)
}
