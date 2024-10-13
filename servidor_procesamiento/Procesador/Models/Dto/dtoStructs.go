package dto

// VMSpecificationsDTO representa las especificaciones de la máquina virtual
type VMSpecificationsDTO struct {
	VMName      string `json:"vm_name"`
	VMOS        string `json:"vm_so"`
	VMOSDistro  string `json:"vm_so_distro"`
	VMRAM       int    `json:"vm_ram"`
	VMCPU       int    `json:"vm_cpu"`
	VMUserEmail string `json:"vm_usr_email"`
	VMHostname  string `json:"vm_hostname"`
}

// CreateVMRequestDTO representa el cuerpo completo de la solicitud
type CreateVMRequestDTO struct {
	Specifications VMSpecificationsDTO `json:"specifications"`
	ClientIP       string              `json:"clientIP"`
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
