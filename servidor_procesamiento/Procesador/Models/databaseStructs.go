package models

import (
	"time"

	"gorm.io/gorm"
)

/*
Clase encarga de contener todas las estructuras que posteriormente
representaran entidades en la base de datos
*/

/*
Estrucutura de datos tipo JSON que contiene los campos necesarios para la gestiòn de usuarios
@Nombre Representa el nombre del usuario
@Apellido Representa el apellido del usuario
@Email Representa el email del usuario
@Contrasenia Representa la contraseña de la cuenta
@Rol Representa el rol que tiene la persona en la plataforma. Puede ser Estudiante o Administrador
*/
type Persona struct {
	gorm.Model
	Nombre      string `json:"nombre"`
	Apellido    string `json:"apellido"`
	Email       string `json:"email"`
	Contrasenia string `json:"contrasenia"`
	Rol         byte   `json:"rol"`
}

/*
Estructura de datos tipo JSOn que contiene los datos de una màquina virtual
@Uuid Representa el uuid de una màqina virtual, el cual es un identificador ùnico
@Nombre Representa el nombre de la MV
@Ram Representa la cantidad de memoria RAM que tiene la màquina virtual
@Cpu Representa la cantidad de unidades de procesamiento que tiene la màquina virtial
@Ip Representa la direcciòn IP de la màquina
@Estado Representa el estado actual de la MV. Puede ser: Encendido, Apagado ò Procesando. Este ùltimo estado indica que la màquina se està encendiendo o apagando
@Hostname Representa el nombre del usuario del sistema operativo
@Persona_email Representa el email de la persona asociada a la MV.
@Host_id Representa el identificador ùnico de la màquina host en la cual està creada la MV
@Disco_id Representa el identificador ùnico del disco al cual està conectada la MV
@Sistema_operativo Represneta el tipo de sistema operativo que tiene la MV. Por ejemplo: Linux o Windows
@Distribucion_sistema_operativo Representa la distribuciòn del sistema operativo que està usando la MV. Por ejemplo: Debian ò 11 Home
*/
type Maquina_virtual struct {
	gorm.Model
	Uuid                           string
	Nombre                         string
	Ram                            int
	Cpu                            int
	Ip                             string
	Estado                         string
	Hostname                       string
	Persona_email                  string
	Host_id                        int
	Disco_id                       int
	Sistema_operativo              string
	Distribucion_sistema_operativo string
	Fecha_creacion                 time.Time
}

/*
Estructura de datos tipo JSON que contiene los campos de un host
@Id Representa el identificador ùnico del host
@Nombre Representa el nombre del host
@Mac Representa la direcciòn fìsica del host
@Ip Representa la direcciòn Ip del host
@Hostname Representa el nombre del host
@Ram_total Representa la cantidad total de memoria RAM que tiene el host. Se representa en mb
@Cpu_total Representa la cantidad de unidades de procesamiento total que tiene el host
@Almacentamiento_total Representa la cantidad total de almacenamiento del host. Se representa en mb
@Ram_usada Representa la cantidad total de memoria RAM que està siendo usada por las màquinas virtuales alojadas en el host. Se representa en mb
@Cpu_usada Representa la cantidad total de unidades de procesamiento que estàn siendo usadas por las MV's alojadas en el host
@Almacenamiento_usado Representa la cantidad de alamacenamiento que està siendo usado por las MV's alojadas en el host. Se representa en mb
@Adaptador_red Representa el nombre del adaptador de red del host
@Estado Representa el estado del host (Disponible o Fuera de servicio)
@Ruta_llave_ssh_pub Representa la ubiaciòn de la llave ssh pùblica
@Sistema_operativo Representa el tipo de sistema operativo del host. Por ejemplo: Windows o Mac
@Distribucion_sistema_operativo Representa el tipo de distribuciòn del sistema operativo que tiene el host. Por ejemplo: 10 Pro o 11 Home
*/
type Host struct {
	gorm.Model
	Id                             int
	Nombre                         string
	Mac                            string
	Ip                             string
	Hostname                       string
	Ram_total                      int
	Cpu_total                      int
	Almacenamiento_total           int
	Ram_usada                      int
	Cpu_usada                      int
	Almacenamiento_usado           int
	Adaptador_red                  string
	Estado                         string
	Ruta_llave_ssh_pub             string
	Sistema_operativo              string
	Distribucion_sistema_operativo string
}

/*
Estructura de datos tipo JSON que contiene los campos para representar una MV del catàlogo
@Nombre Representa el nombre de la MV
@Memoria Representa la cantidad de memoria RAM de la MV
@Cpu Representa la cantidad de unidades de procesamiento de la MV
@Sistema_operativo Representa el tipo de sistema operativo de la Mv
@Distribucion_sistema_operativo Representa la distribuciòn del sistema operativo que tiene la màquina del catàlogo
@Arquitectura Respresenta la arquitectura del sistema operativo. Se presententa en un valor entero. Por ejemplo: 32 o 64
*/
type Catalogo struct {
	gorm.Model
	Id                             int
	Nombre                         string
	Ram                            int
	Cpu                            int
	Sistema_operativo              string
	Distribucion_sistema_operativo string
	Arquitectura                   int
}

/*
Estructura de datos tipo JSON que representa la informaciòn de los discos que tiene la plataforma Desktop Cloud
@Id Representa el identificador ùnico del disco en la base de datos. Este identificador es generado automaticamente por la base de datos
@Nombre Representa el nombre del disco
@Ruta_ubicacion Representa la ubicaciòn de disco en el host.
@Sistema_operativo Representa el tipo de sistema operativo que tiene el disco. Por ejemplo: Linux
@Distribucion_sistema_operativo Representa el tipo de distribuciòn del sistema operativo. Por ejemplo: Debian o Ubuntu
@arquitectura Representa la arquitectura del sistema operativo. Se representa en un valor entero. Por ejemplo: 32 o 64
@Host_id Representa el identificador ùnico del host en el cual està ubicado el disco
*/
type Disco struct {
	gorm.Model
	Id                             int
	Nombre                         string
	Ruta_ubicacion                 string
	Sistema_operativo              string
	Distribucion_sistema_operativo string
	Arquitectura                   int
	Host_id                        int
}

/*
Estructura de datos tipo JSON que representa la informaciòn de las imagenes que tiene la plataforma Desktop Cloud
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
Estructura de datos tipo JSON que representa la informaciòn de los contenedores que tiene la plataforma Desktop Cloud
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
	ConetendorId string
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
