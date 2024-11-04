package hostutilities

import (
	"fmt"
	"math/rand"
	"time"

	config "servidor_procesamiento/Procesador/Config"
	models "servidor_procesamiento/Procesador/Models/Entities"
)

// función que recibe un arreglo de ips de los host a registrar rapido
// y generando nombres aleatorios para los host, los registra junto a sus discos
func FastRegisterHosts(ips []string) {

	var quickHost models.Host

	// Data obtained from the configuration file
	quickHost.Hostname = config.QUICK_HOST_HOSTNAME
	quickHost.Ram_total = config.QUICK_HOST_RAM
	quickHost.Cpu_total = config.QUICK_HOST_CPU
	quickHost.Almacenamiento_total = config.QUICK_HOST_STORAGE
	quickHost.Adaptador_red = config.QUICK_HOST_NETWORK
	quickHost.Sistema_operativo = config.QUICK_HOST_SO
	quickHost.Distribucion_sistema_operativo = config.QUICK_HOST_SO_DISTRO
	
	// Default values for the used aspects in the hosts
	quickHost.Ram_usada = 0
	quickHost.Cpu_usada = 0
	quickHost.Almacenamiento_usado = 0
	quickHost.Estado = "apagado"
	
	for _, ip := range ips {
		quickHost.Ip = ""
		quickHost.Ip = ip
		quickHost.Nombre = ""
		quickHost.Nombre = generateRandomName()

		fmt.Println("--------------------")
		SetUpHostAndDisk(quickHost)

	}
}

// Listas de adjetivos y nombres
var adjectives = []string{
	"amazing", "brave", "clever", "daring", "eager", "fearless", "graceful", "heroic", "intrepid", "jovial",
}

var surnames = []string{
	"curie", "einstein", "newton", "tesla", "darwin", "hawking", "turing", "lovelace", "galileo", "feynman",
}

// Función para generar nombre aleatorio
func generateRandomName() string {
	// Inicializar el generador de números aleatorios con la semilla actual
	randomIndex := generateRandonNumber()

	// Elegir un adjetivo y un apellido aleatorios
	adjective := adjectives[randomIndex]
	surname := surnames[randomIndex]

	// Combinar el adjetivo y el apellido
	return fmt.Sprintf("%s_%s", adjective, surname)
}

func generateRandonNumber() int {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(len(adjectives))
}
