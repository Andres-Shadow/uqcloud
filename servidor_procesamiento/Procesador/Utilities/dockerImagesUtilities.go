package utilities

import (
	"fmt"
	"log"
)

// función que elimia una imagen docker dentro de una maquina virtual
// params: imagen, ip, hostname
// returns: string (mensaje de confirmación)

func EliminarImagen(imagen, ip, hostname string) string {

	fmt.Println("Eliminar Imagen: ", imagen)

	sctlCommand := "docker rmi " + imagen

	config, err := ConfigurarSSHContrasenia(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := EnviarComandoSSH(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return "Comando"
}

// función que elimina todas las imagenes docker dentro de una maquina virtual
// params: ip, hostname
// returns: string (mensaje de confirmación)
func EliminarTodasImagenes(ip, hostname string) string {

	sctlCommand := "docker rmi $(docker images -a -q)"

	config, err := ConfigurarSSHContrasenia(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := EnviarComandoSSH(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return "Comando"
}