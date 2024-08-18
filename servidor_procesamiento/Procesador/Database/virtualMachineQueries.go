package database

import (
	_ "database/sql"
	"errors"
	"log"

	models "servidor_procesamiento/Procesador/Models"
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

	if persona.Rol == 1 {
		err = DATABASE.Table("maquina_virtual as m").
			Select("m.nombre, m.ram, m.cpu, m.ip, m.estado, d.sistema_operativo, d.distribucion_sistema_operativo, m.hostname").
			Joins("INNER JOIN disco as d on m.disco_id = d.id").
			Scan(&machines).Error

	} else {
		err = DATABASE.Model(&models.Maquina_virtual{}).Select("m.nombre, m.ram, m.cpu, m.ip, m.estado, d.sistema_operativo, d.distribucion_sistema_operativo, m.hostname").
			Joins("INNER JOIN disco as d on m.disco_id = d.id").
			Where("m.persona_email = ?", persona.Email).
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
	err := DATABASE.Model(&models.Maquina_virtual{}).Where("nombre = ?", nombre).Update("estado", newState).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateVirtualMachineIP(newIP, nombre string) error {
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

func ExistVirtualMachine(virtualMachineName string)(bool, error) {
	err := DATABASE.Where("nombre = ?", virtualMachineName).First(&models.Maquina_virtual{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
Funciòn que permite obtener todas las màquinas virtuales creadas en la plataforma por los usuarios con rol invitado
@return Retorna un arreglo con todas las màquinas encontradas
*/

// func GetGuestMachines() ([]models.Maquina_virtual, error) {
// 	var maquinas []models.Maquina_virtual

// 	query := "SELECT m.nombre, m.fecha_creacion, m.host_id, m.persona_email FROM maquina_virtual m JOIN persona p ON m.persona_email = p.email WHERE p.rol = ?;"

// 	rows, err := DB.Query(query, "Invitado")
// 	if err != nil {
// 		log.Println("Error al consultar las máquinas de los invitados:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var machine models.Maquina_virtual
// 		var fechaCreacionStr string

// 		if err := rows.Scan(&machine.Nombre, &fechaCreacionStr, &machine.Host_id, &machine.Persona_email); err != nil {
// 			log.Println("Error al escanear la fila:", err)
// 			continue
// 		}

// 		// Convierte la cadena de fecha y hora a un tipo de datos time.Time
// 		fechaCreacion, err := time.Parse("2006-01-02 15:04:05", fechaCreacionStr)
// 		if err != nil {
// 			log.Println("Error al convertir la fecha y hora:", err)
// 			continue
// 		}

// 		// Asigna la fecha y hora convertida a la estructura Maquina_virtual
// 		machine.Fecha_creacion = fechaCreacion

// 		maquinas = append(maquinas, machine)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Println("Error al iterar sobre las filas:", err)
// 		return nil, err
// 	}

// 	return maquinas, nil
// }

/*
Funciòn que permite obtener un host dado su identificador ùnico
@idHost Paràmetro que representa el identificador ùnico del host a buscar
@Return Retorna el host encontrado
*/
// func GetHost(idHost int) (models.Host, error) {

// 	var host models.Host
// 	err := DB.QueryRow("SELECT * FROM host WHERE id = ?", idHost).Scan(&host.Id, &host.Nombre, &host.Mac, &host.Ip, &host.Hostname, &host.Ram_total, &host.Cpu_total, &host.Almacenamiento_total, &host.Ram_usada, &host.Cpu_usada, &host.Almacenamiento_usado, &host.Adaptador_red, &host.Estado, &host.Ruta_llave_ssh_pub, &host.Sistema_operativo, &host.Distribucion_sistema_operativo)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Println("No se encontró el host con el nombre especificado.")
// 		} else {
// 			log.Println("Error al realizar la consulta: ", err)
// 		}
// 		return host, err
// 	}
// 	return host, nil
// }

/*
Funciòn que permite obtener una màquina virtual dado su nombre
@nameVM Paràmetro que representa el nombre de la màquina virtual a buscar
@Retorna la màquina virtual en caso de que exista en la base de datos
*/
// func GetVM(nameVM string) (models.Maquina_virtual, error) {
// 	var maquinaVirtual models.Maquina_virtual

// 	var fechaCreacionStr string
// 	err := DB.QueryRow("SELECT * FROM maquina_virtual WHERE nombre = ?", nameVM).Scan(
// 		&maquinaVirtual.Uuid, &maquinaVirtual.Nombre, &maquinaVirtual.Ram,
// 		&maquinaVirtual.Cpu, &maquinaVirtual.Ip, &maquinaVirtual.Estado,
// 		&maquinaVirtual.Hostname, &maquinaVirtual.Persona_email,
// 		&maquinaVirtual.Host_id, &maquinaVirtual.Disco_id,
// 		&fechaCreacionStr)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Println("No se encontró la máquina virtual con el nombre especificado.")
// 		} else {
// 			log.Println("Hubo un error al realizar la consulta:", err)
// 		}
// 		return maquinaVirtual, err
// 	}

// 	fechaCreacion, err := time.Parse("2006-01-02 15:04:05", fechaCreacionStr)
// 	if err != nil {
// 		log.Println("Error al parsear la fecha de creación:", err)
// 		return maquinaVirtual, err
// 	}
// 	maquinaVirtual.Fecha_creacion = fechaCreacion

// 	return maquinaVirtual, nil
// }

/*
Funciòn que permite conocer las màquinas virtuales que tiene creadas un usuario ò todas las màquinas de la plataforma si es un administrador
@persona Paràmetro que representa un usuario, al cual se le van a buscar las màquinas que le pertenece
@return Retorna un arreglo con las màquinas que le pertenecen al usuario.
*/

// func ConsultMachines(persona models.Persona) ([]models.Maquina_virtual, error) {

// 	email := persona.Email
// 	var query string
// 	var rows *sql.Rows
// 	var err error

// 	if persona.Rol == 1 {
// 		//Consulta todas las màquinas virtuales de la base de datos
// 		query = "SELECT m.nombre, m.ram, m.cpu, m.ip, m.estado, d.sistema_operativo, d.distribucion_sistema_operativo, m.hostname FROM maquina_virtual as m INNER JOIN disco as d on m.disco_id = d.id"
// 		rows, err = DB.Query(query)
// 	} else {
// 		//Consulta las màquinas virtuales de un usuario en la base de datos
// 		query = "SELECT m.nombre, m.ram, m.cpu, m.ip, m.estado, d.sistema_operativo, d.distribucion_sistema_operativo, m.hostname FROM maquina_virtual as m INNER JOIN disco as d on m.disco_id = d.id WHERE m.persona_email = ?"
// 		rows, err = DB.Query(query, email)
// 	}

// 	var machines []models.Maquina_virtual

// 	if err != nil {
// 		log.Println("Error al realizar la consulta de màquinas en la BD", err)
// 		return machines, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var machine models.Maquina_virtual
// 		if err := rows.Scan(&machine.Nombre, &machine.Ram, &machine.Cpu, &machine.Ip, &machine.Estado, &machine.Sistema_operativo, &machine.Distribucion_sistema_operativo, &machine.Hostname); err != nil {
// 			// Manejar el error al escanear la fila
// 			continue
// 		}
// 		machines = append(machines, machine)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Println("Error al iterar sobre las filas ", err)
// 		return machines, err
// 	}

// 	if len(machines) == 0 {
// 		// No se encontraron máquinas virtuales para el usuario
// 		return machines, errors.New("no Machines Found")
// 	}
// 	return machines, nil
// }
