package models

import "gorm.io/gorm"

/*
Estructura de datos tipo JSON que representa la informaciòn de las imagenes que tiene la plataforma QuickCloud
@Repositorio Representa el identificador ùnico del disco en la base de datos. Este identificador es generado automaticamente por la base de datos
@Tag Representa el nombre del disco
@ImagenId Representa la ubicaciòn de disco en el host.
@Creacion Representa el tipo de sistema operativo que tiene el disco. Por ejemplo: Linux
@Tamanio Representa el tipo de distribuciòn del sistema operativo. Por ejemplo: Debian o Ubuntu
*/
type Imagen struct {
	gorm.Model
	Repositorio string
	Tag         string
	ImagenId    string
	Creacion    string
	Tamanio     string
	MaquinaVM   string
}

/*
Estructura de datos tipo JSON que representa la informaciòn de los contenedores que tiene la plataforma QuickCloud
@ConetendorId Representa el identificador ùnico del disco en la base de datos. Este identificador es generado automaticamente por la base de datos
@Imagen Representa el nombre del disco
@Comando Representa la ubicaciòn de disco en el host.
@Creado Representa el tipo de sistema operativo que tiene el disco. Por ejemplo: Linux
@Status Representa el tipo de distribuciòn del sistema operativo. Por ejemplo: Debian o Ubuntu
@Puerto Representa la arquitectura del sistema operativo. Se representa en un valor entero. Por ejemplo: 32 o 64
@Nombre Representa el identificador ùnico del host en el cual està ubicado el disco
*/

type Contenedor struct {
	gorm.Model
	ContenedorId string
	Imagen       string
	Comando      string
	Creado       string
	Status       string
	Puerto       string
	Nombre       string
	MaquinaVM    string
}

type CatalogoDisco struct {
	CatalogoID int
	DiscoID    int
}
