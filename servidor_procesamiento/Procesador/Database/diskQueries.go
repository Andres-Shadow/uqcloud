package database

import (
	_ "database/sql"
	"log"
	models "servidor_procesamiento/Procesador/Models/Entities"
	"strconv"

	"gorm.io/gorm"
)

/*
Funciòn que permite obtener un disco que cumpla con los paràmetros especificados
@sistema_operativo Paràmetro que representa el tipo de sistema operativo que debe tener el disco
@distribucion_sistema_operativo Paràmetro que representa la distribuciòn del sistema operativo
@id_host Paràmetro que representa el identificador ùnico del host en el cual se està buscando el disco
@Return Retorna el disco en caso de que exista y cumpla con las condiciones mencionadas anterormente
*/

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

func CreateDisck(disco models.Disco) error {
	err := DATABASE.Create(&disco).Error
	if err != nil {
		log.Println("Hubo un error al crear el disco: " + err.Error())
		return err
	}
	return nil
}

func ListUniquesDisks() ([]models.Disco, error) {
	var disks []models.Disco
	err := DATABASE.Select("DISTINCT distribucion_sistema_operativo").Find(&disks).Error
	if err != nil {
		log.Println("Hubo un error al listar los discos únicos: " + err.Error())
		return disks, err
	}
	return disks, nil
}

func ListHostWhereDiskExists(diskSo string) ([]models.Host, error) {
	var hosts []models.Host
	err := DATABASE.Joins("JOIN disco ON host.id = disco.host_id").Where("disco.distribucion_sistema_operativo = ?", diskSo).Find(&hosts).Error
	if err != nil {
		log.Println("Hubo un error al listar los hosts donde existe el disco: " + err.Error())
		return hosts, err
	}
	return hosts, nil
}

func DeleteDiskFromHost(hostId string, diskSo string) error {
	id, _ := strconv.Atoi(hostId)
	err := DATABASE.Unscoped().Where("host_id = ? AND distribucion_sistema_operativo = ?", id, diskSo).Delete(&models.Disco{}).Error
	if err != nil {
		log.Println("Hubo un error al eliminar el disco: " + err.Error())
		return err
	}
	return nil
}
