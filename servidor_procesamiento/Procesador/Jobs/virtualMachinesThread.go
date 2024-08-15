package jobs

import (
	"encoding/json"
	"fmt"
	config "servidor_procesamiento/Procesador/Config"
	models "servidor_procesamiento/Procesador/Models"
	utilities "servidor_procesamiento/Procesador/Utilities"
	"time"
)

// Funcion encargada de gestionar los cambios que ocurren la cola de creación y gestión de maquinas virtualesl.

func CheckVirtualMachinesQueueChanges() {
	for {
		// Verifica si el tamaño de la cola de especificaciones ha cambiado.
		mu := config.GetMu()
		mu.Lock()

		currentQueueSize := config.GetMaquina_virtualQueue().Queue.Len()
		//currentQueueSize :=  maquina_virtualesQueue.Queue.Len()
		mu.Unlock()

		if currentQueueSize > 0 {
			// Imprime y elimina el primer elemento de la cola de especificaciones.
			mu.Lock()
			firstElement := config.GetMaquina_virtualQueue().Queue.Front()
			data, dataPresent := firstElement.Value.(map[string]interface{})
			//maquina_virtualesQueue.Queue.Remove(firstElement)
			mu.Unlock()

			if !dataPresent {
				fmt.Println("No se pudo procesar la solicitud")
				mu.Lock()
				config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
				//maquina_virtualesQueue.Queue.Remove(firstElement)
				mu.Unlock()
				continue
			}

			specsMap, _ := data["specifications"].(map[string]interface{})
			specsJSON, err := json.Marshal(specsMap)
			if err != nil {
				fmt.Println("Error al extraer las especificaciones:", err)
				mu.Lock()
				config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
				//maquina_virtualesQueue.Queue.Remove(firstElement)
				mu.Unlock()
				continue
			}

			var specifications models.Maquina_virtual
			err = json.Unmarshal(specsJSON, &specifications)
			if err != nil {
				fmt.Println("Error al deserializar las especificaciones:", err)
				mu.Lock()
				config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
				//maquina_virtualesQueue.Queue.Remove(firstElement)
				mu.Unlock()
				continue
			}

			clientIP := data["clientIP"].(string)

			fmt.Println(clientIP)

			go utilities.CreateVM(specifications, clientIP)
			config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
			//maquina_virtualesQueue.Queue.Remove(firstElement)
			utilities.PrintVirtualMachine(specifications, true)
		}

		// Espera un segundo antes de verificar nuevamente.
		time.Sleep(1 * time.Second)
	}
}
