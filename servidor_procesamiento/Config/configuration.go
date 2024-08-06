package config

import (
	"container/list"
	"fmt"
	models "servidor_procesamiento/Procesador/Models"
	"sync"
)

/*
Clase encargada de contener las variables globales que son utilizadas como elementos en las diferentes
clases del programa
*/

// Variable de entorno que contiene la ruta del archivo de clave privada
var privateKeyPath string


// Declaraciòn de variables globales
var (
	maquina_virtualesQueue models.Maquina_virtualQueue
	managementQueue        models.ManagementQueue
	docker_imagesQueue     models.Docker_imagesQueue
	docker_contenedorQueue models.Docker_contenedorQueue
	mu                     sync.Mutex
	LastQueueSize          int
)

// Funcion que inicializa la ruta del archivo de clave privada
// @param path: string
func InitPrivateKeyPath(path string){
	if path != "" {
		privateKeyPath = path
	}else{
		fmt.Println("No se ha especificado la ruta del archivo de clave privada")
	}
}

// Funcion que retorna la ruta del archivo de clave privada
// @return string
func GetPrivateKeyPath() string {
	return privateKeyPath
}

// Funcion encargada de inicializar los objetos de las colas usadas a lo largo del programa
func InitQueues() {
	maquina_virtualesQueue.Queue = list.New()
	managementQueue.Queue = list.New()
	docker_imagesQueue.Queue = list.New()
	docker_contenedorQueue.Queue = list.New()
}

// Funciones que retornan las variable globales
func GetMaquina_virtualQueue() *models.Maquina_virtualQueue {
	return &maquina_virtualesQueue
}

func GetManagementQueue() *models.ManagementQueue {
	return &managementQueue
}	

func GetDocker_imagesQueue() *models.Docker_imagesQueue {
	return &docker_imagesQueue
}

func GetDocker_contenedorQueue() *models.Docker_contenedorQueue {
	return &docker_contenedorQueue
}

func GetMu() *sync.Mutex {
	return &mu
}

func GetLastQueueSize() int {
	return LastQueueSize
}