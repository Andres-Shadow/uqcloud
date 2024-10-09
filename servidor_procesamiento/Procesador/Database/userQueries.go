package database

import (
	_ "database/sql"
	"fmt"
	"log"
	models "servidor_procesamiento/Procesador/Models"

	"gorm.io/gorm"
)

/*
Clase encarga de contener los elementos relacionados a las consultas sobre la base de datos sobre la tabla de usuarios
*/

func CountAdminsRegistered() bool {
	var count int64
	err := DATABASE.Model(&models.Persona{}).Where("rol = ?", 1).Count(&count).Error
	if err != nil {
		log.Println("Error al contar los administradores que hay en la base de datos: " + err.Error())
		return false
	}
	return count > 0
}

func GetUser(email string) (models.Persona, error) {
	var persona models.Persona
	err := DATABASE.Where("email = ?", email).First(&persona).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No se encontró un usuario con el email especificado")
		} else {
			log.Println("Hubo un error al realizar la consulta: " + err.Error())
		}
		return persona, err
	}
	return persona, nil
}

func GetUserPassword(email string) (string, error) {
	var persona models.Persona
	err := DATABASE.Where("email = ?", email).First(&persona).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No se encontró un usuario con el email especificado")
		} else {
			log.Println("Hubo un error al realizar la consulta: " + err.Error())
		}
		return "", err
	}
	return persona.Contrasenia, nil

}

func GetUserFromEmail(email string) (models.Persona, error) {
	var persona models.Persona
	err := DATABASE.Where("email = ?", email).First(&persona).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No se encontró un usuario con el email especificado")
		} else {
			log.Println("Hubo un error al realizar la consulta: " + err.Error())
		}
		return persona, err
	}
	return persona, nil
}

func CreateNewUser(persona models.Persona) bool {
	err := DATABASE.Create(&persona).Error
	if err != nil {
		log.Println("Error al insertar el nuevo registro de persona en la base de datos: ", err)
		return false
	}

	// Eliminar usuarios invitados antiguos (que tengan mas de 6 horas desde su creación)
	DATABASE.Exec("DELETE FROM persona WHERE rol = 0 AND TIMESTAMPDIFF(HOUR, created_at, NOW()) > 6")

	return true
}

// Funcion para precargar el usuario administrador
func CreateAdmin() {
	if !CountAdminsRegistered() {
		persona := models.Persona{
			Nombre:      "admin",
			Apellido:    "admin",
			Email:       "admin@uqcloud.co",
			Contrasenia: "$2y$10$JGxWitiykfO83Ep8IBab/.3fn.H/DxMjAK8dFTQCPZyJ5EHqZtfji", // Dejar este hash bcrypt para la contraseña "admin"
			Rol:         1,
		}

		DATABASE.Create(&persona)
		fmt.Print("Usuario administrador creado\n")
	}
}

