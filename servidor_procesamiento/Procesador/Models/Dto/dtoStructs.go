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
	Nombre                         string `json:"hst_name" `
	Ip                             string `json:"hst_ip" `
	Hostname                       string `json:"hst_hostname" `
	Ram_total                      int    `json:"hst_ram" `
	Cpu_total                      int    `json:"hst_cpu" `
	Almacenamiento_total           int    `json:"hst_storage" `
	Ram_usada                      int    `json:"hst_used_ram"`
	Cpu_usada                      int    `json:"hst_used_cpu"`
	Almacenamiento_usado           int    `json:"hst_used_storage"`
	Adaptador_red                  string `json:"hst_network" `
	Estado                         string `json:"hst_state" `
	Sistema_operativo              string `json:"hst_so" `
	Distribucion_sistema_operativo string `json:"hst_so_distro" `
}

type ConsultaDiscosDTO struct {
	Distribucion_sistema_operativo string `json:"dsk_so_distro" `
}

type ConsultaHostConDiscoDTO struct {
	Id     int    `json:"hst_id" `
	Nombre string `json:"hst_name" `
}

type DiscoDTO struct {
	Id                             int    `json:"dsk_id"`
	Nombre                         string `json:"dsk_name" `
	Ruta_ubicacion                 string `json:"dsk_route" `
	Sistema_operativo              string `json:"dsk_so" `
	Distribucion_sistema_operativo string `json:"dsk_so_distro" `
	Arquitectura                   int    `json:"dsk_arch" `
	Host_id                        int    `json:"dsk_host_id" `
}
