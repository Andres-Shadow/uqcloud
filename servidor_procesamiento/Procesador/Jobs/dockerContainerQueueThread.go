package jobs

import (
	"fmt"
	config "servidor_procesamiento/Procesador/Config"
	utilities "servidor_procesamiento/Procesador/Utilities"
	"strings"
	"time"
)

// Funciòn que se encarga de gestionar la cola de solicitudes para la gestiòn de Contenedores Docker

func CheckContainerQueueChanges() {
	for {
		config.GetMu().Lock()
		currentQueueSize := config.GetDocker_contenedorQueue().Queue.Len()
		config.GetMu().Unlock()

		if currentQueueSize > 0 {
			config.GetMu().Lock()
			firstElement := config.GetDocker_contenedorQueue().Queue.Front()
			data, dataPresent := firstElement.Value.(map[string]interface{})
			config.GetMu().Unlock()

			if !dataPresent {
				fmt.Println("No se pudo procesar la solicitud")
				config.GetMu().Lock()
				config.GetDocker_contenedorQueue().Queue.Remove(firstElement)
				config.GetMu().Unlock()
				continue
			}

			tipoSolicitud, _ := data["solicitud"].(string)

			switch strings.ToLower(tipoSolicitud) {
			case "correr":

				contenedor := data["contenedor"].(string)
				ip := data["ip"].(string)
				hostname := data["hostname"].(string)

				go utilities.RunContainer(contenedor, ip, hostname)

			case "pausar":

				contenedor := data["contenedor"].(string)
				ip := data["ip"].(string)
				hostname := data["hostname"].(string)
				go utilities.StopContainer(contenedor, ip, hostname)

			case "reiniciar":

				contenedor := data["contenedor"].(string)
				ip := data["ip"].(string)
				hostname := data["hostname"].(string)
				go utilities.RestartContainer(contenedor, ip, hostname)

			case "borrar":

				contenedor := data["contenedor"].(string)
				ip := data["ip"].(string)
				hostname := data["hostname"].(string)
				go utilities.DeleteContainer(contenedor, ip, hostname)

			case "eliminar":

				ip := data["ip"].(string)
				hostname := data["hostname"].(string)
				go utilities.DeleteAllContainers(ip, hostname)

			default:
				fmt.Println("Tipo de solicitud no válido:", tipoSolicitud)
			}

			config.GetMu().Lock()
			config.GetDocker_contenedorQueue().Queue.Remove(firstElement)
			config.GetMu().Unlock()
		}

		time.Sleep(1 * time.Second) //Espera 1 segundo para volver a verificar la cola
	}
}
