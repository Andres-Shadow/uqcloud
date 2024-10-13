package database

import (
	"log"
	models "servidor_procesamiento/Procesador/Models/Entities"
)

/*
Funciòn que se encarga de obtener diversas mètricas para el monitoreo de la plataforma
@return Retorna un mapa con los valores obtenidos en las consultas realizadas a la base de datos
*/

func GetMetrics() (map[string]interface{}, error) {
	// Inicializar el mapa
	metricas := make(map[string]interface{})

	// Obtiene la cantidad total de màquinas virtuales que hay en la base de datos
	var total_maquinas_creadas int64
	if err := DATABASE.Model(&models.Maquina_virtual{}).Count(&total_maquinas_creadas).Error; err != nil {
		log.Println("Error al contar las màquinas creadas en la base de datos: " + err.Error())
		return nil, err
	}

	// Obtiene la cantidad total de màquinas virtuales encendidas que hay en la plataforma
	var total_maquinas_encendidas int64
	if err := DATABASE.Model(&models.Maquina_virtual{}).Where("estado = ?", "Encendido").Count(&total_maquinas_encendidas).Error; err != nil {
		log.Println("Error al contar las màquinas encendidas que hay en la plataforma: " + err.Error())
		return nil, err
	}

	// Obtiene la cantidad total de usuarios registradas en la base de datos
	var total_usuarios int64
	if err := DATABASE.Model(&models.Persona{}).Count(&total_usuarios).Error; err != nil {
		log.Println("Error al contar los usuarios totales registrados: " + err.Error())
		return nil, err
	}

	// Obtiene la cantidad total de usuarios con rol "estudiante"
	var total_estudiantes int64
	if err := DATABASE.Model(&models.Persona{}).Where("rol = ?", "Estudiante").Count(&total_estudiantes).Error; err != nil {
		log.Println("Error al contar los usuarios con rol estudiante: " + err.Error())
		return nil, err
	}

	// Obtiene la cantidad total de usuarios con rol "invitado"
	var total_invitados int64
	if err := DATABASE.Model(&models.Persona{}).Where("rol = ?", "Invitado").Count(&total_invitados).Error; err != nil {
		log.Println("Error al contar los usuarios con rol invitado: " + err.Error())
		return nil, err
	}

	// Obtiene la cantidad total de memoria RAM que tiene la plataforma
	var total_RAM int64
	if err := DATABASE.Model(&models.Host{}).Select("COALESCE(SUM(ram_total), 0)").Scan(&total_RAM).Error; err != nil {
		log.Println("Error al contar el total de memoria RAM que tiene disponible la plataforma: " + err.Error())
		return nil, err
	}

	// Obtiene la cantidad total de memoria RAM que estàn usando las màquinas virtuales encendidas
	var total_RAM_usada int64
	if err := DATABASE.Model(&models.Maquina_virtual{}).Select("COALESCE(SUM(ram), 0)").Scan(&total_RAM_usada).Error; err != nil {
		log.Println("Error al contar el total de memoria RAM que estàn usando las màquinas encendidas: " + err.Error())
		return nil, err
	}

	// Obtiene la cantidad total de CPU que tiene la plataforma
	var total_CPU int64
	if err := DATABASE.Model(&models.Host{}).Select("COALESCE(SUM(cpu_total), 0)").Scan(&total_CPU).Error; err != nil {
		log.Println("Error al contar el total de CPU que tiene disponible la plataforma: " + err.Error())
		return nil, err
	}

	// Obtiene la cantidad total de CPU que estàn usando las màquinas virtuales encendidas
	var total_CPU_usada int64
	if err := DATABASE.Model(&models.Maquina_virtual{}).Select("COALESCE(SUM(cpu), 0)").Scan(&total_CPU_usada).Error; err != nil {
		log.Println("Error al contar el total de CPU que estàn usando las màquinas encendidas: " + err.Error())
		return nil, err
	}

	metricas["total_maquinas_creadas"] = total_maquinas_creadas
	metricas["total_maquinas_encendidas"] = total_maquinas_encendidas
	metricas["total_usuarios"] = total_usuarios
	metricas["total_estudiantes"] = total_estudiantes
	metricas["total_invitados"] = total_invitados
	metricas["total_RAM"] = total_RAM / 1024             // Se divide por 1024 para pasar de Mb a Gb
	metricas["total_RAM_usada"] = total_RAM_usada / 1024 // Se divide por 1024 para pasar de Mb a Gb
	metricas["total_CPU"] = total_CPU
	metricas["total_CPU_usada"] = total_CPU_usada

	return metricas, nil
}
