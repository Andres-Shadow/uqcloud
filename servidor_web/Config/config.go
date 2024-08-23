package Config

import (
	"fmt"
	"os"
)

// Se declara una variable global para almacenar la URL del servidor de procesamiento
var ServidorProcesamientoRoute string

// TODO: Deberia moverse a un package solo de config
// Se inicializa la variable global en la función init
func init() {
	// Verifica si la variable de entorno "servidor_procesamiento" está definida
	ServidorProcesamientoRoute = os.Getenv("servidor_procesamiento")
	fmt.Println("Servidor de procesamiento: ", ServidorProcesamientoRoute)
	if ServidorProcesamientoRoute == "" {
		ServidorProcesamientoRoute = "localhost"
	}
}
