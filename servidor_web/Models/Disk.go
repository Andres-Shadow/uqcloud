package Models

type Disk struct {
	Name              string `json:"dsk_name"`
	Ruta_Ubicacion    string `json:"dsk_route"`
	Sistema_Operativo string `json:"dsk_so"`
	Distrubucion_SO   string `json:"dsk_so_distro"`
	Arquitectura      int    `json:"dsk_arch"`
	Host_id           int    `json:"dsk_host_id"`
}
