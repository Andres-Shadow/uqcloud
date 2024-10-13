package utilities

import (
	"log"
	"math/rand"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models/Entities"
	"time"

	"golang.org/x/crypto/bcrypt"
)

/*
Clase encargada de contener las herramientas sobre las máquinas virtuales para invitados
*/

/*
Funciòn que se encarga de crear cuentas temporales para usuarios invitados
Cuando se crea la cuenta e ingresa a la base de datos, se encarga de invocar la funciòn para crear una màquina virtual temporal.
@clientIP Paràmetro que contiene la direcciòn IP desde la cual se està realizando la solicitud de crear la cuenta temporal
*/

func CreateTempAccount() string {

	persona := models.Persona{
		Nombre:   "Usuario",
		Apellido: "Invitado",
		Email:    generateRandomEmail(),
		Rol:      0,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("guestuqcloud"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al encriptar la contraseña: %v\n", err)
		return ""
	}

	persona.Contrasenia = string(hashedPassword)

	if !database.CreateNewUser(persona) {
		log.Printf("Error al crear la cuenta temporal para el usuario invitado: %s - %s\n", persona.Email, persona.Nombre)
		return ""
	}

	return persona.Email
}

/*
Funciòn que se encarga de generar un correo aleatorio para las cuentas de lo sinvitados las cuales son temporales
*/

func generateRandomEmail() string {
	email := GenerateRandomString(5) + "@temp.com"
	return email
}

/*
Funciòn que genera una cadena alfanumèrica aleatoria
@length Paràmetro que contiene la longitud de la cadena a generar.
@Return Retorna la cadena aleatoria generada
*/
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// seededRand utiliza un generador de números aleatorios con una semilla basada en el tiempo actual en nanosegundos.
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		// Selecciona un carácter aleatorio del conjunto de caracteres
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
