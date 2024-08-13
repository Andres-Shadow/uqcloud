package database

import (
	_"database/sql"
	"errors"
	"log"
	"math/rand"
	models "servidor_procesamiento/Procesador/Models"
	"time"
)

/*
Clase encargada de contener los elementos relacioados con las consultas sobre la base de datos
relacionados con la tabla Host de la base de datos
*/

// Función que consulta la lista completa de hsot en la base de datos
// Return lista de objetos host
// func ConsultHosts() ([]models.Host, error) {

// 	var query string
// 	var rows *sql.Rows
// 	var err error

// 	query = "SELECT id, nombre from host"
// 	rows, err = DB.Query(query)

// 	var hosts []models.Host

// 	if err != nil {
// 		log.Println("Error al realizar la consulta de màquinas en la BD", err)
// 		return hosts, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var host models.Host
// 		if err := rows.Scan(&host.Id, &host.Nombre); err != nil {
// 			// Manejar el error al escanear la fila
// 			continue
// 		}
// 		hosts = append(hosts, host)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Println("Error al iterar sobre las filas ", err)
// 		return hosts, err
// 	}

// 	if len(hosts) == 0 {
// 		// No se encontraron máquinas virtuales para el usuario
// 		return hosts, errors.New("no Machines Found")
// 	}
// 	return hosts, nil
// }


/*
Funciòn que contiene el algoritmo de asignaciòn tipo aleatorio. Se encarga de escoger un host de la base de datos al azar
Return host seleccionado por el algoritmo
*/
// func SelectHost() (models.Host, error) {

// 	var host models.Host
// 	// Consulta para contar el número de registros en la tabla "host"
// 	var count int
// 	err := DB.QueryRow("SELECT COUNT(*) FROM host").Scan(&count)
// 	if err != nil {
// 		log.Println("Error al realizar la consulta: " + err.Error())
// 		return host, err
// 	}

// 	// Genera un número aleatorio dentro del rango de registros
// 	rand.New(rand.NewSource(time.Now().Unix())) // Seed para generar números aleatorios diferentes en cada ejecución
// 	randomIndex := rand.Intn(count)

// 	// Consulta para seleccionar un registro aleatorio de la tabla "host"
// 	err = DB.QueryRow("SELECT * FROM host ORDER BY RAND() LIMIT 1 OFFSET ?", randomIndex).Scan(&host.Id, &host.Nombre, &host.Mac, &host.Ip, &host.Hostname, &host.Ram_total, &host.Cpu_total, &host.Almacenamiento_total, &host.Ram_usada, &host.Cpu_usada, &host.Almacenamiento_usado, &host.Adaptador_red, &host.Estado, &host.Ruta_llave_ssh_pub, &host.Sistema_operativo, &host.Distribucion_sistema_operativo)
// 	if err != nil {
// 		log.Println("Error al realizar la consulta sql: ", err)
// 		return host, err
// 	}

// 	// Imprime el registro aleatorio seleccionado
// 	fmt.Printf("Registro aleatorio seleccionado: ")
// 	fmt.Printf("ID: %d, Nombre: %s, IP: %s\n", host.Id, host.Nombre, host.Ip)

// 	return host, nil
// }


func ConsultHosts() ([]models.Host, error) {

	var hosts []models.Host
	err := DATABASE.Select("id, nombre").Find(&hosts).Error

	if err != nil {
		log.Println("Error al realizar la consulta de máquinas en la BD", err)
		return hosts, err
	}

	if len(hosts) == 0 {
		// No se encontraron máquinas virtuales para el usuario
		return hosts, errors.New("no Machines Found")
	}
	return hosts, nil
}

func SelectHost() (models.Host, error) {

	var host models.Host
	var count int64

	err := DATABASE.Model(&models.Host{}).Count(&count).Error
	if err != nil {
		log.Println("Error al realizar la consulta: " + err.Error())
		return host, err
	}

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