package utilities

import (
	"fmt"
	"log"
	"regexp"
	models "servidor_procesamiento/Procesador/Models/Entities"
)

// Función que crea una imagen docker desde dockerhub dentro de una maquina virtual
// params: imagen, version, ip, hostname
// returns: string (mensaje de confirmación)
func CreateImageDockerHub(imagen, version, ip, hostname string) string {

	sctlCommand := "docker pull " + imagen + ":" + version

	fmt.Println(hostname)

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	respuesta, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return respuesta
}

//-------------------------------------------------------------------
//Desde aqui empieza el codigo para el funcionamiento de Docker UQ
//-------------------------------------------------------------------

// Función que crea una imagen docker desde un archivo tar dentro de una maquina virtual
// params: nombreArchivo, ip, hostname
// returns: string (mensaje de confirmación)
func CreateImageTarFile(nombreArchivo, ip, hostname string) string {

	sctlCommand := "docker load < " + nombreArchivo

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return "Comando Envidado con exito"
}

// Función que crea una imagen docker desde un archivo Dockerfile dentro de una maquina virtual
// params: nombreArchivo, nombreImagen, ip, hostname
// returns: string (mensaje de confirmación)
func CreateImageDockerFile(nombreArchivo, nombreImagen, ip, hostname string) string {

	sctlCommand := "mkdir /home/" + hostname + "/" + nombreImagen + "&&" + " unzip " + nombreArchivo + " -d /home/" + hostname + "/" + nombreImagen

	fmt.Println(hostname)

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	fmt.Println("dockerFile")

	fmt.Println(nombreArchivo)

	sctlCommand = "cd /home/" + hostname + "/" + nombreImagen + "&&" + " docker build -t " + nombreImagen + " ."

	fmt.Println(hostname)

	respuesta, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return respuesta
}

// Función que lista las imagenes docker dentro de una maquina virtual
// params: ip, hostname
// returns: lista de imagenes encontradas
func ListImages(ip, hostname string) ([]models.Imagen, error) {

	fmt.Println("Revisar Imagenes:", ip, hostname)

	sctlCommand := "docker images --format " + "{{.Repository}},{{.Tag}},{{.ID}},{{.CreatedAt}},{{.Size}}"

	config, err := ConfigureSSHPassword(hostname)

	fmt.Println("hostname:", hostname)

	if err != nil {
		log.Println("Fallo en la ejecucion", err)
		return nil, err
	}

	lista, err3 := SendSSHCommand(ip, sctlCommand, config)

	fmt.Println("Ip:", ip)

	if err3 != nil {
		log.Println("Fallo en la ejecucion", err)
		return nil, err
	}

	res := SplitWord(lista)

	tabla := 0
	datos := make([]string, 5)
	var imagenes []models.Imagen
	maquinaVM := ip + " - " + hostname

	for i := 0; i < len(res); i++ {
		if tabla == 4 {
			datos[tabla] = res[i]
			imagenes = append(imagenes, RegisterImageData(datos, maquinaVM))
			tabla = 0
			datos = make([]string, 5)
		} else {
			datos[tabla] = res[i]
			tabla++
		}

	}

	return imagenes, nil

}

// Función que crea un contenedor dentro de una maquina virtual
// params: imagen, comando, ip, hostname
// returns: string (mensaje de confirmación)
func CreateContainer(imagen, comando, ip, hostname string) string {

	sctlCommand := comando + " " + imagen

	fmt.Println("\n" + sctlCommand)

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return "Comando Enviado con Exito"

}

// Función que lista los contenedores dentro de una maquina virtual
// params: ip, hostname
// returns: lista de contenedores encontrados
func ListContainers(ip, hostname string) ([]models.Contenedor, error) {

	fmt.Println("Revisar Contenedores")

	sctlCommand := "docker ps -a --format  '{{.ID}},{{.Image}},{{.Command}},{{.CreatedAt}},{{.Status}},{{if .Ports}}{{.Ports}}{{else}}No ports exposed{{end}},{{.Names}}'"

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return nil, err
	}

	lista, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return nil, err
	}

	res := SplitWord(lista)

	tabla := 0
	datos := make([]string, 7)
	var contenedores []models.Contenedor
	conetendor := 1
	maquinaVM := ip + " - " + hostname

	for i := 0; i < len(res); i++ {
		if tabla == 6 {
			datos[tabla] = res[i]
			contenedores = append(contenedores, RegisterContainerData(datos, maquinaVM))
			tabla = 0
			conetendor++
			datos = make([]string, 7)
		} else {
			datos[tabla] = res[i]
			tabla++
		}
	}

	return contenedores, err

}

// Función que corre un contenedor dentro de una maquina virtual
// params: contenedor, ip, hostname
// returns: string (mensaje de confirmación)
func RunContainer(contenedor, ip, hostname string) string {

	fmt.Println("Correr Contenedor")

	sctlCommand := "docker start " + contenedor

	fmt.Println("\n" + sctlCommand)

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return "Comando Enviado con Exito"
}

// Función que detiene un contenedor dentro de una maquina virtual
// params: contenedor, ip, hostname
// returns: string (mensaje de confirmación)
func StopContainer(contenedor, ip, hostname string) string {

	fmt.Println("Detener Contenedor")

	sctlCommand := "docker stop " + contenedor

	fmt.Println("\n" + sctlCommand)

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return "Comando Enviado con Exito"
}

// Función que reinicia un contenedor dentro de una maquina virtual
// params: contenedor, ip, hostname
// returns: string (mensaje de confirmación)
func RestartContainer(contenedor, ip, hostname string) string {

	fmt.Println("Reiniciar Contenedor")

	sctlCommand := "docker restart " + contenedor

	fmt.Println("\n" + sctlCommand)

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return "Comando Enviado con Exito"
}

// Función que elimina un contenedor dentro de una maquina virtual
// params: contenedor, ip, hostname
// returns: string (mensaje de confirmación)
func DeleteContainer(contenedor, ip, hostname string) string {

	fmt.Println("Eliminar Contenedor")

	sctlCommand := "docker rm " + contenedor

	fmt.Println("\n" + sctlCommand)

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	return "Comando Enviado con Exito"
}

// Función que elimina todos los contenedores dentro de una maquina virtual
// params: ip, hostname
// returns: string (mensaje de confirmación)
func DeleteAllContainers(ip, hostname string) string {

	sctlCommand := "docker rm $(docker ps -a -q)"

	config, err := ConfigureSSHPassword(hostname)

	if err != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}

	_, err3 := SendSSHCommand(ip, sctlCommand, config)

	if err3 != nil {
		log.Println("Error al configurar SSH:", err)
		return "Error al configurar la conexiòn SSH"
	}
	return "Comando Enviado con Exito"
}

//Funciones complementarias

// Función que divide una cadena de texto en un arreglo de strings
// params: word
// returns: arreglo de strings
func SplitWord(word string) []string {
	array := regexp.MustCompile("[,,\n]+").Split(word, -1)
	return array
}

// Función que divide una cadena de texto en un arreglo de strings para la creación de un nuevo registro
// de una imagen docker en la base de datos
// params: datos, maquinaVM
// returns: modelo de imagen
func RegisterImageData(datos []string, maquinaVM string) models.Imagen {

	nuevaImagen := models.Imagen{
		Repositorio: datos[0],
		Tag:         datos[1],
		ImagenId:    datos[2],
		Creacion:    datos[3],
		Tamanio:     datos[4],
		MaquinaVM:   maquinaVM,
	}

	return nuevaImagen

}

// Función que divide una cadena de texto en un arreglo de strings para la creación de un nuevo registro
// de un contenedor docker en la base de datos
// params: datos, maquinaVM
// returns: modelo de contenedor
func RegisterContainerData(datos []string, maquinaVM string) models.Contenedor {

	nuevaContenedor := models.Contenedor{
		ContenedorId: datos[0],
		Imagen:       datos[1],
		Comando:      datos[2],
		Creado:       datos[3],
		Status:       datos[4],
		Puerto:       datos[5],
		Nombre:       datos[6],
		MaquinaVM:    maquinaVM,
	}

	return nuevaContenedor

}
