package Models

// Clase encargada de manejar los host en donde se crean las Maquina virtuales
type Host struct {
	Id                   int    `json:"id"`
	Name                 string `json:"hst_name"`
	Mac                  string `json:"hst_mac"`
	Ip                   string `json:"hst_ip"`
	Hostname             string `json:"hst_hostname"`
	Ram_total            int    `json:"hst_ram"`
	Cpu_total            int    `json:"hst_cpu"`
	Almacenamiento_total int    `json:"hst_storage"`
	Adaptador_red        string `json:"hst_network"`
	Estado               string `json:"hst_state"`
	Ruta_llave_ssh_pub   string `json:"hst_sshroute"`
	Sistema_operativo    string `json:"hst_so"`
	Distribucion_SO      string `json:"hst_so_distro"`
}
