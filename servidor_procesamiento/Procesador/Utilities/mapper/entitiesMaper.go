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

func ToPeronsaFromDTO(dto dto.PeronsaDTO) entities.Persona {
	return entities.Persona{
		Nombre:      dto.Nombre,
		Apellido:    dto.Apellido,
		Email:       dto.Email,
		Contrasenia: dto.Contrasenia,
		Rol:         dto.Rol,
	}
}

func ToMaquina_virtualFromDTO(dto dto.Maquina_virtualDTO) entities.Maquina_virtual {
	return entities.Maquina_virtual{
		Nombre:                         dto.Nombre,
		Ram:                            dto.Ram,
		Cpu:                            dto.Cpu,
		Ip:                             dto.Ip,
		Estado:                         dto.Estado,
		Hostname:                       dto.Hostname,
		Persona_email:                  dto.Persona_email,
		Host_id:                        dto.Host_id,
		Disco_id:                       dto.Disco_id,
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


