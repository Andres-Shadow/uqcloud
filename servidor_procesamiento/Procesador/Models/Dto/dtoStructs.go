package dto

import "time"

type PeronsaDTO struct {
	Nombre      string `json:"usr_name"`
	Apellido    string `json:"usr_surname"`
	Email       string `json:"usr_email" gorm:"unique"`
	Contrasenia string `json:"usr_password"`
	Rol         byte   `json:"usr_role"`
}

type Maquina_virtualDTO struct {
	Nombre                         string    `json:"vm_name" gorm:"unique"`
	Ram                            int       `json:"vm_ram" gorm:"not null"`
	Cpu                            int       `json:"vm_cpu" gorm:"not null"`
	Ip                             string    `json:"vm_ip"`
	Estado                         string    `json:"vm_state" gorm:"not null"`
	Hostname                       string    `json:"vm_hostname" gorm:"not null"`
	Persona_email                  string    `json:"vm_usr_email" gorm:"not null"`
	Host_id                        int       `json:"vm_host_id" gorm:"not null"`
	Disco_id                       int       `json:"vm_disk_id" gorm:"not null"`
	Sistema_operativo              string    `json:"vm_so" gorm:"not null"`
	Distribucion_sistema_operativo string    `json:"vm_so_distro" gorm:"not null"`
	Fecha_creacion                 time.Time `json:"vm_creation_date" gorm:"not null"`
}

type HostDTO struct {
	Id                             int    `json:"hst_id"`
	Nombre                         string `json:"hst_name" gorm:"not null"`
	Ip                             string `json:"hst_ip" gorm:"unique, not null"`
	Hostname                       string `json:"hst_hostname" gorm:"not null"`
	Ram_total                      int    `json:"hst_ram" gorm:"not null"`
	Cpu_total                      int    `json:"hst_cpu" gorm:"not null"`
	Almacenamiento_total           int    `json:"hst_storage" gorm:"not null"`
	Ram_usada                      int    `json:"hst_used_ram"`
	Cpu_usada                      int    `json:"hst_used_cpu"`
	Almacenamiento_usado           int    `json:"hst_used_storage"`
	Adaptador_red                  string `json:"hst_network" gorm:"not null"`
	Estado                         string `json:"hst_state" gorm:"not null"`
	Sistema_operativo              string `json:"hst_so" gorm:"not null"`
	Distribucion_sistema_operativo string `json:"hst_so_distro" gorm:"not null"`
}

type DiscoDTO struct {
	Id                             int    `json:"dsk_id"`
	Nombre                         string `json:"dsk_name" gorm:"not null"`
	Ruta_ubicacion                 string `json:"dsk_route" gorm:"not null"`
	Sistema_operativo              string `json:"dsk_so" gorm:"not null"`
	Distribucion_sistema_operativo string `json:"dsk_so_distro" gorm:"not null"`
	Arquitectura                   int    `json:"dsk_arch" gorm:"not null"`
	Host_id                        int    `json:"dsk_host_id" gorm:"not null"`
}
