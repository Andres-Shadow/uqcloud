package apiutilities

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

func ToDTOFromDiskDistroList(disks []entities.Disco) []dto.ConsultaDiscosDTO {
	var diskDTOList []dto.ConsultaDiscosDTO
	for _, disk := range disks {
		diskDTOList = append(diskDTOList, dto.ConsultaDiscosDTO{
			Distribucion_sistema_operativo: disk.Distribucion_sistema_operativo,
		})
	}
	return diskDTOList
}

func ToDTOFromHostWithDiskList(hosts []entities.Host) []dto.ConsultaHostConDiscoDTO {
	var hostDTOList []dto.ConsultaHostConDiscoDTO
	for _, host := range hosts {
		hostDTOList = append(hostDTOList, dto.ConsultaHostConDiscoDTO{
			Id:     host.Id,
			Nombre: host.Nombre,
		})
	}
	return hostDTOList
}