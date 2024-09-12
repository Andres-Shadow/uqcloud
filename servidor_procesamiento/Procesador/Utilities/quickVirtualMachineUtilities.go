package utilities

import (
	"log"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"

	"golang.org/x/crypto/bcrypt"
)

func CreateQuickVirtualMachine(clientIP string) string {

	distro := "Alpine"

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

	generateQuickVirtualMachine(persona.Email, clientIP, distro)

	return persona.Email
}

func generateQuickVirtualMachine(email string, clientIP string, distro string) {

	maquina_virtual := models.Maquina_virtual{
		Nombre:                         "Guest",
		Sistema_operativo:              "Linux",
		Distribucion_sistema_operativo: distro,
		Ram:                            1024,
		Cpu:                            1,
		Persona_email:                  email,
		Hostname:                       "aleatorio",
	}

	payload := map[string]interface{}{
		"specifications": maquina_virtual,
		"clientIP":       clientIP,
	}

	CreateVirtualMachineFromSpecifications(payload)
}
