package utilities

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
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

func CreateTempAccount(clientIP string, distribucionSO string) string {
	
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

	CreateTempVM(persona.Email, clientIP, distribucionSO)

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

/*
Funciòn que permite crear màquina virtuales temporales para los usuarios con rol "invitado"
Esta funciòn crea las especificaciones para crear una màquina virtual con recursos mìnimos
finalmente encola la solicitud de creaciòn

@email Paràmetro que contiene el email del usuario al cual le va a pertencer la MV
@clientIP Paràmetro que contiene la direcciòn IP desde la cual se està generando la peticiòn
*/

func CreateTempVM(email string, clientIP string, distribucion_SO string) {

	maquina_virtual := models.Maquina_virtual{
		Nombre:                         "Guest",
		Sistema_operativo:              "Linux",
		Distribucion_sistema_operativo: distribucion_SO,
		Ram:                            1024,
		Cpu:                            2,
		Persona_email:                  email,
	}

	payload := map[string]interface{}{
		"specifications": maquina_virtual,
		"clientIP":       clientIP,
	}

	jsonData, _ := json.Marshal(payload) //Se codifica en formato JSON

	var decodedPayload map[string]interface{}
	err := json.Unmarshal(jsonData, &decodedPayload) //Se decodifica para meterlo en la cola
	if err != nil {
		fmt.Println("Error al decodificar el JSON:", err)
		// Manejar el error según tus necesidades
		return
	}

	// Encola la peticiòn
	config.GetMu().Lock()
	config.GetMaquina_virtualQueue().Queue.PushBack(decodedPayload)
	config.GetMu().Unlock()
}
