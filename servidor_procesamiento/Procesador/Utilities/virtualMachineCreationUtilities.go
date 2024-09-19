package utilities

import (
	"fmt"
	"log"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"

	"strconv"
	"time"
)

/*
	Esta funciòn permite enviar los comandos VBoxManage necesarios para crear una nueva màquina virtual
	Se encarga de verificar si el usuario aùn puede crear màquinas virtuales, dependiendo de su rol: màximo 5 para estudiantes y màximo 3 para invitados
	Se encarga de escoger el host dependiendo desde donde se hace la solicitud: algoritmo "here" si se realiza desde un computador que pertenece a los host ò aleatorio en caso contrario
	Valida si el host que se escogiò tiene recursos disponibles para crear la MV solicitada
	Finalmente, crea y actualiza los registros necesarios en la base de datos

@spects Paràmetro que contiene la configuraciòn enviada por el usuario para crear la MV
@clientIP Paràmetro que contiene la direcciòn IP de la màquina desde la cual se està realizando la peticiòn para crear la MV
*/
func CreateVM(specs models.Maquina_virtual, clientIP string) string {

	var disco models.Disco
	var err error
	var host models.Host
	caracteres := GenerateRandomString(4) //Genera 4 caracteres alfanumèricos para concatenarlos al nombre de la MV
	var virtualMachineName string = specs.Nombre + "_" + caracteres
	var availableResources bool = false
	var estadossh bool = true
	var validation bool = false

	// Verifica la disponibilidad del nombre de la MV
	dispoible, mensaje := verifyVirtualMachineExistence(virtualMachineName) //Verifica si existe una MV con ese nombre

	if !dispoible {
		return mensaje
	}

	if specs.Hostname == "aleatorio" {

		//Obtiene la cantidad total de hosts que hay en la base de datos
		count, err := database.CountRegisteredHosts()
		if err != nil {
			log.Println("Error al contar los host que hay en la base de datos: " + err.Error())
			return "Error al contar los host que hay en la base de datos"
		}

		// Bucle para seleccionar un host con recursos disponibles
		for !availableResources && count > 0 {
			aux, _ := database.SelectHost()
			estadossh = Pacemaker(config.GetPrivateKeyPath(), aux.Hostname, aux.Ip)

			// Si el estado de SSH es falso, rompe el bucle
			if !estadossh {
				log.Println("Conexión SSH fallida con el host:", aux.Hostname)
			} else {
				availableResources = ValidateHostResourceAvailability(specs.Cpu, specs.Ram, aux)
				if availableResources {
					host = aux
					break
				}
			}

			count--
		}

		if !availableResources {
			fmt.Println("No hay recursos disponibles el Desktop Cloud para crear la màquina virtual. Intente màs tarde")
			return "No hay recursos disponibles el Desktop Cloud para crear la màquina virtual. Intente màs tarde"
		}

		disco, err = database.GetDisk(specs.Sistema_operativo, specs.Distribucion_sistema_operativo, host.Id)
		if err != nil {
			log.Println("Error al obtener el disco:", err)
			return "Error al obtener el disco"
		}

	} else if specs.Hostname == "roundrobin" {

		// Bucle para seleccionar un host con recursos disponibles
		for !availableResources {

			fmt.Println("Obteniendo el host siguiente: " + config.RoundRobinManager.GetNextHost().Nombre)
			aux := config.RoundRobinManager.GetNextHost()
			estadossh = Pacemaker(config.GetPrivateKeyPath(), aux.Hostname, aux.Ip)

			// Si el estado de SSH es falso, rompe el bucle
			if !estadossh {
				log.Println("Conexión SSH fallida con el host:", aux.Hostname)
			} else {
				availableResources = ValidateHostResourceAvailability(specs.Cpu, specs.Ram, aux)
				if availableResources {
					host = aux
					break
				}
			}
		}

		if !availableResources {
			fmt.Println("No hay recursos disponibles el Desktop Cloud para crear la màquina virtual. Intente màs tarde")
			return "No hay recursos disponibles el Desktop Cloud para crear la màquina virtual. Intente màs tarde"
		}

		disco, err = database.GetDisk(specs.Sistema_operativo, specs.Distribucion_sistema_operativo, host.Id)

		if err != nil {
			log.Println("Error al obtener el disco:", err)
			return "Error al obtener el disco"
		}

	} else {

		host, err = GetHostByName(specs.Hostname)

		if err != nil {
			log.Println("Error al obtener el host por nombre:", err)
			return "Error al obtener el host por nombre " + specs.Hostname
		}

		log.Println("Host seleccionado: ", host.Nombre)

		//se verifica el ssh de la maquina fisica con el marcapasos
		estadossh := Pacemaker(config.GetPrivateKeyPath(), host.Hostname, host.Ip)
		if estadossh {
			//Obtenemos el disco multiconexion del host previamente obtenido
			disco, err = database.GetDisk(specs.Sistema_operativo, specs.Distribucion_sistema_operativo, host.Id)
			if err != nil {
				log.Println("Error al obtener el disco:", err)
				return "Error al obtener el disco"
			}

		} else {
			return "Error al verificar la conexión con el host seleccionado"
		}
	}

	validation = configureAndCreateVM(host, specs, virtualMachineName, disco)
	if !validation {
		return "Error al ejecutar los comandos para crear la máquina virtual mediante ssh"
	}

	err = createDatabaseRecords(host, specs, virtualMachineName, disco)
	if err != nil {
		return "Error al crear los registros de la máquina virtual en la base de datos"
	}

	// En caso de la solicitud sea desde el boton de "Lanzar Máquina Ahora"
	if specs.Nombre == "QuickGuest" {
		log.Println("Iniciando automáticamente la máquina virtual rápida en 3 segundos...")
		// Iniciar la máquina virtual después de 3 segundos de forma asincrona
		go StartVM(virtualMachineName, host.Hostname)
		return "Máquina virtual creada y en proceso de inicio"
	}

	return "solicitud invalida"
}

func verifyVirtualMachineExistence(virtualMachineName string) (bool, string) {
	//Consulta si existe una MV con ese nombre
	existe, error1 := ExistVM(virtualMachineName)
	if error1 != nil {
		log.Println("Error al consultar si existe una MV con el nombre indicado: ", error1)
		return false, "Error al consultar si existe una MV con el nombre indicado"

	} else if existe {
		fmt.Println("El nombre " + virtualMachineName + " no està disponible, por favor ingrese otro.")
		return false, "Nombre de la MV no disponible"
	}
	return true, ""
}

func configureAndCreateVM(host models.Host, specs models.Maquina_virtual, nameVM string, disco models.Disco) bool {
	// Configurar SSH y enviar comandos para crear y configurar la MV
	config, err := ConfigureSSH(host.Hostname, config.GetPrivateKeyPath())
	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return false
	}

	commands := []string{
		"VBoxManage createvm --name \"" + nameVM + "\" --ostype  \"" + disco.Distribucion_sistema_operativo + "_" + strconv.Itoa(disco.Arquitectura) + "\" --register",
		"VBoxManage modifyvm \"" + nameVM + "\" --memory " + strconv.Itoa(specs.Ram),
		"VBoxManage storagectl \"" + nameVM + "\" --name hardisk --add sata",
		"VBoxManage storageattach \"" + nameVM + "\" --storagectl hardisk --port 0 --device 0 --type hdd --medium \"" + disco.Ruta_ubicacion + "\"",
		"VBoxManage modifyvm \"" + nameVM + "\" --cpus " + strconv.Itoa(specs.Cpu),
		"VBoxManage modifyvm \"" + nameVM + "\" --nic1 bridged --bridgeadapter1 \"" + host.Adaptador_red + "\"",
	}

	for _, command := range commands {
		if _, err := SendSSHCommand(host.Ip, command, config); err != nil {
			log.Println("Error al ejecutar comando:", err)
			return false
		}
	}

	return true
}

func createDatabaseRecords(host models.Host, specs models.Maquina_virtual, nameVM string, disco models.Disco) error {
	// Lógica para crear los registros en la base de datos y actualizar recursos
	currentTime := time.Now().UTC()

	log.Println("Creando registros en la base de datos...")
	log.Println("Host seleccionado: ", host.Id)

	nuevaMaquinaVirtual := models.Maquina_virtual{
		Uuid:                           nameVM + "_uuid",
		Nombre:                         nameVM,
		Ram:                            specs.Ram,
		Cpu:                            specs.Cpu,
		Ip:                             "",
		Estado:                         "Apagado",
		Hostname:                       host.Hostname,
		Persona_email:                  specs.Persona_email,
		Host_id:                        host.Id,
		Disco_id:                       disco.Id,
		Sistema_operativo:              specs.Sistema_operativo,
		Distribucion_sistema_operativo: specs.Distribucion_sistema_operativo,
		Fecha_creacion:                 currentTime}

	if err := database.CreateVirtualMachine(nuevaMaquinaVirtual); err != nil {
		log.Println("Error al crear el registro en la base de datos:", err)
		return err
	}

	usedCpu := host.Cpu_usada + specs.Cpu
	usedRam := host.Ram_usada + specs.Ram
	if err := database.UpdateHostRamAndCPU(host.Id, usedRam, usedCpu); err != nil {
		log.Println("Error al actualizar el host en la base de datos:", err)
		return err
	}

	return nil
}
