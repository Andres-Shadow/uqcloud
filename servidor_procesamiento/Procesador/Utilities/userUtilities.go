package utilities

import (
	"fmt"
	models "servidor_procesamiento/Procesador/Models"
)

/*
Clase encargada de contener las funciones que se relacionan con las herramientas
que enseñan o manipulan informacion de los usuarios
*/

// Función que imprime la información de un usuario.
// params: estructura de persona
func PrintUserAccount(account models.Persona) {
	// Imprime la cuenta recibida.
	fmt.Printf("-------------------------\n")
	fmt.Printf("Nombre de Usuario: %s\n", account.Nombre)
	fmt.Printf("Contraseña: %s\n", account.Contrasenia)
	fmt.Printf("Email: %s\n", account.Email)

}



