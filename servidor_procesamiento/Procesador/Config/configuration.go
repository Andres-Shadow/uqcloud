package config

import (
	"container/list"
	"fmt"
	"net/http"
	models "servidor_procesamiento/Procesador/Models"
	"sync"
	"time"
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

// Estructura que contiene los hosts y el host actual
// para la asignación mediante round robin
type RoundRobin struct {
	Hosts       []models.Host
	CurrentHost int
}

var RoundRobinManager *RoundRobin = nil

// Constructor para la estructura RoundRobin
func NewRoundRobin(hosts []models.Host) *RoundRobin {
	return &RoundRobin{
		Hosts:       hosts,
		CurrentHost: -1, // Inicializamos en -1 para que la primera vez sea 0
	}
}

// // Función para obtener el siguiente host
func (rr *RoundRobin) GetNextHost() models.Host {
	// Avanzamos al siguiente host en la lista
	rr.CurrentHost = (rr.CurrentHost + 1) % len(rr.Hosts)
	return rr.Hosts[rr.CurrentHost]
}

// Funcion para actualizar la lista de host en la estructura RoundRobin
func (rr *RoundRobin) UpdateHosts(hosts []models.Host) {
	rr.Hosts = hosts
}

// Funcion que inicializa la ruta del archivo de clave privada
// @param path: string
func InitPrivateKeyPath(path string) {
	if path != "" {
		privateKeyPath = path
	} else {
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

// Función para recargar la configuración de Prometheus al iniciar el servidor
func ReloadPrometheusConfig() {
	updateConfigPrometheusURL := "http://prometheus:9090/-/reload"

	// Realiza la solicitud POST
	req, err := http.NewRequest("POST", updateConfigPrometheusURL, nil)
	if err != nil {
		fmt.Println("Error creando la solicitud:", err)
		return
	}

	go func() {
		// Esperar 4 segundos antes de enviar la solicitud
		time.Sleep(4 * time.Second)

		// Envía la solicitud
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error haciendo la solicitud:", err)
			return
		}
		defer resp.Body.Close()

		// Verifica el estado de la respuesta
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: Prometheus devolvió el código de estado %d\n", resp.StatusCode)
			return
		}

		fmt.Println("Configuración de Prometheus recargada exitosamente", resp)
	}()
}
