package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	config "servidor_procesamiento/Procesador/Config"
	models "servidor_procesamiento/Procesador/Models/Entities"
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
		mu.Unlock()

		if currentQueueSize > 0 {
			// Imprime y elimina el primer elemento de la cola de especificaciones.
			mu.Lock()
			firstElement := config.GetMaquina_virtualQueue().Queue.Front()
			if firstElement == nil {
				log.Println("Error: No se pudo obtener el primer elemento de la cola")
				mu.Unlock()
				continue
			}

			data, dataPresent := firstElement.Value.(map[string]interface{})
			mu.Unlock()

			if !dataPresent {
				log.Println("No se pudo procesar la solicitud: el primer elemento no contiene datos válidos")
				mu.Lock()
				config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
				mu.Unlock()
				continue
			}

			specsMap, _ := data["specifications"].(map[string]interface{})
			specsJSON, err := json.Marshal(specsMap)
			if err != nil {
				log.Println("Error al extraer las especificaciones:", err)
				mu.Lock()
				config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
				mu.Unlock()
				continue
			}

			var specifications models.Maquina_virtual
			err = json.Unmarshal(specsJSON, &specifications)
			if err != nil {
				log.Println("Error al deserializar las especificaciones:", err)
				mu.Lock()
				config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
				mu.Unlock()
				continue
			}

			clientIP, ok := data["clientIP"].(string)
			if !ok {
				log.Println("Error: La IP del cliente no está presente o no es de tipo string")
				mu.Lock()
				config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
				mu.Unlock()
				continue
			}

			fmt.Println(clientIP)

			go func(specifications models.Maquina_virtual, clientIP string) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recuperado de un pánico en CreateVM: %v", r)
					}
				}()
				utilities.CreateVM(specifications, clientIP)
			}(specifications, clientIP)

			mu.Lock()
			config.GetMaquina_virtualQueue().Queue.Remove(firstElement)
			mu.Unlock()

			utilities.PrintVirtualMachine(specifications, true)
		}

		// Espera un segundo antes de verificar nuevamente.
		time.Sleep(1 * time.Second)
	}
}
