package Models

// Clase que gestiona las maquinas virtuales que se crean
type VirtualMachine struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Ram               int    `json:"ram"`
	Cpu               int    `json:"cpu"`
	Ip                string `json:"ip"`
	Estado            string `json:"estado"`
	Hostname          string `json:"hostname"`
	Person_Email      string `json:"person_email"`
	Host_id           int    `json:"host_id"`
	Disco_id          int    `json:"disco_id"`
	Sistema_operativo string `json:"sistema_operativo"`
	Distribucion_SO   string `json:"distribucion_SO"`
}
