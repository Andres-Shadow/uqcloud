package database

import (
	"errors"
	"log"
	models "servidor_procesamiento/Procesador/Models"
)

/*
Funciòn que permite consultar el catàlogo de servicios de la plataforma
@Return Retorna un arreglo con los catàlogos disponibles de la plataforma
*/
func ConsultCatalog() ([]models.Catalogo, error) {

	var listaCatalogo []models.Catalogo

    // Realiza la consulta con GORM usando Join
    err := DATABASE.Table("catalogo_disco cd").
        Select("c.id, c.nombre, c.ram, c.cpu, d.sistema_operativo, d.distribucion_sistema_operativo, d.arquitectura").
        Joins("JOIN catalogo c ON cd.catalogo_id = c.id").
        Joins("JOIN disco d ON cd.disco_id = d.id").
        Scan(&listaCatalogo).Error

    if err != nil {
        log.Println("Error al realizar la consulta del catálogo en la base de datos:", err)
        return listaCatalogo, err
    }

    if len(listaCatalogo) == 0 {
        log.Println("El catálogo está vacío")
        return nil, errors.New("el catalogo se encuentra vacío")
    }

    return listaCatalogo, nil
}