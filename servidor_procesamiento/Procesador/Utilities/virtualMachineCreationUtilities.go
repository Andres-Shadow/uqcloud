package utilities

import (
	"database/sql"
	"fmt"
	"log"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
	"strconv"
	"strings"
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

			caracteres := GenerateRandomString(4) //Genera 4 caracteres alfanumèricos para concatenarlos al nombre de la MV

			nameVM := specs.Nombre + "_" + caracteres

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
			//Configura la conexiòn SSH con el host Obtenido
			config, err := ConfigureSSH(host.Hostname, config.GetPrivateKeyPath())
			if err != nil {
				log.Println("Error al configurar SSH:", err)
				return "Error al configurar la conexiòn SSH"
			}

			//Comando para crear una màquina virtual
			createVM := "VBoxManage createvm --name " + "\"" + nameVM + "\"" + " --ostype " + disco.Distribucion_sistema_operativo + "_" + strconv.Itoa(disco.Arquitectura) + " --register"
			uuid, err1 := SendSSHCommand(host.Ip, createVM, config)
			if err1 != nil {
				log.Println("Error al ejecutar el comando para crear y registrar la MV:", err1)
				return "Error al crear la MV"
			}

			//Comando para asignar la memoria RAM a la MV
			memoryCommand := "VBoxManage modifyvm " + "\"" + nameVM + "\"" + " --memory " + strconv.Itoa(specs.Ram)
			_, err2 := SendSSHCommand(host.Ip, memoryCommand, config)
			if err2 != nil {
				log.Println("Error ejecutar el comando para asignar la memoria a la MV:", err2)
				return "Error al asignar la memoria a la MV"
			}

			//Comando para agregar el controlador de almacenamiento
			sctlCommand := "VBoxManage storagectl " + "\"" + nameVM + "\"" + " --name hardisk --add sata"
			_, err3 := SendSSHCommand(host.Ip, sctlCommand, config)
			if err3 != nil {
				log.Println("Error al ejecutar el comando para asignar el controlador de almacenamiento a la MV:", err3)
				return "Error al asignar el controlador de almacenamiento a la MV"
			}

			//Comando para conectar el disco multiconexiòn a la MV
			sattachCommand := "VBoxManage storageattach " + "\"" + nameVM + "\"" + " --storagectl hardisk --port 0 --device 0 --type hdd --medium " + "\"" + disco.Ruta_ubicacion + "\""
			_, err4 := SendSSHCommand(host.Ip, sattachCommand, config)
			if err4 != nil {
				log.Println("Error al ejecutar el comando para conectar el disco a la MV: ", err4)
				return "Error al conectar el disco a la MV"
			}

			//Comando para asignar las unidades de procesamiento
			cpuCommand := "VBoxManage modifyvm " + "\"" + nameVM + "\"" + " --cpus " + strconv.Itoa(specs.Cpu)
			_, err5 := SendSSHCommand(host.Ip, cpuCommand, config)
			if err5 != nil {
				log.Println("Error al ejecutar el comando para asignar la cpu a la MV:", err5)
				return "Error al asignar la cpu a la MV"
			}

			//Comando para poner el adaptador de red en modo puente (Bridge)
			redAdapterCommand := "VBoxManage modifyvm " + "\"" + nameVM + "\"" + " --nic1 bridged --bridgeadapter1 " + "\"" + host.Adaptador_red + "\""
			_, err6 := SendSSHCommand(host.Ip, redAdapterCommand, config)
			if err6 != nil {
				log.Println("Error al ejecutar el comando para configurar el adaptador de red de la MV:", err6)
				return "Error al configurar el adaptador de red de la MV"
			}

			//Obtiene el UUID de la màquina virtual creda
			lines := strings.Split(string(uuid), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "UUID:") {
					uuid = strings.TrimPrefix(line, "UUID:")
				}
			}
			currentTime := time.Now().UTC()

			nuevaMaquinaVirtual := models.Maquina_virtual{
				Uuid:              uuid,
				Nombre:            nameVM,
				Sistema_operativo: specs.Sistema_operativo,
				Ram:               specs.Ram,
				Cpu:               specs.Cpu,
				Estado:            "Apagado",
				Hostname:          "uqcloud",
				Persona_email:     specs.Persona_email,
				Fecha_creacion:    currentTime,
			}

			//Crea el registro de la nueva MV en la base de datos
			_, err7 := database.DB.Exec("INSERT INTO maquina_virtual (uuid, nombre,  ram, cpu, ip, estado, hostname, persona_email, host_id, disco_id, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				nuevaMaquinaVirtual.Uuid, nuevaMaquinaVirtual.Nombre, nuevaMaquinaVirtual.Ram, nuevaMaquinaVirtual.Cpu,
				nuevaMaquinaVirtual.Ip, nuevaMaquinaVirtual.Estado, nuevaMaquinaVirtual.Hostname, nuevaMaquinaVirtual.Persona_email,
				host.Id, disco.Id, nuevaMaquinaVirtual.Fecha_creacion)
			if err7 != nil {
				log.Println("Error al crear el registro en la base de datos:", err7)
				return "Error al crear el registro en la base de datos"
			}

			//Calcula la CPU y RAM usada en el host
			usedCpu := host.Cpu_usada + specs.Cpu
			usedRam := host.Ram_usada + (specs.Ram)

			//Actualiza la informaciòn de los recursos usados en el host
			_, err8 := database.DB.Exec("UPDATE host SET ram_usada = ?, cpu_usada = ? where id = ?", usedRam, usedCpu, host.Id)
			if err8 != nil {
				log.Println("Error al actualizar el host en la base de datos: ", err8)
				return "Error al actualizar el host en la base de datos"
			}

			fmt.Println("Màquina virtual creada con èxito")
			StartVM(nameVM, clientIP)
			return "Màquina virtual creada con èxito"

		}

	} else {
		//Creacion de la Maquina con Algoritmo aleatorio
		//Obtiene el usuario
		user, error0 := database.GetUser(specs.Persona_email)
		if error0 != nil {
			log.Println("Error al obtener el usuario")
			return ""
		}

		if user.Rol == "Estudiante" {
			/*if cantidad >= 5 {
				fmt.Println("El usuario " + user.Nombre + " no puede crear màs de 5 màquinas virtuales.")
				return "El usuario " + user.Nombre + " no puede crear màs de 5 màquinas virtuales."
			}*/
		}

		caracteres := GenerateRandomString(4) //Genera 4 caracteres alfanumèricos para concatenarlos al nombre de la MV

		nameVM := specs.Nombre + "_" + caracteres

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
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM host").Scan(&count)
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
		//Configura la conexiòn SSH con el host
		config, err := ConfigureSSH(host.Hostname, config.GetPrivateKeyPath())
		if err != nil {
			log.Println("Error al configurar SSH:", err)
			return "Error al configurar la conexiòn SSH"
		}

		//Comando para crear una màquina virtual
		createVM := "VBoxManage createvm --name " + "\"" + nameVM + "\"" + " --ostype " + disco.Distribucion_sistema_operativo + "_" + strconv.Itoa(disco.Arquitectura) + " --register"
		uuid, err1 := SendSSHCommand(host.Ip, createVM, config)
		if err1 != nil {
			log.Println("Error al ejecutar el comando para crear y registrar la MV:", err1)
			return "Error al crear la MV"
		}

		//Comando para asignar la memoria RAM a la MV
		memoryCommand := "VBoxManage modifyvm " + "\"" + nameVM + "\"" + " --memory " + strconv.Itoa(specs.Ram)
		_, err2 := SendSSHCommand(host.Ip, memoryCommand, config)
		if err2 != nil {
			log.Println("Error ejecutar el comando para asignar la memoria a la MV:", err2)
			return "Error al asignar la memoria a la MV"
		}

		//Comando para agregar el controlador de almacenamiento
		sctlCommand := "VBoxManage storagectl " + "\"" + nameVM + "\"" + " --name hardisk --add sata"
		_, err3 := SendSSHCommand(host.Ip, sctlCommand, config)
		if err3 != nil {
			log.Println("Error al ejecutar el comando para asignar el controlador de almacenamiento a la MV:", err3)
			return "Error al asignar el controlador de almacenamiento a la MV"
		}

		//Comando para conectar el disco multiconexiòn a la MV
		sattachCommand := "VBoxManage storageattach " + "\"" + nameVM + "\"" + " --storagectl hardisk --port 0 --device 0 --type hdd --medium " + "\"" + disco.Ruta_ubicacion + "\""
		_, err4 := SendSSHCommand(host.Ip, sattachCommand, config)
		if err4 != nil {
			log.Println("Error al ejecutar el comando para conectar el disco a la MV: ", err4)
			return "Error al conectar el disco a la MV"
		}

		//Comando para asignar las unidades de procesamiento
		cpuCommand := "VBoxManage modifyvm " + "\"" + nameVM + "\"" + " --cpus " + strconv.Itoa(specs.Cpu)
		_, err5 := SendSSHCommand(host.Ip, cpuCommand, config)
		if err5 != nil {
			log.Println("Error al ejecutar el comando para asignar la cpu a la MV:", err5)
			return "Error al asignar la cpu a la MV"
		}

		//Comando para poner el adaptador de red en modo puente (Bridge)
		redAdapterCommand := "VBoxManage modifyvm " + "\"" + nameVM + "\"" + " --nic1 bridged --bridgeadapter1 " + "\"" + host.Adaptador_red + "\""
		_, err6 := SendSSHCommand(host.Ip, redAdapterCommand, config)
		if err6 != nil {
			log.Println("Error al ejecutar el comando para configurar el adaptador de red de la MV:", err6)
			return "Error al configurar el adaptador de red de la MV"
		}

		//Obtiene el UUID de la màquina virtual creda
		lines := strings.Split(string(uuid), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "UUID:") {
				uuid = strings.TrimPrefix(line, "UUID:")
			}
		}
		currentTime := time.Now().UTC()

		nuevaMaquinaVirtual := models.Maquina_virtual{
			Uuid:              uuid,
			Nombre:            nameVM,
			Sistema_operativo: specs.Sistema_operativo,
			Ram:               specs.Ram,
			Cpu:               specs.Cpu,
			Estado:            "Apagado",
			Hostname:          "uqcloud",
			Persona_email:     specs.Persona_email,
			Fecha_creacion:    currentTime,
		}

		//Crea el registro de la nueva MV en la base de datos
		_, err7 := database.DB.Exec("INSERT INTO maquina_virtual (uuid, nombre,  ram, cpu, ip, estado, hostname, persona_email, host_id, disco_id, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			nuevaMaquinaVirtual.Uuid, nuevaMaquinaVirtual.Nombre, nuevaMaquinaVirtual.Ram, nuevaMaquinaVirtual.Cpu,
			nuevaMaquinaVirtual.Ip, nuevaMaquinaVirtual.Estado, nuevaMaquinaVirtual.Hostname, nuevaMaquinaVirtual.Persona_email,
			host.Id, disco.Id, nuevaMaquinaVirtual.Fecha_creacion)
		if err7 != nil {
			log.Println("Error al crear el registro en la base de datos:", err7)
			return "Error al crear el registro en la base de datos"
		}

		//Calcula la CPU y RAM usada en el host
		usedCpu := host.Cpu_usada + specs.Cpu
		usedRam := host.Ram_usada + (specs.Ram)

		//Actualiza la informaciòn de los recursos usados en el host
		_, err8 := database.DB.Exec("UPDATE host SET ram_usada = ?, cpu_usada = ? where id = ?", usedRam, usedCpu, host.Id)
		if err8 != nil {
			log.Println("Error al actualizar el host en la base de datos: ", err8)
			return "Error al actualizar el host en la base de datos"
		}

		fmt.Println("Màquina virtual creada con èxito")
		StartVM(nameVM, clientIP)
		return "Màquina virtual creada con èxito"

	}
	return "solicitud invalida"
}
