package utilities

import (
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

/*
Clase encarga de contener todos los metodos que representan una funcionalidad
interna de la aplicacion que son llamadas dentro del programa y se encuentran relacionaodas con la
gestion de archivos
*/

// Función para cargar la llave privada desde un archivo
// @file Parámetro que contiene la ruta de la llave privada
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("Error al leer la llave privada:", err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatal("Error al analizar la llave privada:", err)
	}

	return ssh.PublicKeys(key)
}

/*
Esta funciòn carga y devuelve la llave privada SSH desde la ruta especificada
@file Paràmetro que contiene la ruta de la llave privada
*/
func privateKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(key), nil
}
