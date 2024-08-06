package database

import (
	"database/sql"
	"log"
	models "servidor_procesamiento/Procesador/Models"
)

/*
Clase encarga de contener los elementos relacionados a las consultas sobre la base de datos sobre la tabla de usuarios
*/

/*
Funciòn que permite obtener un usuario dado su identificador ùnico, es decir, su email
@email Paràmetro que representa el email del usuario a buscar
@Return Retorna el usuario (Persona) en caso de que exista un usuario con ese email
*/
func GetUser(email string) (models.Persona, error) {

	var persona models.Persona
	err := DB.QueryRow("SELECT * FROM persona WHERE email = ?", email).Scan(&persona.Email, &persona.Nombre, &persona.Apellido, &persona.Contrasenia, &persona.Rol)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No se encontrò un usuario con el email especificado")
		} else {
			log.Println("Hubo un error al realizar la consulta: " + err.Error())
		}
		return persona, err
	}
	return persona, nil
}

/*
Funciòn que permite eliminar una cuenta de un usuario de la base de datos
@email Paràmetro que contiene el email del usuario a eliminar
*/

// func deleteAccount(email string) {

// 	//Elimina la cuenta de la base de datos
// 	err := DB.QueryRow("DELETE FROM persona WHERE email = ?", email)
// 	if err == nil {
// 		log.Println("Error al eliminar el registro de la base de datos: ", err)
// 	}
// }

/*
Funciòn que permite conocer el total de màquianas virtuales que tiene creadas un usuario
@email Paràmetro que contiene el email del usuario al cual se le va a contar las mpaquinas que tiene creadas
@return retorna un entero con el nùmero de màquinas creadas
*/

// func countUserMachinesCreated(email string) (int, error) {

// 	//Obtiene la cantidad total de hosts que hay en la base de datos
// 	var count int
// 	err := DB.QueryRow("SELECT COUNT(*) FROM maquina_virtual where persona_email = ?", email).Scan(&count)
// 	if err != nil {
// 		log.Println("Error al contar las màquinas del usuario que hay en la base de datos: " + err.Error())
// 		return 0, err
// 	}

// 	return count, nil
// }



// func GetUser(email string) (models.Persona, error) {
// 	var persona models.Persona
// 	err := GormDB.Where("email = ?", email).First(&persona).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			log.Println("No se encontró un usuario con el email especificado")
// 		} else {
// 			log.Println("Hubo un error al realizar la consulta: " + err.Error())
// 		}
// 		return persona, err
// 	}
// 	return persona, nil
// }

// func deleteAccount(email string) {
// 	err := GormDB.Where("email = ?", email).Delete(&models.Persona{}).Error
// 	if err != nil {
// 		log.Println("Error al eliminar el registro de la base de datos: ", err)
// 	}
// }

// func countUserMachinesCreated(email string) (int64, error) {
// 	var count int64
// 	err := GormDB.Model(&models.MaquinaVirtual{}).Where("persona_email = ?", email).Count(&count).Error
// 	if err != nil {
// 		log.Println("Error al contar las máquinas del usuario que hay en la base de datos: " + err.Error())
// 		return 0, err
// 	}
// 	return count, nil
// }