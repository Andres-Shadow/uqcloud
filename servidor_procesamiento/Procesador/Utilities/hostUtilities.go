package utilities

import models "servidor_procesamiento/Procesador/Models"

/*
Funciòn que permite validar si un host tiene los recursos (CPU y RAM) que se estàn solicitando
@cpuRequerida Paràmetro que representa al cantidad de CPU requerida en el host
@ramRequerida Paràmetro que representa la cantidad de memoria RAM requerdida en el host
@host Paràmetro que representa el host en el cual se quiere realizar la validaciòn
@Return Retorna true en caso de que el host tenga libre los recursos solicitados, o false, en caso contrario
*/
func ValidateHostResourceAvailability(cpuRequerida int, ramRequerida int, host models.Host) bool {

	recursosDisponibles := false

	var cpuNecesitada int
	cpuDisponible := float64(host.Cpu_total) * 0.75 //Obtiene el 75% de la cpu total del host

	if cpuRequerida != 0 {
		cpuNecesitada = cpuRequerida + host.Cpu_usada
	}

	var ramNecesitada int
	ramDisponible := float64(host.Ram_total) * 0.75 //Obtiene el 75% de la ram total del host

	if ramRequerida != 0 {
		ramNecesitada = ramRequerida + host.Ram_usada
	}

	if cpuNecesitada != 0 && cpuNecesitada < int(cpuDisponible) {
		recursosDisponibles = true
	}
	if ramNecesitada != 0 && ramNecesitada < int(ramDisponible) {
		recursosDisponibles = true
	}
	return recursosDisponibles
}
