package utilities

import (
	"fmt"
	"log"
	"regexp"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"

	"time"

	"golang.org/x/crypto/ssh"
)

/*
Clase encargada de contener las funciones relacionadas con la gestion de maquinas virtuales
*/
var privateKeyPath string

func initPrivateKey(path string) {
	privateKeyPath = path
}

/*
Funciòn que verifica el tiempo de creaciòn de las màquinas de los usuarios invitados con el fin de determinar si se ha pasado o no del tiempo lìmite (2.5horas)
En caso de que una màquina se haya pasado del tiempo, se procederà a eliminarla.
*/

func CheckMachineTime(privateKey string) {

	initPrivateKey(privateKey)

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
			config, err2 := ConfigureSSH(host.Hostname, privateKeyPath)
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
	config, err2 := ConfigureSSH(host.Hostname, privateKeyPath)
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
		// err6 := database.DB.QueryRow("DELETE FROM maquina_virtual WHERE NOMBRE = ?", nameVM)
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
	if err != nil {
		log.Println("Error al realizar la consulta: ", err)
		return existe, err
	}
	return existe, err
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
