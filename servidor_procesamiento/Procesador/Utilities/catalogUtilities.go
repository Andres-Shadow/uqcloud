package utilities

import (
	"fmt"
	"log"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
)

/*
Clase encargada de contener las funcionalidades a realizar sobre los catálogos
*/

/*
Funciòn que permite consultar el catàlogo de servicios de la plataforma
@Return Retorna un arreglo con los catàlogos disponibles de la plataforma
*/
func ConsultCatalog() ([]models.Catalogo, error) {

	var catalogo models.Catalogo
	var listaCatalogo []models.Catalogo

	query := "SELECT c.id, c.nombre, c.ram, c.cpu, d.sistema_operativo, d.distribucion_sistema_operativo, d.arquitectura FROM catalogo_disco cd JOIN catalogo c ON cd.catalogo_id = c.id JOIN disco d ON cd.disco_id = d.id"
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Error al realizar la consulta del catàlogo en la base de datos")
		return listaCatalogo, err
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&catalogo.Id, &catalogo.Nombre, &catalogo.Ram, &catalogo.Cpu, &catalogo.Sistema_operativo, &catalogo.Distribucion_sistema_operativo, &catalogo.Arquitectura); err != nil {
			log.Println("Error al obtener la fila")
			continue
		}
		listaCatalogo = append(listaCatalogo, catalogo)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error al obtener el catàlogo de màquinas virtuales")
		return listaCatalogo, err
	}

	if len(listaCatalogo) == 0 {
		fmt.Println("El catàlogo està vacìo")
	}

	return listaCatalogo, nil
}