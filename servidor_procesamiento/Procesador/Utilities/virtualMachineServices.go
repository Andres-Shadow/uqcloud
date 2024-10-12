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
Clase encargada de contener las funalidades de encedido y apagado de maquinas virtuales
*/

/*
Funciòn que permite iniciar una màquina virtual en modo "headless", lo que indica que se inicia en segundo plano
para que el usuario de la màquina fìsica no se vea afectado, realiza las comprobaciones iniciales y posteriormente
llama a la funciòn "startVMSequence" para realizar el encendido de la màquina virtual
*/

func StartVM(nameVM string, clientIP string) string {

	var err error
	var host models.Host
	var running bool
	var conf *ssh.ClientConfig
	var maquinaVirtual models.Maquina_virtual

	//Obtiene el objeto "maquina_virtual"
	maquinaVirtual, err = database.GetVM(nameVM)
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
		log.Println("Error al obtener el estado de la MV:", err)
		return "Error al obtener el estado de la MV" + err.Error()
	} else if running {
		return "La máquina ya está encendida"
	}

	fmt.Println("Encendiendo la màquina " + nameVM + "...")

	return startVMSequence(nameVM, clientIP, host, conf)

}

/* Funciòn que permite enviar el comando PowerOff para apagar una màquina virtual, realiza las comprobaciones iniciales
y posteriormente llama a la funciòn "powerOffVMSequence" para realizar el apagado de la màquina virtual
*/

func TurnOffVM(nameVM string, clientIP string) string {
	var err error
	var maquinaVirtual models.Maquina_virtual
	var host models.Host
	var conf *ssh.ClientConfig
	var running bool
	//Obtiene la màquina vitual a apagar
	maquinaVirtual, err = database.GetVM(nameVM)
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

	running, err = IsRunning(nameVM, host.Ip, conf)
	if err != nil {
		log.Println("Error al obtener el estado de la MV:", err)
		return "Error al obtener el estado de la MV"
	}

	if !running {
		log.Println("La máquina ya está apagada")
		return "La máquina ya está apagada"
	}

	return powerOffVMSequence(nameVM, host, conf)
}

// Funcion que se encarga de ejecutar la secuencia de encendido de la màquina virtual
func startVMSequence(nameVM, clientIP string, host models.Host, conf *ssh.ClientConfig) string {

	var err error
	var ipAddress string

	startVMCommand := getStartVMCommand(clientIP, nameVM)

	// verifica el envio del comando ssh
	if _, err = SendSSHCommand(host.Ip, startVMCommand, conf); err != nil {
		log.Println("Error al enviar el comando para encender la MV:", err)
		return "Error al enviar el comando para encender la MV: " + err.Error()
	}

	//Actualiza el estado de la màquina virtual
	if err = database.UpdateVirtualMachineState(nameVM, "Procesando"); err != nil {
		log.Println("Error al actualizar el estado de la MV:", err)
		return "Error al actualizar el estado de la MV: " + err.Error()
	}

	//Obtiene la direcciòn IP de la màquina virtual
	ipAddress, err = getVMIPAddress(nameVM, host, conf)
	if err != nil {
		return err.Error()
	}

	//Actualiza el estado de la màquina virtual y su direcciòn IP en la base de datos
	if err = finalizeVMStart(nameVM, ipAddress); err != nil {
		return err.Error()
	}

	return ipAddress
}

// funcion que retorna el comando para encender la màquina virtual según las condiciones deseadas
func getStartVMCommand(clientIP, nameVM string) string {
	if _, err := IsAHostIp(clientIP); err == nil {
		return "VBoxManage startvm " + "\"" + nameVM + "\""
	}
	return "VBoxManage startvm " + "\"" + nameVM + "\"" + " --type headless"
}

// funcion muy importante encargada de ejecutar el proceso de obtenciòn de la direcciòn IP de la màquina virtual
// mediante la repeticion de comandos y la espera de la respuesta correcta
func getVMIPAddress(nameVM string, host models.Host, conf *ssh.ClientConfig) (string, error) {

	// comandos necesarios para obtener la direccion IP de la màquina virtual mediante
	// las guestaditions de virtualbox
	var getIpCommand string = "VBoxManage guestproperty get " + "\"" + nameVM + "\"" + " /VirtualBox/GuestInfo/Net/0/V4/IP"
	var rebootCommand string = "VBoxManage controlvm " + "\"" + nameVM + "\"" + " reset"

	// variables necesarias para el control del proceso
	// dentro de la funcion
	var restarted bool = false
	var ipAddress string
	maxWait := time.Now().Add(2 * time.Minute)

	for time.Now().Before(maxWait) {
		fmt.Println("Obteniendo direcciòn IP de la màquina " + nameVM + "...")

		// envia el comando ssh para obtener la direcciòn IP o reiniciar si es necesario
		ipAddress, _ = SendSSHCommand(host.Ip, getIpCommand, conf)

		// limpia el output de virtual box para obtener la direcciòn IP
		// VirtualBox -> "Value: 192.168.x.x"
		//Limpieza-> "192.168.x.x"
		ipAddress = strings.TrimSpace(strings.TrimPrefix(ipAddress, "Value:"))

		// verifica si la direcciòn IP es correcta
		if ipAddress != "" && !strings.HasPrefix(ipAddress, "169") && ipAddress != "No value set!" {
			return ipAddress, nil
		}

		// verifica si se ha superado el tiempo de espera (variable maxWait)
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

// funcion encargada de reiniar la màquina virtual
func rebootVM(hostIP, rebootCommand string, conf *ssh.ClientConfig) error {
	_, err := SendSSHCommand(hostIP, rebootCommand, conf)
	if err != nil {
		log.Println("Error al reiniciar la MV:", err)
	}
	return err
}

// funcion encargada de finalizar el proceso de encendido de la màquina virtual
// actualizando su estado y su direcciòn IP en la base de datos
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

// funcion encargada de ejecutar el proceso de apagado de la màquina virtual
// siguiendo los pasos necesarios para garantizar el apagado correcto
func powerOffVMSequence(nameVM string, host models.Host, config *ssh.ClientConfig) string {

	// variables necesarias para el control del proceso
	var err error

	// comando necesario para apgar la màquina virtual
	var powerOffCommand string = "VBoxManage controlvm \"" + nameVM + "\" poweroff"

	fmt.Println("Apagando máquina " + nameVM + "...")

	// actualiza el estado de la màquina virtual en una primera instancia
	err = database.UpdateVirtualMachineState(nameVM, "Procesando")
	if err != nil {
		log.Println("Error al actualizar el estado de la MV:", err)
		return "Error al actualizar el estado de la MV"
	}

	// envia el comando ssh para apagar la màquina virtual
	_, err = SendSSHCommand(host.Ip, powerOffCommand, config)
	if err != nil {
		log.Println("Error al enviar el comando para apagar la MV:", err)
		return "Error al enviar el comando para apagar la MV" + err.Error()
	}

	// verifica si la màquina virtual se ha apagado correctamente
	if !waitForShutdown(nameVM, host, config) {
		if err := forceShutdown(nameVM, host, config); err != nil {
			return err.Error()
		}
	}

	// actualiza el estado de la màquina virtual en la base de datos a apagado
	// de forma definitiva
	err = database.UpdateVirtualMachineState(nameVM, "Apagado")
	if err != nil {
		log.Println("Error al actualizar el estado de la MV:", err)
		return "Error al actualizar el estado de la MV"
	}

	fmt.Println("Máquina apagada con éxito")
	return ""
}

// funcion encargada de verificar si la màquina virtual se ha apagado correctamente
// esperando x tiempo (variable maxWait) y verificando el estado de la màquina virtual
func waitForShutdown(nameVM string, host models.Host, config *ssh.ClientConfig) bool {
	var running bool
	var err error
	maxWait := time.Now().Add(5 * time.Minute)
	for time.Now().Before(maxWait) {
		running, err = IsRunning(nameVM, host.Ip, config)
		if err != nil {
			log.Println("Error al obtener el estado de la MV:", err)
			return false
		}
		if !running {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

// funcion encargada de forzar el apagado de la màquina virtual en caso de que no se haya apagado correctamente
// siguiendo los pasos listado en la funcion TurnOffVMSquence
func forceShutdown(nameVM string, host models.Host, config *ssh.ClientConfig) error {
	var powerOffCommand string = "VBoxManage controlvm \"" + nameVM + "\" poweroff"
	var err error
	_, err = SendSSHCommand(host.Ip, powerOffCommand, config)
	if err != nil {
		log.Println("Error al enviar el comando para apagar la MV:", err)
		return fmt.Errorf("error al enviar el comando para apagar la MV")
	}
	return nil
}
