package utilities

import (
	"log"
	database "servidor_procesamiento/Procesador/Database"
)

/*
Funciòn que se encarga de obtener diversas mètricas para el monitoreo de la plataforma
@return Retorna un mapa con los valores obtenidos en las consultas realizadas a la base de datos
*/

func GetMetrics() (map[string]interface{}, error) {

	var metricas map[string]interface{}

	// Inicializar el mapa
	metricas = make(map[string]interface{})

	//Obtiene la cantidad total de màquinas virtuales que hay en la base de datos
	var total_maquinas_creadas int
	err := database.DB.QueryRow("SELECT COALESCE(COUNT(*),0) FROM maquina_virtual").Scan(&total_maquinas_creadas)
	if err != nil {
		log.Println("Error al contar las màquinas creadas hay en la base de datos: " + err.Error())
		return nil, err
	}

	//Obtiene la cantidad total de màquinas virtuales encendidas que hay en la plataforma
	var total_maquinas_encendidas int
	err1 := database.DB.QueryRow("SELECT COALESCE(COUNT(*),0) FROM maquina_virtual where estado = 'Encendido'").Scan(&total_maquinas_encendidas)
	if err1 != nil {
		log.Println("Error al contar las màquinas encendidas que hay en la plataforma: " + err1.Error())
		return nil, err1
	}

	//Obtiene la cantidad total de usuarios registradas en la base de datos
	var total_usuarios int
	err2 := database.DB.QueryRow("SELECT COALESCE(COUNT(*),0) FROM persona").Scan(&total_usuarios)
	if err2 != nil {
		log.Println("Error al contar los usuarios totales registrados: " + err2.Error())
		return nil, err2
	}

	//Obtiene la cantidad total de usuarios con rol "estudiante"
	var total_estudiantes int
	err3 := database.DB.QueryRow("SELECT COALESCE(COUNT(*),0) FROM persona where rol = 'Estudiante'").Scan(&total_estudiantes)
	if err3 != nil {
		log.Println("Error al contar los usuarios con rol estudiante: " + err3.Error())
		return nil, err3
	}

	//Obtiene la cantidad total de usuarios con rol "invitado"
	var total_invitados int
	err4 := database.DB.QueryRow("SELECT COALESCE(COUNT(*),0) FROM persona where rol = 'Invitado'").Scan(&total_invitados)
	if err4 != nil {
		log.Println("Error al contar las màquinas encendidas que hay en la plataforma: " + err4.Error())
		return nil, err4
	}

	//Obtiene la cantidad total de memoria RAM que tiene la plataforma
	var total_RAM int
	err5 := database.DB.QueryRow("SELECT COALESCE(SUM(ram_total),0) AS total_ram FROM host;").Scan(&total_RAM)
	if err5 != nil {
		log.Println("Error al contar el total de memoria RAM que tiene disponible la plataforma: " + err5.Error())
		return nil, err5
	}

	//Obtiene la cantidad total de memoria RAM que estàn usando las màquinas virtuales encendidas
	var total_RAM_usada int
	err6 := database.DB.QueryRow("SELECT COALESCE(SUM(ram),0) AS total_ram_usada FROM maquina_virtual ;").Scan(&total_RAM_usada)
	if err6 != nil {
		log.Println("Error al contar el total de memoria RAM que estàn usando las màquinas encendidas: " + err6.Error())
		return nil, err6
	}

	//Obtiene la cantidad total de CPU que tiene la plataforma
	var total_CPU int
	err7 := database.DB.QueryRow("SELECT COALESCE(SUM(cpu_total),0) AS total_cpu FROM host;").Scan(&total_CPU)
	if err7 != nil {
		log.Println("Error al contar el total de CPU que tiene disponible la plataforma: " + err7.Error())
		return nil, err7
	}

	//Obtiene la cantidad total de CPU que estàn usando las màquinas virtuales encendidas
	var total_CPU_usada int
	err8 := database.DB.QueryRow("SELECT COALESCE(SUM(cpu),0) AS total_cpu_usada FROM maquina_virtual ;").Scan(&total_CPU_usada)
	if err8 != nil {
		log.Println("Error al contar el total de CPU que estàn usando las màquinas encendidas: " + err8.Error())
		return nil, err8
	}

	metricas["total_maquinas_creadas"] = total_maquinas_creadas
	metricas["total_maquinas_encendidas"] = total_maquinas_encendidas
	metricas["total_usuarios"] = total_usuarios
	metricas["total_estudiantes"] = total_estudiantes
	metricas["total_invitados"] = total_invitados
	metricas["total_RAM"] = total_RAM / 1024             //Se divide por 1024 para pasar de Mb a Gb
	metricas["total_RAM_usada"] = total_RAM_usada / 1024 //Se divide por 1024 para pasar de Mb a Gb
	metricas["total_CPU"] = total_CPU
	metricas["total_CPU_usada"] = total_CPU_usada

	return metricas, nil
}
