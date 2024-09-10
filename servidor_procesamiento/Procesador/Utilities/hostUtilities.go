package utilities

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
)

/*
Funciòn que permite validar si un host tiene los recursos (CPU y RAM) que se estàn solicitando
@cpuRequerida Paràmetro que representa al cantidad de CPU requerida en el host
@ramRequerida Paràmetro que representa la cantidad de memoria RAM requerdida en el host
@host Paràmetro que representa el host en el cual se quiere realizar la validaciòn
@Return Retorna true en caso de que el host tenga libre los recursos solicitados, o false, en caso contrario
*/
func ValidateHostResourceAvailability(cpuRequerida int, ramRequerida int, host models.Host) bool {

	recursosDisponibles := false

	var cpuNecesitada int
	cpuDisponible := float64(host.Cpu_total) * 0.75 //Obtiene el 75% de la cpu total del host

	if cpuRequerida != 0 {
		cpuNecesitada = cpuRequerida + host.Cpu_usada
	}

	var ramNecesitada int
	ramDisponible := float64(host.Ram_total) * 0.75 //Obtiene el 75% de la ram total del host

	if ramRequerida != 0 {
		ramNecesitada = ramRequerida + host.Ram_usada
	}

	if cpuNecesitada != 0 && cpuNecesitada < int(cpuDisponible) {
		recursosDisponibles = true
	}
	if ramNecesitada != 0 && ramNecesitada < int(ramDisponible) {
		recursosDisponibles = true
	}
	return recursosDisponibles
}

/*
Función que precarga la información tomada de los hosts por medio de los jsons generados
por el script de obtenerdatoshost.ps1
*/
func PreregisterHostJsonData() {
	// Ruta donde se encuentran los archivos JSON
	jsonPath := "./DatosHostJson/"

	// Buscar todos los archivos JSON en la carpeta
	files, err := filepath.Glob(filepath.Join(jsonPath, "datosHost-*.json"))
	if err != nil {
		log.Fatalf("Error al buscar archivos JSON: %v", err)
	}

	if len(files) == 0 {
		log.Println("No se encontraron archivos JSON en la ruta especificada.")
		return
	}

	log.Printf("Archivos JSON encontrados: %v", files)

	// Procesar cada archivo JSON
	for _, file := range files {
		// Leer el archivo JSON
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Error al leer el archivo %s: %v", file, err)
			continue
		}

		// Eliminar BOM si existe
		data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

		// Deserializar el contenido del JSON en un struct Host
		var host models.Host
		err = json.Unmarshal(data, &host)
		if err != nil {
			log.Printf("Error al deserializar el archivo %s: %v", file, err)
			continue
		}

		// Crear una rutina de go para el siguiente fragmento de codigo
		go func() {
			// Verificar si el host se encuentra disponible utilizando Pacemaker mandando la ruta de la llave privada y la info del host
			if !Pacemaker(config.GetPrivateKeyPath(), host.Hostname, host.Ip) {
				log.Printf("Host no disponible: %s", host.Nombre)
			} else {
				// Guardar el host en la base de datos usando GORM
				if err := database.DATABASE.Create(&host).Error; err != nil {
					log.Printf("Error al guardar el host en la base de datos: %v", err)
				} else {
					log.Printf("Host registrado: %s", host.Nombre)
				}
			}
		}()
	}
}

/*
Función que dado el nombre de un host retorna el objeto
*/

func GetHostByName(name string) (models.Host, error) {
	var host models.Host
	err := database.DATABASE.Where("nombre = ?", name).First(&host).Error
	return host, err
}

func GetHostById(id int) (models.Host, error) {
	var host models.Host
	err := database.DATABASE.Where("id = ?", id).First(&host).Error
	return host, err
}
