package Models

/*
Esta estructura temporal ayuda a manjerar las configuraciones de las maquinas virtuales de pues de ser creada, los campos de puede
Utilizar de manera que mejor convenga, es decir solo utilizar los campo que se crean necesarios deacuerdo a como lo requiera
La función
*/
type VirtualMachineTemp struct {
	Name              string `json:"vm_name"`
	Ram               int    `json:"vm_ram"`
	Cpu               int    `json:"vm_cpu"`
	Hostname          string `json:"vm_hostname"`
	Person_Email      string `json:"vm_usr_email"`
	Sistema_operativo string `json:"vm_so"`
	Distribucion_SO   string `json:"vm_so_distro"`
}

type StateMachineRequest struct {
	NombreMaquina string `json:"vm_name"`
}
