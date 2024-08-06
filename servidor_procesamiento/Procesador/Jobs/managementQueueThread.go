package jobs

import (
	"encoding/json"
	"fmt"
	config "servidor_procesamiento/Procesador/Config"
	models "servidor_procesamiento/Procesador/Models"
	utilities "servidor_procesamiento/Procesador/Utilities"
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
			case "modify":
				specsMap, _ := data["specifications"].(map[string]interface{})
				specsJSON, err := json.Marshal(specsMap)
				if err != nil {
					fmt.Println("Error al serializar las especificaciones:", err)
					config.GetMu().Lock()
					config.GetManagementQueue().Queue.Remove(firstElement)
					config.GetMu().Unlock()
					continue
				}

				var specifications models.Maquina_virtual
				err = json.Unmarshal(specsJSON, &specifications)
				if err != nil {
					fmt.Println("Error al deserializar las especificaciones:", err)
					config.GetMu().Lock()
					config.GetManagementQueue().Queue.Remove(firstElement)
					config.GetMu().Unlock()
					continue
				}

				go utilities.ModifyVirtualMachine(specifications)

			case "delete":
				nameVM, _ := data["nombreVM"].(string)
				go utilities.DeleteVM(nameVM)

			case "start":
				nameVM, _ := data["nombreVM"].(string)
				clientIP, _ := data["clientIP"].(string)
				go utilities.StartVM(nameVM, clientIP)

			case "stop":
				nameVM, _ := data["nombreVM"].(string)
				clientIP, _ := data["clientIP"].(string)
				go utilities.TurnOffVM(nameVM, clientIP)

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
