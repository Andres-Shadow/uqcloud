package mapper

import (
	dto "servidor_procesamiento/Procesador/Models/Dto"
	entities "servidor_procesamiento/Procesador/Models/Entities"
)

func ToHostFromDTO(dto dto.HostDTO) entities.Host {
	return entities.Host{
		Id:                             dto.Id,
		Nombre:                         dto.Nombre,
		Ip:                             dto.Ip,
		Hostname:                       dto.Hostname,
		Ram_total:                      dto.Ram_total,
		Cpu_total:                      dto.Cpu_total,
		Almacenamiento_total:           dto.Almacenamiento_total,
		Ram_usada:                      dto.Ram_usada,
		Cpu_usada:                      dto.Cpu_usada,
		Almacenamiento_usado:           dto.Almacenamiento_usado,
		Adaptador_red:                  dto.Adaptador_red,
		Estado:                         dto.Estado,
		Sistema_operativo:              dto.Sistema_operativo,
		Distribucion_sistema_operativo: dto.Distribucion_sistema_operativo,
	}
}


func ToDiscoFromDTO(dto dto.DiscoDTO) entities.Disco {
	return entities.Disco{
		Id:                             dto.Id,
		Nombre:                         dto.Nombre,
		Ruta_ubicacion:                 dto.Ruta_ubicacion,
		Sistema_operativo:              dto.Sistema_operativo,
		Distribucion_sistema_operativo: dto.Distribucion_sistema_operativo,
		Arquitectura:                   dto.Arquitectura,
		Host_id:                        dto.Host_id,
	}
}
