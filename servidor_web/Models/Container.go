package Models

// Clase encargada de almacenar los contenedores Docker
type Container struct {
	IdContainer string `json:"id_container"`
	Imagen      string `json:"imagen"`
	Comando     string `json:"comando"`
	Creado      string `json:"creado"`
	Status      string `json:"status"`
	Puerto      string `json:"puerto"`
	Name        string `json:"name"`
	MaquinaVM   string `json:"maquina_vm"`
}
