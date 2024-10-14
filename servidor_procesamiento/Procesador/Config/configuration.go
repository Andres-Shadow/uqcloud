package config

import (
	"container/list"
	"fmt"
	"net/http"
	models "servidor_procesamiento/Procesador/Models/Entities"
	internal "servidor_procesamiento/Procesador/Models/Internal"

	"sync"
	"time"
)

/*
Configuración global del servidor de procesamiento.
Esta clase contiene variables globales y funciones que se utilizan en diferentes partes del programa.
*/

// Variables de entorno
var privateKeyPath string // Ruta del archivo de clave privada

// Variables globales del sistema de colas
var (
	maquina_virtualesQueue internal.Maquina_virtualQueue
	managementQueue        internal.ManagementQueue
	docker_imagesQueue     internal.Docker_imagesQueue
	docker_contenedorQueue internal.Docker_contenedorQueue
	mu                     sync.Mutex
	LastQueueSize          int
)

// RoundRobin: Gestión de hosts para el balanceo de carga
type RoundRobin struct {
	Hosts       []models.Host
	CurrentHost int
}

var RoundRobinManager *RoundRobin = nil

// ====== Funciones relacionadas con RoundRobin ======

// NewRoundRobin: Constructor para inicializar la estructura de RoundRobin
func NewRoundRobin(hosts []models.Host) *RoundRobin {

	return &RoundRobin{
		Hosts:       hosts,
		CurrentHost: -1, // Iniciamos en -1 para que la primera asignación sea el host 0
	}
}

// GetNextHost: Devuelve el siguiente host en el ciclo Round Robin
func (rr *RoundRobin) GetNextHost() models.Host {
	rr.CurrentHost = (rr.CurrentHost + 1) % len(rr.Hosts)
	return rr.Hosts[rr.CurrentHost]
}

// UpdateHosts: Actualiza la lista de hosts en la estructura RoundRobin
func (rr *RoundRobin) UpdateHosts(hosts []models.Host) {
	rr.Hosts = hosts
}

// ====== Funciones de inicialización y acceso a variables ======

// InitPrivateKeyPath: Inicializa la ruta del archivo de clave privada
func InitPrivateKeyPath(path string) {
	if path != "" {
		privateKeyPath = path
	} else {
		fmt.Println("No se ha especificado la ruta del archivo de clave privada")
	}
}

// GetPrivateKeyPath: Retorna la ruta actual del archivo de clave privada
func GetPrivateKeyPath() string {
	return privateKeyPath
}

// InitQueues: Inicializa las colas usadas en el programa
func InitQueues() {
	maquina_virtualesQueue.Queue = list.New()
	managementQueue.Queue = list.New()
	docker_imagesQueue.Queue = list.New()
	docker_contenedorQueue.Queue = list.New()
}

// ====== Getters para las variables globales ======

func GetMaquina_virtualQueue() *internal.Maquina_virtualQueue {
	return &maquina_virtualesQueue
}

func GetManagementQueue() *internal.ManagementQueue {
	return &managementQueue
}

func GetDocker_imagesQueue() *internal.Docker_imagesQueue {
	return &docker_imagesQueue
}

func GetDocker_contenedorQueue() *internal.Docker_contenedorQueue {
	return &docker_contenedorQueue
}

func GetMu() *sync.Mutex {
	return &mu
}

func GetLastQueueSize() int {
	return LastQueueSize
}

// ====== Función de recarga de configuración de Prometheus ======

// ReloadPrometheusConfig: Realiza una solicitud POST para recargar la configuración de Prometheus
func ReloadPrometheusConfig() {
	updateConfigPrometheusURL := "http://prometheus:9090/-/reload"

	req, err := http.NewRequest("POST", updateConfigPrometheusURL, nil)
	if err != nil {
		fmt.Println("Error creando la solicitud:", err)
		return
	}

	go func() {
		time.Sleep(4 * time.Second) // Espera antes de enviar la solicitud

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error haciendo la solicitud:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: Prometheus devolvió el código de estado %d\n", resp.StatusCode)
			return
		}

		fmt.Println("Configuración de Prometheus recargada exitosamente", resp)
	}()
}
