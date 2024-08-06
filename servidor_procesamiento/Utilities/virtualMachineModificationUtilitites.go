package utilities

import (
	"fmt"
	"log"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
	"strconv"
)

/* Funciòn que contiene los comandos necesarios para modificar una màquina virtual. Primero verifica
si la màquina esta encendida o apagada. En caso de que estè encendida, se le notifica el usuario que debe apagar la màquina.
Tambièn, verifica si el host tiene los recusos necesarios solicitados por el cliente, en caso de que requiera aumentar lso recursos.
@specs Paràmetro que contiene las especificaciones a modificar en la màquina virtual
*/

func ModifyVirtualMachine(specs models.Maquina_virtual) string {

	//Obtiene la màquina virtual a modificar
	maquinaVirtual, err1 := database.GetVM(specs.Nombre)
	if err1 != nil {
		log.Println("Error al obtener la MV:", err1)
		return "Error al obtener la MV"
	}

	//Obtiene el host en el cual està alojada la MV
	host, err2 := database.GetHost(maquinaVirtual.Host_id)
	if err2 != nil {
		log.Println("Error al obtener el host:", err2)
		return "Error al obtener el host"
	}

	//Configura la conexiòn SSH con el host
	config, err := ConfigurarSSH(host.Hostname, config.GetPrivateKeyPath())
	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar SSH"
	}

	//Comando para modificar la memoria RAM a la MV
	memoryCommand := "VBoxManage modifyvm " + "\"" + specs.Nombre + "\"" + " --memory " + strconv.Itoa(specs.Ram)

	//Comando para modificar las unidades de procesamiento
	cpuCommand := "VBoxManage modifyvm " + "\"" + specs.Nombre + "\"" + " --cpus " + strconv.Itoa(specs.Cpu)

	//Variable que contiene el estado de la MV (Encendida o apagada)
	running, err3 := IsRunning(specs.Nombre, host.Ip, config)
	if err3 != nil {
		log.Println("Error al obtener el estado de la MV:", err3)
		return "Error al obtener el estado de la MV"
	}

	if running {
		fmt.Println("Para modificar la màquina primero debe apagarla")
		return "Para modificar la màquina primero debe apagarla"
	}

	if specs.Cpu != 0 && specs.Cpu != maquinaVirtual.Cpu {
		var cpu_host_usada int
		flagCpu := true

		if specs.Cpu < maquinaVirtual.Cpu {
			cpu_host_usada = host.Cpu_usada - (maquinaVirtual.Cpu - specs.Cpu)
		} else {
			cpu_host_usada = host.Cpu_usada + (specs.Cpu - maquinaVirtual.Cpu)
			recDisponibles := ValidarDisponibilidadRecursosHost((specs.Cpu - maquinaVirtual.Cpu), 0, host) //Valida si el host tiene recursos disponibles
			if !recDisponibles {
				fmt.Println("No se pudo aumentar la cpu, no hay recursos disponibles en el host")
				flagCpu = false
			}
		}
		if flagCpu {
			//ACtualiza la CPU usada en el host
			_, er := database.DB.Exec("UPDATE host SET cpu_usada = ? where id = ?", cpu_host_usada, host.Id)
			if er != nil {
				log.Println("Error al actualizar la cpu_usada del host en la base de datos: ", er)
				return "Error al actualizar el host en la base de datos"
			}
			_, err11 := EnviarComandoSSH(host.Ip, cpuCommand, config)
			if err11 != nil {
				log.Println("Error al realizar la actualizaciòn de la cpu", err11)
				return "Error al realizar la actualizaciòn de la cpu"
			}
			//Actualiza la CPU que tiene la MV
			_, err1 := database.DB.Exec("UPDATE maquina_virtual set cpu = ? WHERE NOMBRE = ?", strconv.Itoa(specs.Cpu), specs.Nombre)
			if err1 != nil {
				log.Println("Error al realizar la actualizaciòn de la CPU", err1)
				return "Error al realizar la actualizaciòn de la CPU"
			}
			fmt.Println("Se modificò con èxito la CPU")
		}
	}

	if specs.Ram != 0 && specs.Ram != maquinaVirtual.Ram {
		var ram_host_usada int
		flagRam := true

		if specs.Ram < maquinaVirtual.Ram {
			ram_host_usada = host.Ram_usada - (maquinaVirtual.Ram - specs.Ram)
		} else {
			ram_host_usada = host.Ram_usada + (specs.Ram - maquinaVirtual.Ram)
			recDisponibles := ValidarDisponibilidadRecursosHost(0, (specs.Ram - maquinaVirtual.Ram), host) //Valida si el host tiene RAM disponible para realizar el aumento de recursos
			if !recDisponibles {
				fmt.Println("No se modificò la ram porque el host no tiene recursos disponibles")
				flagRam = false
			}
		}
		if flagRam {
			//Actualiza la RAM usada en el host
			_, er := database.DB.Exec("UPDATE host SET ram_usada = ? where id = ?", ram_host_usada, host.Id)
			if er != nil {
				log.Println("Error al actualizar la ram_usada del host en la base de datos: ", er)
				return "Error al actualizar el host en la base de datos"
			}
			_, err22 := EnviarComandoSSH(host.Ip, memoryCommand, config)
			if err22 != nil {
				log.Println("Error al realizar la actualizaciòn de la memoria", err22)
				return "Error al realizar la actualizaciòn de la memoria"
			}
			//Actualiza la RAM de la MV
			_, err2 := database.DB.Exec("UPDATE maquina_virtual set ram = ? WHERE NOMBRE = ?", strconv.Itoa(specs.Ram), specs.Nombre)
			if err2 != nil {
				log.Println("Error al realizar la actualizaciòn de la memoria en la base de datos", err2)
				return "Error al realizar la actualizaciòn de la memoria en la base de datos"
			}
			fmt.Println("Se modificò con èxito la RAM")
		}
	}
	return "Modificaciones realizadas con èxito"
}
