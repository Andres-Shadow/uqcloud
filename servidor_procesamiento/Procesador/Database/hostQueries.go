package database

import (
	_ "database/sql"
	"errors"
	"log"
	"math/rand"
	models "servidor_procesamiento/Procesador/Models/Entities"
	"time"
)

/*
Clase encargada de contener los elementos relacioados con las consultas sobre la base de datos
relacionados con la tabla Host de la base de datos
*/

func ConsultHosts() ([]map[string]interface{}, error) {
	//mapa que almacena el id y el nombre de las máquinas
	//id: x
	//nombre: y
	//para posteriormente ser utliizado en la respuesta
	var results []map[string]interface{}

	// Realiza la consulta y guarda los resultados directamente en una lista de mapas
	// Nota: Está estructura del map se debe usar así para obtener los valores de la BD
	// 		 Luego, el map cambia cuando toca pasar los key:valor a la estructura del servidor web
	err := DATABASE.Model(&models.Host{}).Select("id, nombre, ip, ram_total, cpu_total, almacenamiento_total, estado").Find(&results).Error

	if err != nil {
		log.Println("Error al realizar la consulta de máquinas en la BD:", err)
		return nil, err
	}

	if len(results) == 0 {
		log.Println("No se encontraron máquinas registradas en la base de datos")
		return nil, errors.New("no Machines Found")
	}

	// ESO ERAAAAAAAA JASJSAJASJASJ DOS DIAS PAARA SACAR EL NOMBRE DEL HOST, DIOSMIO AJSJSAJJSJ SOY ELMOEJR JULIOPROYT777 REPORTANDO
	// Se cambia la key 'nombre' a 'hst_name', para no tener los errores de:
	// RESPONSE:  [{"id":1,"nombre":"prueba pc personal 3"}] 			|  MAL  Model.Host struct NO LO ACEPTA
	// RESPONSE:  [{"id":1,"hst_name":"prueba pc personal 3"}]			|  BIEN Model.Host struct LO ACEPTA
	for _, result := range results {
		if nombre, ok := result["nombre"]; ok {
			result["hst_name"] = nombre
			delete(result, "nombre")
		}
	}

	return results, nil
}

// funcion que registra los host en la base de datos
func AddHost(host models.Host) error {
	err := DATABASE.Create(&host).Error
	if err != nil {
		log.Println("Error al registrar el host.")
		return err
	} else {
		log.Println("Registro del host exitoso")
	}
	return nil
}

func SelectHost() (models.Host, error) {

	var host models.Host
	var count int64

	err := DATABASE.Model(&models.Host{}).Count(&count).Error
	if err != nil {
		log.Println("Error al realizar la consulta: " + err.Error())
		return host, err
	}

	log.Println("Número de hosts: ", count)

	if count == 0 {
		log.Println("No existen hosts disponibles para la creación de la VM")
		return host, err
	} else {
		// Genera un número aleatorio dentro del rango de registros
		rand.Seed(time.Now().Unix()) // Seed para generar números aleatorios diferentes en cada ejecución
		randomIndex := rand.Intn(int(count))

		// Consulta para seleccionar un registro aleatorio de la tabla "host"
		err = DATABASE.Offset(randomIndex).Limit(1).Find(&host).Error
		if err != nil {
			log.Println("Error al realizar la consulta sql: ", err)
			return host, err
		}

		// Imprime el registro aleatorio seleccionado
		log.Printf("Registro aleatorio seleccionado: ID: %d, Nombre: %s, IP: %s\n", host.Id, host.Nombre, host.Ip)
		return host, nil
	}
}

func GetHostByIp(ip string) (models.Host, error) {
	var host models.Host

	err := DATABASE.Where("ip = ?", ip).First(&host).Error
	if err != nil {
		log.Println("Error al realizar la consulta: ", err)
		return host, err
	}

	return host, nil
}

func UpdateHostRamAndCPU(idHost int, ram int, cpu int) error {
	err := DATABASE.Model(&models.Host{}).Where("id = ?", idHost).Update("ram_usada", ram).Update("cpu_usada", cpu).Error
	if err != nil {
		log.Println("Error al actualizar la información del host: ", err)
		return err
	}

	return nil
}

func CountRegisteredHosts() (int64, error) {
	var count int64

	err := DATABASE.Model(&models.Host{}).Count(&count).Error
	if err != nil {
		log.Println("Error al realizar la consulta: " + err.Error())
		return 0, err
	}
	log.Println("Número de hosts registrados: ", count)
	return count, nil
}

func UpdateHostUsedCpu(hostId int, newUserCpu int) error {
	err := DATABASE.Model(&models.Host{}).Where("id = ?", hostId).Update("cpu_usada", newUserCpu).Error
	if err != nil {
		log.Println("Error al actualizar la información del host: ", err)
		return err
	}
	return nil
}

func UpdateHostUsedRam(hostId int, newUserRam int) error {
	err := DATABASE.Model(&models.Host{}).Where("id = ?", hostId).Update("ram_usada", newUserRam).Error
	if err != nil {
		log.Println("Error al actualizar la información del host: ", err)
		return err
	}
	return nil
}

func GetHosts() []models.Host {
	var hosts []models.Host
	err := DATABASE.Find(&hosts).Error
	if err != nil {
		log.Println("Error al obtener los hosts: ", err)
	}
	return hosts
}

func DeleteHostById(id int) error {
	err := DATABASE.Where("id = ?", id).Delete(&models.Host{}).Error
	if err != nil {
		log.Println("Error al eliminar el host: ", err)
		return err
	}
	return nil
}
