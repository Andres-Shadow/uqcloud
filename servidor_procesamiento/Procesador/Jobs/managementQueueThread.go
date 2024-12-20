package jobs

import (
	"fmt"
	config "servidor_procesamiento/Procesador/Config"
	virtualmachineutilities "servidor_procesamiento/Procesador/Utilities/VirtualMachineUtilities"
	"strings"
	"time"
)

//Funciòn que se encarga de gestionar la cola de solicitudes para la gestiòn de màquinas virtuales

func CheckManagementQueueChanges() {
	for {
		config.GetMu().Lock()
		currentQueueSize := config.GetManagementQueue().Queue.Len()
		config.GetMu().Unlock()

		if currentQueueSize > 0 {
			config.GetMu().Lock()
			firstElement := config.GetManagementQueue().Queue.Front()
			data, dataPresent := firstElement.Value.(map[string]interface{})
			config.GetMu().Unlock()

			if !dataPresent {
				fmt.Println("No se pudo procesar la solicitud")
				config.GetMu().Lock()
				config.GetManagementQueue().Queue.Remove(firstElement)
				config.GetMu().Unlock()
				continue
			}

			tipoSolicitud, _ := data["tipo_solicitud"].(string)

			switch strings.ToLower(tipoSolicitud) {
			case "delete":
				nameVM, _ := data["nombreVM"].(string)
				go virtualmachineutilities.DeleteVM(nameVM)

			case "start":
				nameVM, _ := data["nombreVM"].(string)
				clientIP, _ := data["clientIP"].(string)
				go virtualmachineutilities.StartVM(nameVM, clientIP)

			case "stop":
				nameVM, _ := data["nombreVM"].(string)
				clientIP, _ := data["clientIP"].(string)
				go virtualmachineutilities.TurnOffVM(nameVM, clientIP)

			default:
				fmt.Println("Tipo de solicitud no válido:", tipoSolicitud)
			}

			config.GetMu().Lock()
			config.GetManagementQueue().Queue.Remove(firstElement)
			config.GetMu().Unlock()
		}

		time.Sleep(1 * time.Second) //Espera 1 segundo para volver a verificar la cola
	}
}
