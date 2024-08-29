package utilities

import (
	"fmt"
	"log"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	"strings"
	"time"
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

	//Obtiene el objeto "maquina_virtual"
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
	//Configura la conexiòn SSH con el host aaaaa
	config, err2 := ConfigureSSH(host.Hostname, config.GetPrivateKeyPath()) //RERVISAR QUE SI LE LLEGUE
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

	if running {
		TurnOffVM(nameVM, clientIP) //En caso de que la MV ya estè encendida, entonces se invoca el mètodo para apagar la MV
		return ""
	} else {
		fmt.Println("Encendiendo la màquina " + nameVM + "...")

		// Comando para encender la máquina virtual en segundo planto
		startVMHeadlessCommand := "VBoxManage startvm " + "\"" + nameVM + "\"" + " --type headless"

		//Comnado para encender la màquina virtual con GUI
		startVMGUICommand := "VBoxManage startvm " + "\"" + nameVM + "\""

		_, er := IsAHostIp(clientIP) //Verifica si la solicitud se està realizando desde un host registrado en la BD
		if er == nil {
			//Envìa el comando para encender la MV con GUI
			_, err4 := SendSSHCommand(host.Ip, startVMGUICommand, config)
			if err4 != nil {
				log.Println("Error al enviar el comando para encender la MV:", err4)
				return "Error al enviar el comando para encender la MV"
			}
		} else {
			//Envìa el comando para encender la MV en segundo plano
			_, err4 := SendSSHCommand(host.Ip, startVMHeadlessCommand, config)
			if err4 != nil {
				log.Println("Error al enviar el comando para encender la MV:", err4)
				return "Error al enviar el comando para encender la MV"
			}
		}

		fmt.Println("Obteniendo direcciòn IP de la màquina " + nameVM + "...")
		//Actualiza el estado de la MV en la base de datos
		err5 := database.UpdateVirtualMachineState(nameVM, "Procesando")
		if err5 != nil {
			log.Println("Error al realizar la actualizaciòn del estado", err5)
			return "Error al realizar la actualizaciòn del estado"
		}
		// Espera 10 segundos para que la máquina virtual inicie
		time.Sleep(10 * time.Second)

		//Comando para obtener la dirección IP de la máquina virtual
		getIpCommand := "VBoxManage guestproperty get " + "\"" + nameVM + "\"" + " /VirtualBox/GuestInfo/Net/0/V4/IP"
		//Comando para reiniciar la MV
		rebootCommand := "VBoxManage controlvm " + "\"" + nameVM + "\"" + " reset"

		var ipAddress string

		// Establece un temporizador de espera máximo de 2 minutos
		maxEspera := time.Now().Add(2 * time.Minute)
		restarted := false

		for ipAddress == "" || ipAddress == "No value set!" || strings.HasPrefix(strings.TrimSpace(ipAddress), "169") {
			if time.Now().Before(maxEspera) {
				if ipAddress == "No value set!" {
					time.Sleep(5 * time.Second) // Espera 5 segundos antes de intentar nuevamente
					fmt.Println("Obteniendo dirección IP de la màquina " + nameVM + "...")
				}
				//Envìa el comando para obtener la IP
				ipAddress, _ = SendSSHCommand(host.Ip, getIpCommand, config)

				ipAddress = strings.TrimSpace(ipAddress) //Elimina espacios en blanco al final de la cadena
				ipParts := strings.Split(ipAddress, ":")
				if len(ipParts) > 1 {
					ipParts := strings.Split(ipParts[1], ".")
					if strings.TrimSpace(ipParts[0]) == "169" {
						ipAddress = strings.TrimSpace(ipParts[0])
						time.Sleep(5 * time.Second) // Espera 5 segundos antes de intentar nuevamente
						fmt.Println("Obteniendo dirección IP de la màquina " + nameVM + "...")
					}
				}

			} else {
				if restarted {
					log.Println("No se logrò obtener la direcciòn IP de la màquina: " + nameVM)
					//Actualiza el estado de la MV en la base de datos
					err9 := database.UpdateVirtualMachineState(nameVM, "Apagado")
					if err9 != nil {
						log.Println("Error al realizar la actualizaciòn del estado", err9)
						return "Error al realizar la actualizaciòn del estado"
					}
					return "No se logrò obtener la direcciòn IP, por favor contacte al administrador"
				}
				//Envìa el comando para reiniciar la MV
				reboot, error := SendSSHCommand(host.Ip, rebootCommand, config)
				if error != nil {
					log.Println("Error al reinciar la MV:", reboot)
					return "Error al reinciar la MV"
				}
				fmt.Println("Reiniciando la màquina: " + nameVM)
				maxEspera = time.Now().Add(2 * time.Minute) //Agrega dos minutos de tiempo màximo para obtener la IP cuando se reincia la MV
				restarted = true
			}
		}

		//Almacena la direccion ip de la maquina virtual
		ipAddress = strings.TrimSpace(strings.TrimPrefix(ipAddress, "Value:"))

		//Actualiza el estado de la MV en la base de datos
		err9 := database.UpdateVirtualMachineState(nameVM, "Encendido")
		if err9 != nil {
			log.Println("Error al realizar la actualizaciòn del estado", err9)
			return "Error al realizar la actualizaciòn del estado"
		}
		//Actualiza la direcciòn IP de la MV en la base de datos
		err10 := database.UpdateVirtualMachineIP(nameVM, ipAddress)
		if err10 != nil {
			log.Println("Error al realizar la actualizaciòn de la IP", err10)
			return "Error al realizar la actualizaciòn de la IP"
		}
		fmt.Println("Màquina encendida, la direcciòn IP es: " + ipAddress)
		return ipAddress
	}
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

	if !running { //En caso de que la MV estè apagada, entonces se invoca el mètodo para encenderla
		StartVM(nameVM, clientIP)
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
