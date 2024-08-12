package Models

/*
Clase que guarda la imagen de los docker
*/
type Imagen struct {
	IdImagen    string `json:"id_imagen"`
	Repositorio string `json:"repositorio"`
	Tag         string `json:"tag"`
	Creacion    string `json:"creacion"`
	Tamanio     string `json:"tamanio"`
	MaquinaVM   string `json:"maquina_vm"`
}
