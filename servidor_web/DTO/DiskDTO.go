package DTO

type DiskNameDTO struct {
	NameDisk string `json:"dsk_so_distro"`
}

type DiskResponseDTO struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Count   int           `json:"count"`
	Data    []DiskNameDTO `json:"data"`
}

type HostsOfDisksResponseDTO struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    []HostInfoDTO `json:"data"`
}

type HostInfoDTO struct {
	HostID   int    `json:"hst_id"`
	HostName string `json:"hst_name"`
}
