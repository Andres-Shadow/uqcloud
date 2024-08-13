package database

import (
	_"database/sql"
	"log"
	models "servidor_procesamiento/Procesador/Models"

	"gorm.io/gorm"
)

/*
Funciòn que permite obtener un disco que cumpla con los paràmetros especificados
@sistema_operativo Paràmetro que representa el tipo de sistema operativo que debe tener el disco
@distribucion_sistema_operativo Paràmetro que representa la distribuciòn del sistema operativo
@id_host Paràmetro que representa el identificador ùnico del host en el cual se està buscando el disco
@Return Retorna el disco en caso de que exista y cumpla con las condiciones mencionadas anterormente
*/
// func GetDisk(sistema_operativo string, distribucion_sistema_operativo string, id_host int) (models.Disco, error) {

// 	var disco models.Disco
// 	err := database.DB.QueryRow("Select * from disco where sistema_operativo = ? and distribucion_sistema_operativo =? and host_id = ?", sistema_operativo, distribucion_sistema_operativo, id_host).Scan(&disco.Id, &disco.Nombre, &disco.Ruta_ubicacion, &disco.Sistema_operativo, &disco.Distribucion_sistema_operativo, &disco.arquitectura, &disco.Host_id)
// 	err := DB.QueryRow("Select * from disco where sistema_operativo = ? and distribucion_sistema_operativo =? and host_id = ?", sistema_operativo, distribucion_sistema_operativo, id_host).Scan(&disco.Id, &disco.Nombre, &disco.Ruta_ubicacion, &disco.Sistema_operativo, &disco.Distribucion_sistema_operativo, &disco.Arquitectura, &disco.Host_id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Println("No se encontrò un disco: " + sistema_operativo + " " + distribucion_sistema_operativo)
// 		} else {
// 			log.Println("Hubo un error al realizar la consulta: " + err.Error())
// 		}
// 		return disco, err
// 	}
// 	return disco, nil
// }


func GetDisk(sistema_operativo string, distribucion_sistema_operativo string, id_host int) (models.Disco, error) {

	var disco models.Disco
	err := DATABASE.Where("sistema_operativo = ? AND distribucion_sistema_operativo = ? AND host_id = ?", sistema_operativo, distribucion_sistema_operativo, id_host).First(&disco).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No se encontró un disco: " + sistema_operativo + " " + distribucion_sistema_operativo)
		} else {
			log.Println("Hubo un error al realizar la consulta: " + err.Error())
		}
		return disco, err
	}
	return disco, nil
}
