package utilities

import (
	"fmt"
	"log"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

/*
Clase encargada de contener las funalidades asociadas a la realizacion de acciones
dentro de las maquinas virtuales, diferente a las utilidades encargadas de
realizar acciones de interaccion con las mismas
*/

/*
Funciòn que permite iniciar una màquina virtual en modo "headless", lo que indica que se inicia en segundo plano
para que el usuario de la màquina fìsica no se vea afectado
@nameVM Paràmetro que contiene el nombre de la màquina virtual a encender
*/

func StartVM(nameVM string, clientIP string) string {

	var err error
	var host models.Host
	var running bool
	// var ipAddress string
	// var restarted bool
	// var maxEspera time.Time
	var conf *ssh.ClientConfig

	// // Comando para encender la máquina virtual en segundo planto
	// var startVMHeadlessCommand string = "VBoxManage startvm " + "\"" + nameVM + "\"" + " --type headless"

	// //Comnado para encender la màquina virtual con GUI
	// var startVMGUICommand string = "VBoxManage startvm " + "\"" + nameVM + "\""

	//Comando para obtener la dirección IP de la máquina virtual
	// var getIpCommand string = "VBoxManage guestproperty get " + "\"" + nameVM + "\"" + " /VirtualBox/GuestInfo/Net/0/V4/IP"

	// //Comando para reiniciar la MV
	// var rebootCommand string = "VBoxManage controlvm " + "\"" + nameVM + "\"" + " reset"

	//Obtiene el objeto "maquina_virtual"
	maquinaVirtual, err := database.GetVM(nameVM)
	if err != nil {
		log.Println("Error al obtener la MV:", err)
		return "Error al obtener la MV"
	}

	//Obtiene el host en el cual està alojada la MV
	host, err = database.GetHost(maquinaVirtual.Host_id)
	if err != nil {
		log.Println("Error al obtener el host:", err)
		return "Error al obtener el host"
	}

	//Configura la conexiòn SSH con el host
	conf, err = ConfigureSSH(host.Hostname, config.GetPrivateKeyPath())
	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar SSH"
	}

	if running, err = IsRunning(nameVM, host.Ip, conf); err != nil {
		return logError("Error al obtener el estado de la MV:", err)
	} else if running {
		return "La máquina ya está encendida"
	}

	fmt.Println("Encendiendo la màquina " + nameVM + "...")

	return startVMSequence(nameVM, clientIP, host, conf)

	// ------------------------------- //

	// time.Sleep(10 * time.Second)

	// // Establece un temporizador de espera máximo de 2 minutos
	// maxEspera = time.Now().Add(2 * time.Minute)
	// restarted = false

	// for ipAddress == "" || ipAddress == "No value set!" || strings.HasPrefix(strings.TrimSpace(ipAddress), "169") {
	// 	if time.Now().Before(maxEspera) {
	// 		if ipAddress == "No value set!" {
	// 			time.Sleep(5 * time.Second) // Espera 5 segundos antes de intentar nuevamente
	// 			fmt.Println("Obteniendo dirección IP de la màquina " + nameVM + "...")
	// 		}
	// 		//Envìa el comando para obtener la IP
	// 		ipAddress, _ = SendSSHCommand(host.Ip, getIpCommand, conf)

	// 		ipAddress = strings.TrimSpace(ipAddress) //Elimina espacios en blanco al final de la cadena
	// 		ipParts := strings.Split(ipAddress, ":")
	// 		if len(ipParts) > 1 {
	// 			ipParts := strings.Split(ipParts[1], ".")
	// 			if strings.TrimSpace(ipParts[0]) == "169" {
	// 				ipAddress = strings.TrimSpace(ipParts[0])
	// 				time.Sleep(5 * time.Second) // Espera 5 segundos antes de intentar nuevamente
	// 				fmt.Println("Obteniendo dirección IP de la màquina " + nameVM + "...")
	// 			}
	// 		}

	// 	} else {
	// 		if restarted {
	// 			log.Println("No se logrò obtener la direcciòn IP de la màquina: " + nameVM)
	// 			//Actualiza el estado de la MV en la base de datos
	// 			err9 := database.UpdateVirtualMachineState(nameVM, "Apagado")
	// 			if err9 != nil {
	// 				log.Println("Error al realizar la actualizaciòn del estado", err9)
	// 				return "Error al realizar la actualizaciòn del estado"
	// 			}
	// 			return "No se logrò obtener la direcciòn IP, por favor contacte al administrador"
	// 		}
	// 		//Envìa el comando para reiniciar la MV
	// 		reboot, error := SendSSHCommand(host.Ip, rebootCommand, conf)
	// 		if error != nil {
	// 			log.Println("Error al reinciar la MV:", reboot)
	// 			return "Error al reinciar la MV"
	// 		}
	// 		fmt.Println("Reiniciando la máquina: " + nameVM)
	// 		maxEspera = time.Now().Add(2 * time.Minute) //Agrega dos minutos de tiempo màximo para obtener la IP cuando se reincia la MV
	// 		restarted = true
	// 	}
	// }

	// //Almacena la direccion ip de la maquina virtual
	// ipAddress = strings.TrimSpace(strings.TrimPrefix(ipAddress, "Value:"))

	// log.Println("Dirección IP de la máquina " + nameVM + ": " + ipAddress)

	// //Actualiza el estado de la MV en la base de datos
	// err9 := database.UpdateVirtualMachineState(nameVM, "Encendido")
	// if err9 != nil {
	// 	log.Println("Error al realizar la actualizaciòn del estado", err9)
	// 	return "Error al realizar la actualizaciòn del estado"
	// }
	// //Actualiza la direcciòn IP de la MV en la base de datos
	// err10 := database.UpdateVirtualMachineIP(nameVM, ipAddress)
	// if err10 != nil {
	// 	log.Println("Error al realizar la actualizaciòn de la IP", err10)
	// 	return "Error al realizar la actualizaciòn de la IP"
	// }
	// fmt.Println("Màquina encendida, la direcciòn IP es: " + ipAddress)
	// return ipAddress

}

// startVMSequence handles the steps to start the VM and retrieve its IP address.
func startVMSequence(nameVM, clientIP string, host models.Host, conf *ssh.ClientConfig) string {
	startVMCommand := getStartVMCommand(clientIP, nameVM)
	if err := executeStartCommand(host.Ip, startVMCommand, conf); err != nil {
		return logError("Error al enviar el comando para encender la MV:", err)
	}

	if err := database.UpdateVirtualMachineState(nameVM, "Procesando"); err != nil {
		return logError("Error al actualizar el estado de la MV:", err)
	}

	ipAddress, err := getVMIPAddress(nameVM, host, conf)
	if err != nil {
		return err.Error()
	}

	if err := finalizeVMStart(nameVM, ipAddress); err != nil {
		return err.Error()
	}

	return ipAddress
}

// getStartVMCommand selects the proper start command based on the client IP.
func getStartVMCommand(clientIP, nameVM string) string {
	if _, err := IsAHostIp(clientIP); err == nil {
		return "VBoxManage startvm " + "\"" + nameVM + "\""
	}
	return "VBoxManage startvm " + "\"" + nameVM + "\"" + " --type headless"
}

// executeStartCommand sends the command to start the VM.
func executeStartCommand(hostIP, command string, conf *ssh.ClientConfig) error {
	_, err := SendSSHCommand(hostIP, command, conf)
	return err
}

// getVMIPAddress attempts to retrieve the VM's IP address within a time limit.
func getVMIPAddress(nameVM string, host models.Host, conf *ssh.ClientConfig) (string, error) {

	fmt.Println("Obteniendo direcciòn IP de la màquina " + nameVM + "...")

	getIpCommand := "VBoxManage guestproperty get " + "\"" + nameVM + "\"" + " /VirtualBox/GuestInfo/Net/0/V4/IP"
	rebootCommand := "VBoxManage controlvm " + "\"" + nameVM + "\"" + " reset"
	maxWait := time.Now().Add(2 * time.Minute)
	restarted := false
	var ipAddress string

	for time.Now().Before(maxWait) {
		ipAddress, _ = SendSSHCommand(host.Ip, getIpCommand, conf)
		ipAddress = cleanIPAddress(ipAddress)

		if isValidIPAddress(ipAddress) {
			return ipAddress, nil
		}

		if time.Now().After(maxWait) {
			if restarted {
				return "", fmt.Errorf("no se logró obtener la dirección IP. Contacte al administrador")
			}
			if err := rebootVM(host.Ip, rebootCommand, conf); err != nil {
				return "", err
			}
			restarted = true
			maxWait = time.Now().Add(2 * time.Minute) // Extiende el tiempo después del reinicio
		}
		time.Sleep(5 * time.Second) // Espera entre intentos
	}

	return "", fmt.Errorf("no se logró obtener la dirección IP")
}

// cleanIPAddress removes unnecessary parts of the IP address string.
func cleanIPAddress(ipAddress string) string {
	return strings.TrimSpace(strings.TrimPrefix(ipAddress, "Value:"))
}

// isValidIPAddress checks if the IP address is valid.
func isValidIPAddress(ipAddress string) bool {
	return ipAddress != "" && !strings.HasPrefix(ipAddress, "169") && ipAddress != "No value set!"
}

// rebootVM reboots the virtual machine if needed.
func rebootVM(hostIP, rebootCommand string, conf *ssh.ClientConfig) error {
	_, err := SendSSHCommand(hostIP, rebootCommand, conf)
	if err != nil {
		log.Println("Error al reiniciar la MV:", err)
	}
	return err
}

// finalizeVMStart updates the database with the VM state and IP address.
func finalizeVMStart(nameVM, ipAddress string) error {
	if err := database.UpdateVirtualMachineState(nameVM, "Encendido"); err != nil {
		return fmt.Errorf("error al actualizar el estado de la MV: %v", err)
	}

	if err := database.UpdateVirtualMachineIP(nameVM, ipAddress); err != nil {
		return fmt.Errorf("error al actualizar la IP de la MV: %v", err)
	}

	log.Println("Máquina encendida, la dirección IP es:", ipAddress)
	return nil
}

// logError logs an error message and returns it as a string.
func logError(message string, err error) string {
	log.Println(message, err)
	return message
}

/* Funciòn que permite enviar el comando PowerOff para apagar una màquina virtual
@nameVM Paràmetro que contiene el nombre de la màquina virtual a apagar
@clientIP Paràmetro que contiene la direcciòn IP del cliente desde el cual se realiza la solicitud
*/

func TurnOffVM(nameVM string, clientIP string) string {

	//Obtiene la màquina vitual a apagar
	maquinaVirtual, err := database.GetVM(nameVM)
	if err != nil {
		log.Println("Error al obtener la MV:", err)
		return "Error al obtener la MV"
	}
	//Obtiene el host en el cual està alojada la MV
	host, err1 := database.GetHost(maquinaVirtual.Host_id)
	if err1 != nil {
		log.Println("Error al obtener el host:", err1)
		return "Error al obtener el host"
	}
	//Configura la conexiòn SSH con el host
	config, err2 := ConfigureSSH(host.Hostname, config.GetPrivateKeyPath())
	if err2 != nil {
		log.Println("Error al configurar SSH:", err2)
		return "Error al configurar SSH"
	}
	//Variable que contiene el estado de la MV (Encendida o apagada)
	running, err3 := IsRunning(nameVM, host.Ip, config)
	if err3 != nil {
		log.Println("Error al obtener el estado de la MV:", err3)
		return "Error al obtener el estado de la MV"
	}

	if !running { //En caso de que la MV estè apagada, no haga nada
		log.Println("La máquina ya está apagada")
		return "La máquina ya está apagada"
	} else {

		//Comando para apagar la màquina virtual
		powerOffCommand := "VBoxManage controlvm " + "\"" + nameVM + "\"" + " poweroff"

		fmt.Println("Apagando màquina " + nameVM + "...")
		//Actualza el estado de la MV en la base de datos

		err4 := database.UpdateVirtualMachineState(nameVM, "Procesando")
		if err4 != nil {
			log.Println("Error al realizar la actualizaciòn del estado", err4)
			return "Error al realizar la actualizaciòn del estado"
		}
		//Envìa el comando para apagar la MV a travès de un ACPI
		_, err5 := SendSSHCommand(host.Ip, powerOffCommand, config)
		if err5 != nil {
			log.Println("Error al enviar el comando para apagar la MV:", err5)
			return "Error al enviar el comando para apagar la MV"
		}
		// Establece un temporizador de espera máximo de 5 minutos
		maxEspera := time.Now().Add(5 * time.Minute)

		// Espera hasta que la máquina esté apagada o haya pasado el tiempo máximo de espera
		for time.Now().Before(maxEspera) {
			status, err6 := IsRunning(nameVM, host.Ip, config)
			if err6 != nil {
				log.Println("Error al obtener el estado de la MV:", err6)
				return "Error al obtener el estado de la MV"
			}
			if !status {
				break
			}
			// Espera un 1 segundo antes de volver a verificar el estado de la màquina
			time.Sleep(1 * time.Second)
		}

		//Consulta si la MV està encendida
		status, err7 := IsRunning(nameVM, host.Ip, config)
		if err7 != nil {
			log.Println("Error al obtener el estado de la MV:", err7)
			return "Error al obtener el estado de la MV"
		}
		if status {
			_, err8 := SendSSHCommand(host.Ip, powerOffCommand, config) //Envìa el comando para apagar la MV a travès de un Power Off
			if err8 != nil {
				log.Println("Error al enviar el comando para apagar la MV:", err8)
				return "Error al enviar el comando para apagar la MV"
			}
		}
		//Actualiza el estado de la MV en la base de datos
		err9 := database.UpdateVirtualMachineState(nameVM, "Apagado")
		if err9 != nil {
			log.Println("Error al realizar la actualizaciòn del estado", err9)
			return "Error al realizar la actualizaciòn del estado"
		}

		fmt.Println("Màquina apagada con èxito")
	}
	return ""
}
