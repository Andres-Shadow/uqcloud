package handlers

import "os"

// Se declara una variable global para almacenar la URL del servidor de procesamiento
var ServidorProcesamientoRoute string

// Se inicializa la variable global en la función init
func init() {
	// Verifica si la variable de entorno "servidor_procesamiento" está definida
	ServidorProcesamientoRoute = os.Getenv("servidor_procesamiento")
	if ServidorProcesamientoRoute == "" {
		ServidorProcesamientoRoute = "localhost"
	}
}
