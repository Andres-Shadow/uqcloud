package Models

import "time"

// Clase que gestiona las maquinas virtuales que se crean
type VirtualMachine struct {
	Id                string    `json:"vm_uuid"`
	Name              string    `json:"vm_name"`
	Ram               int       `json:"vm_ram"`
	Cpu               int       `json:"vm_cpu"`
	Ip                string    `json:"vm_ip"`
	Estado            string    `json:"vm_state"`
	Hostname          string    `json:"vm_hostname"`
	Person_Email      string    `json:"vm_usr_email"`
	Host_id           int       `json:"vm_host_id"`
	Disco_id          int       `json:"vm_disk_id"`
	Sistema_operativo string    `json:"vm_so"`
	Distribucion_SO   string    `json:"vm_so_distro"`
	Fecha_creacion    time.Time `json:"vm_creation_date"`
}
