package database

import (
	_ "database/sql"
	"errors"
	"log"

	models "servidor_procesamiento/Procesador/Models/Entities"
	"time"

	"gorm.io/gorm"
)

/*
Clase enccarga de contener las funciones relacionadas con la gestion de peticiones sobre las tablas
de maquina virtual y asociados
*/

func GetGuestMachines() ([]models.Maquina_virtual, error) {
	var maquinas []models.Maquina_virtual

	err := DATABASE.Model(&models.Maquina_virtual{}).Select("m.nombre, m.fecha_creacion, m.host_id, m.persona_email").
		Joins("JOIN persona p ON m.persona_email = p.email").
		Where("p.rol = ?", "Invitado").
		Scan(&maquinas).Error

	if err != nil {
		log.Println("Error al consultar las máquinas de los invitados:", err)
		return nil, err
	}

	for i := range maquinas {
		fechaCreacionStr := maquinas[i].Fecha_creacion.Format("2006-01-02 15:04:05")
		fechaCreacion, err := time.Parse("2006-01-02 15:04:05", fechaCreacionStr)
		if err != nil {
			log.Println("Error al convertir la fecha y hora:", err)
			continue
		}
		maquinas[i].Fecha_creacion = fechaCreacion
	}

	return maquinas, nil
}

func GetHost(idHost int) (models.Host, error) {
	var host models.Host
	err := DATABASE.Where("id = ?", idHost).First(&host).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No se encontró el host con el nombre especificado.")
		} else {
			log.Println("Error al realizar la consulta: ", err)
		}
		return host, err
	}
	return host, nil
}

func GetVM(nameVM string) (models.Maquina_virtual, error) {
	var maquinaVirtual models.Maquina_virtual
	err := DATABASE.Where("nombre = ?", nameVM).First(&maquinaVirtual).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No se encontró la máquina virtual con el nombre especificado.")
		} else {
			log.Println("Hubo un error al realizar la consulta:", err)
		}
		return maquinaVirtual, err
	}
	return maquinaVirtual, nil
}

func ConsultMachines(persona models.Persona) ([]models.Maquina_virtual, error) {
	var machines []models.Maquina_virtual
	var err error

	// log.Println("email: ", persona.Email)

	if persona.Rol == 1 {
		err = DATABASE.Table("maquina_virtual as m").
			Select("m.id, m.nombre, m.ram, m.cpu, m.ip, m.estado, d.sistema_operativo, d.distribucion_sistema_operativo, m.hostname, m.fecha_creacion, m.host_id, m.disco_id, m.persona_email").
			Joins("INNER JOIN disco as d on m.disco_id = d.id").
			Scan(&machines).Error

	} else {
		err = DATABASE.Table("(?) AS m", DATABASE.Model(&models.Maquina_virtual{})).Select("m.id, m.nombre, m.ram, m.cpu, m.ip, m.estado, d.sistema_operativo, d.distribucion_sistema_operativo, m.hostname, m.fecha_creacion, m.host_id, m.disco_id, m.persona_email").
			Joins("INNER JOIN disco as d on m.disco_id = d.id").
			Where("m.persona_email = ? and m.deleted_at IS NULL", persona.Email).
			Scan(&machines).Error
	}

	if err != nil {
		log.Println("Error al realizar la consulta de máquinas en la BD", err)
		return machines, err
	}

	if len(machines) == 0 {
		return machines, errors.New("no Machines Found")
	}
	return machines, nil
}

func GetStateFromVirtualMachineName(name string) (string, error) {
	var maquina models.Maquina_virtual
	err := DATABASE.Where("nombre = ?", name).First(&maquina).Error
	if err != nil {
		log.Println("Ha ocurrido un error alconsultar el estado de la máquina virtual:", err)
		return "", err
	}
	return maquina.Estado, nil
}

func CreateVirtualMachine(machine models.Maquina_virtual) error {
	err := DATABASE.Create(&machine).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateVirtualMachineCPU(newCpu int, nombre string) error {
	err := DATABASE.Model(&models.Maquina_virtual{}).Where("nombre = ?", nombre).Update("cpu", newCpu).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateVirtualMachineRam(newRam int, nombre string) error {
	err := DATABASE.Model(&models.Maquina_virtual{}).Where("nombre = ?", nombre).Update("ram", newRam).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateVirtualMachineState(nombre, newState string) error {

	var ip string
	var err error

	if newState == "Apagado" {
		ip = ""
	}

	err = DATABASE.Model(&models.Maquina_virtual{}).Where("nombre = ?", nombre).Update("estado", newState).Error
	if err != nil {
		return err
	}

	err = DATABASE.Model(&models.Maquina_virtual{}).Where("nombre = ?", nombre).Update("ip", ip).Error

	if err != nil {
		return err
	}
	return nil
}

func UpdateVirtualMachineIP(nombre, newIP string) error {
	err := DATABASE.Model(&models.Maquina_virtual{}).Where("nombre = ?", nombre).Update("ip", newIP).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteVirtualMachine(nameVM string) error {
	err := DATABASE.Where("nombre = ?", nameVM).Delete(&models.Maquina_virtual{}).Error
	if err != nil {
		return err
	}
	return nil
}

func ExistVirtualMachine(virtualMachineName string) (bool, error) {
	err := DATABASE.Where("nombre = ?", virtualMachineName).First(&models.Maquina_virtual{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetAllVirtualMachines() ([]models.Maquina_virtual, error) {
	var maquinas []models.Maquina_virtual
	err := DATABASE.Find(&maquinas).Error
	if err != nil {
		return nil, err
	}
	return maquinas, nil
}
