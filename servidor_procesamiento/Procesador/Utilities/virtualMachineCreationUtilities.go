package utilities

import (
	"database/sql"
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

	caracteres := GenerateRandomString(4) //Genera 4 caracteres alfanumèricos para concatenarlos al nombre de la MV

	if specs.Host_id > 0 {
		// Creacion de Maquina Virtual con seleccion de usuario
		// Obtenemeos el host por medio del indice que es previamente
		mihost, err := database.GetHost(specs.Host_id)
		if err != nil {
			fmt.Println("Error al obtener el host", err)

		}

		//se verifica el ssh de la maquina fisica con el marcapasos
		estadossh := Pacemaker(config.GetPrivateKeyPath(), mihost.Hostname, mihost.Ip)
		if estadossh {

			nameVM := specs.Nombre + "_" + caracteres

			//Consulta si existe una MV con ese nombre
			existe, error1 := ExistVM(nameVM)
			if error1 != nil {
				fmt.Println("Erro al verificar si existe un registro en las maquinas virtuales con el nombre " + nameVM + " en la base de datos: " + error1.Error())
				return "Error a la hora de verificar si existe un registro en las maquinas virtuales con el nombre: " + nameVM
			} else if existe {
				fmt.Println("El nombre " + nameVM + " no està disponible, por favor ingrese otro.")
				return "Nombre de la MV no disponible"
			}

			//Inicializamos la creacion de una variable tipo host
			var host models.Host

			//Obtenemos el host
			host, _ = database.GetHost(specs.Host_id)

			//Obtenemos el disco multiconexion del host previamente obtenido
			disco, err20 := database.GetDisk(specs.Sistema_operativo, specs.Distribucion_sistema_operativo, host.Id)
			if err20 != nil {
				log.Println("Error al obtener el disco:", err20)
				return "Error al obtener el disco"
			}

			//Obtiene el UUID de la màquina virtual creda
			// lines := strings.Split(string(uuid), "\n")
			// for _, line := range lines {
			// 	if strings.HasPrefix(line, "UUID:") {
			// 		uuid = strings.TrimPrefix(line, "UUID:")
			// 	}
			// }

			validation := configureAndCreateVM(host, specs, nameVM, disco)

			if !validation {
				return "Error al ejecutar los comandos para crear la máquina virtual mediante ssh"
			}

			validation2 := createDatabaseRecords(host, specs, nameVM)
			if validation2 != nil {
				return "Error al crear los registros de la máquina virtual en la base de datos"
			}
		}

	} else {
		//Creacion de la Maquina con Algoritmo aleatorio
		//Obtiene el usuario

		

		nameVM := specs.Nombre + "_" + caracteres

		fmt.Println("llego hasta aqui")
		//Consulta si existe una MV con ese nombre
		existe, error1 := ExistVM(nameVM)
		if error1 != nil {
			if error1 != sql.ErrNoRows {
				log.Println("Error al consultar si existe una MV con el nombre indicado: ", error1)
				return "Error al consultar si existe una MV con el nombre indicado"
			}
		} else if existe {
			fmt.Println("El nombre " + nameVM + " no està disponible, por favor ingrese otro.")
			return "Nombre de la MV no disponible"
		}

		var host models.Host
		availableResources := false
		host, er := IsAHostIp(clientIP) //Consulta si la ip de la peticiòn proviene de un host registrado en la BD
		if er == nil {                  //nil = El host existe
			availableResources = ValidateHostResourceAvailability(specs.Cpu, specs.Ram, host) //Verifica si el host tiene recursos disponibles
		}
		fmt.Print("available", availableResources)

		//Obtiene la cantidad total de hosts que hay en la base de datos
		count, err := database.CountRegisteredHosts()
		if err != nil {
			log.Println("Error al contar los host que hay en la base de datos: " + err.Error())
			return "Error al contar los gost que hay en la base de datos"
		}
		count += 5 //Para dar n+5 iteraciones en busca de hosts con recursos disponibles, donde n es el total de hosts guardados en la bse de datos
		fmt.Print("count :", count)
		estadossh := Pacemaker(config.GetPrivateKeyPath(), host.Hostname, host.Ip)
		//Escoge hosts al azar en busca de alguno que tenga recursos disponibles para crear la MV
		log.Println(estadossh)
		for !estadossh && count > 0 {
			//Selecciona un host al azar

			host, _ = database.SelectHost()
			estadossh = Pacemaker(config.GetPrivateKeyPath(), host.Hostname, host.Ip)
			if err != nil {
				log.Println("Error al seleccionar el host:", err)

			}
			availableResources = ValidateHostResourceAvailability(specs.Cpu, specs.Ram, host) //Verifica si el host tiene recursos disponibles
			//fmt.Print("resources :", availableResources)
			count--
		}

		if !availableResources {
			fmt.Println("No hay recursos disponibles el Desktop Cloud para crear la màquina virtual. Intente màs tarde")
		}

		disco, err20 := database.GetDisk(specs.Sistema_operativo, specs.Distribucion_sistema_operativo, host.Id)
		if err20 != nil {
			log.Println("Error al obtener el disco:", err20)

			return "Error al obtener el disco"
		}

		validation := configureAndCreateVM(host, specs, nameVM, disco)
		if !validation {
			return "Error al ejecutar los comandos para crear la máquina virtual mediante ssh"
		}
		// //Obtiene el UUID de la màquina virtual creda
		// lines := strings.Split(string(uuid), "\n")
		// for _, line := range lines {
		// 	if strings.HasPrefix(line, "UUID:") {
		// 		uuid = strings.TrimPrefix(line, "UUID:")
		// 	}
		// }

		validation2 := createDatabaseRecords(host, specs, nameVM)
		if validation2 != nil {
			return "Error al crear los registros de la máquina virtual en la base de datos"
		}

	}
	return "solicitud invalida"
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

func createDatabaseRecords(host models.Host, specs models.Maquina_virtual, nameVM string) error {
	// Lógica para crear los registros en la base de datos y actualizar recursos
	currentTime := time.Now().UTC()

	nuevaMaquinaVirtual := models.Maquina_virtual{
		Uuid:              nameVM + "_uuid",
		Nombre:            nameVM,
		Sistema_operativo: specs.Sistema_operativo,
		Ram:               specs.Ram,
		Cpu:               specs.Cpu,
		Estado:            "Apagado",
		Hostname:          "uqcloud",
		Persona_email:     specs.Persona_email,
		Fecha_creacion:    currentTime,
	}

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
