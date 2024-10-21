package models

import (
	"container/list"
	"sync"
)

/*
Clase que se encarga de contener las colas de gestion de
solicitudes que llegan al servidor de procesamiento
*/

// Cola de especificaciones para la gestiòn de màquinas virtuales
// La gestiòn puede ser: modificar, eliminar, iniciar, detener una MV.
type ManagementQueue struct {
	sync.Mutex
	Queue *list.List
}


type Maquina_virtualQueue struct {
	sync.Mutex
	Queue *list.List
}


type Docker_imagesQueue struct {
	sync.Mutex
	Queue *list.List
}

type Docker_contenedorQueue struct {
	sync.Mutex
	Queue *list.List
}