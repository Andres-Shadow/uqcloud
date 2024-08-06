package jobs

import (
	"fmt"
	config "servidor_procesamiento/Procesador/Config"
	utilities "servidor_procesamiento/Procesador/Utilities"
	"strings"
	"time"
)

/* Funciòn que se encarga de gestionar la cola de solicitudes para la gestiòn de Imagenes Docker  */

func CheckImagesQueueChanges() {
	for {
		config.GetMu().Lock()
		currentQueueSize := config.GetDocker_imagesQueue().Queue.Len()
		config.GetMu().Unlock()

		if currentQueueSize > 0 {
			config.GetMu().Lock()
			firstElement := config.GetDocker_imagesQueue().Queue.Front()
			data, dataPresent := firstElement.Value.(map[string]interface{})
			config.GetMu().Unlock()

			if !dataPresent {
				fmt.Println("No se pudo procesar la solicitud")
				config.GetMu().Lock()
				config.GetDocker_imagesQueue().Queue.Remove(firstElement)
				config.GetMu().Unlock()
				continue
			}

			tipoSolicitud, _ := data["solicitud"].(string)

			fmt.Println(tipoSolicitud)

			switch strings.ToLower(tipoSolicitud) {

			case "borar":
				fmt.Println("Borrar")
				imagen := data["imagen"].(string)
				ip := data["ip"].(string)
				hostname := data["hostname"].(string)

				go utilities.EliminarImagen(imagen, ip, hostname)

			case "eliminar":
				ip := data["ip"].(string)
				hostname := data["hostname"].(string)
				go utilities.EliminarTodasImagenes(ip, hostname)

			default:
				fmt.Println("Tipo de solicitud no válido:", tipoSolicitud)
			}

			config.GetMu().Lock()
			config.GetDocker_imagesQueue().Queue.Remove(firstElement)
			config.GetMu().Unlock()
		}

		time.Sleep(1 * time.Second) //Espera 1 segundo para volver a verificar la cola
	}
}