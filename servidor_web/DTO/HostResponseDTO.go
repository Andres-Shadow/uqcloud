package DTO

type HostDTO struct {
	AlmacenamientoTotal int    `json:"almacenamiento_total"`
	CPUTotal            int    `json:"cpu_total"`
	Estado              string `json:"estado"`
	HostName            string `json:"hst_name"`
	ID                  int    `json:"id"`
	IP                  string `json:"ip"`
	RAMTotal            int    `json:"ram_total"`
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
