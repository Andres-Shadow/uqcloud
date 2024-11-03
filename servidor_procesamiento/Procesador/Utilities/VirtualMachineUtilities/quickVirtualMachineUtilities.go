package virtualmachineutilities

import (
	"log"
	"os"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models/Entities"
	userutilities "servidor_procesamiento/Procesador/Utilities/UserUtilities"

	"golang.org/x/crypto/bcrypt"
)

func CreateQuickVirtualMachine(clientIP string) string {

	persona := models.Persona{
		Nombre:   "Usuario",
		Apellido: "Invitado",
		Email:    userutilities.GenerateRandomEmail(),
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

	generateQuickVirtualMachine(persona.Email, clientIP)

	return persona.Email
}

func generateQuickVirtualMachine(email string, clientIP string) {

	distro := os.Getenv("DEFAULT_QUICK_VM_DISTRO")
	ram := config.DEFAULT_QUICK_VM_RAM
	cpu := config.DEFAULT_QUICK_VM_CPU

	maquina_virtual := models.Maquina_virtual{
		Nombre:                         "QuickGuest",
		Sistema_operativo:              "Linux",
		Distribucion_sistema_operativo: distro,
		Ram:                            ram,
		Cpu:                            cpu,
		Persona_email:                  email,
		Hostname:                       "aleatorio",
	}

	payload := map[string]interface{}{
		"specifications": maquina_virtual,
		"clientIP":       clientIP,
	}
	CreateVirtualMachineFromSpecifications(payload)
}
