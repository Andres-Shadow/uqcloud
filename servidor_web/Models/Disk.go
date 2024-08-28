package Models

type Disk struct {
	Name              string `json:"name"`
	Ruta_Ubicacion    string `json:"ruta_ubicacion"`
	Sistema_Operativo string `json:"sistema_operativo"`
	Distrubucion_SO   string `json:"distrubucion_so"`
	Arquitectura      int    `json:"arquitectura"`
	Host_id           int    `json:"host_id"`
}
