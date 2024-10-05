package utilities

import (
	"fmt"
	"math/rand"
	"time"

	"os"
	models "servidor_procesamiento/Procesador/Models"
	"strconv"
)

// función que recibe un arreglo de ips de los host a registrar rapido
// y generando nombres aleatorios para los host, los registra junto a sus discos
func FastRegisterHosts(ips []string) {

	//TODO hacer la implementación de la función
	var quickHost models.Host
	quickHost.Mac = "0A-00-27-00-00-0A"
	quickHost.Hostname = os.Getenv("QUICK_HOST_HOSTNAME")
	quickHost.Ram_total, _ = strconv.Atoi(os.Getenv("QUICK_HOST_RAM"))
	quickHost.Cpu_total, _ = strconv.Atoi(os.Getenv("QUICK_HOST_CPU"))
	quickHost.Almacenamiento_total, _ = strconv.Atoi(os.Getenv("QUICK_HOST_ALMACENAMIENTO"))
	quickHost.Ram_usada = 0
	quickHost.Cpu_usada = 0
	quickHost.Almacenamiento_usado = 0
	quickHost.Adaptador_red = os.Getenv("QUICK_HOST_NETWORK")
	quickHost.Ruta_llave_ssh_pub = os.Getenv("QUICK_HOST_SSH_ROUTE")
	quickHost.Estado = "apagado"
	quickHost.Sistema_operativo = os.Getenv("QUICK_HOST_SO")
	quickHost.Distribucion_sistema_operativo = os.Getenv("QUICK_HOST_SO_DISTRO")

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
