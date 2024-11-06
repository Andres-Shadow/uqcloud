package DTO

type HostDTO struct {
	AdaptadorRed        string `json:"adaptador_red"`
	AlmacenamientoTotal int    `json:"almacenamiento_total"`
	CPUTotal            int    `json:"cpu_total"`
	Estado              string `json:"estado"`
	UserName            string `json:"hostname"`
	HostName            string `json:"hst_name"`
	ID                  int    `json:"id"`
	IP                  string `json:"ip"`
	RAMTotal            int    `json:"ram_total"`
	SistemaOperativo    string `json:"sistema_operativo"`
}

type HostsResponseDTO struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Count   int       `json:"count"`
	Data    []HostDTO `json:"data"`
}

type HostIDDTO struct {
	HostIds []int `json:"hostIds"`
}
