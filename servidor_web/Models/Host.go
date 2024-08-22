package Models

// Clase encargada de manejar los host en donde se crean las Maquina virtuales
type Host struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	Mac                  string `json:"mac"`
	Ip                   string `json:"ip"`
	Hostname             string `json:"hostname"`
	Ram_total            int    `json:"ram_total"`
	Cpu_total            int    `json:"cpu_total"`
	Almacenamiento_total int    `json:"almacenamiento_total"`
	Adaptador_red        string `json:"adaptador_red"`
	Estado               string `json:"estado"`
	Ruta_llave_ssh_pub   string `json:"ruta_llave_ssh_pub"`
	Sistema_operativo    string `json:"sistema_operativo"`
	Distribucion_SO      string `json:"distribucion_SO"`
}
