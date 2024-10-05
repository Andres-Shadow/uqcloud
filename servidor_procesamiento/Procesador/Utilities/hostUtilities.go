package utilities

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
	"strconv"
	"strings"
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

func GetHostWithMostResources() (models.Host, error) {
	var registeredHosts []models.Host = []models.Host{}
	var selectedHost models.Host

	// Obtener todos los hosts registrados en la base de datos
	registeredHosts = database.GetHosts()

	if len(registeredHosts) == 0 {
		return selectedHost, nil
	}

	// Recorrer todos los host disponibles para verificar cual tiene
	// la mayor cantidad de ram disponible para la creación de máquinas virtuales

	var maxRam float64 = 0

	for _, host := range registeredHosts {
		//verficar que el host este vivo con el marcapasos
		if Pacemaker(config.GetPrivateKeyPath(), host.Hostname, host.Ip) {
			//hace una llamada http a la ip de los host:9182/metrics
			//para obtener la cantidad de ram usada
			request, err := http.NewRequest("GET", "http://"+host.Ip+":9182/metrics", nil)
			if err != nil {
				log.Printf("Error al crear la petición HTTP: %v", err)
				return selectedHost, err
			}
			client := &http.Client{}
			response, err := client.Do(request)

			if err != nil {
				log.Printf("Error al realizar la petición HTTP: %v", err)
				return selectedHost, err
			}

			defer response.Body.Close()

			// Leer la respuesta línea por línea
			reader := bufio.NewReader(response.Body)
			var totalRAM, availableRAM float64
			for {
				//se busca la cantidad de ram usada en el cuerpo de la respuesta
				// etiqueta con la cantidad total de ram  -> windows_cs_physical_memory_bytes
				// etiqueta con la cantidad de ram usada -> windows_os_physical_memory_free_bytes

				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Printf("Error al leer el cuerpo de la respuesta: %v", err)
					return selectedHost, err
				}

				// Buscar la línea que contiene la cantidad total de RAM
				if strings.Contains(line, "windows_cs_physical_memory_bytes") {
					totalRAM = extractValue(line)
				}

				// Buscar la línea que contiene la RAM disponible
				if strings.Contains(line, "windows_os_physical_memory_free_bytes") {
					availableRAM = extractValue(line)
				}
			}

			// Calcular la RAM utilizada
			usedRAM := totalRAM - availableRAM
			// fmt.Printf("Total RAM: %.0f bytes\n", totalRAM)
			// fmt.Printf("Available RAM: %.0f bytes\n", availableRAM)
			// fmt.Printf("Used RAM: %.0f bytes\n", usedRAM)
			if usedRAM > maxRam {
				maxRam = usedRAM
				selectedHost = host
			}
		}
	}

	return selectedHost, nil
}

// Función para extraer el valor numérico de una línea
func extractValue(line string) float64 {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0
	}
	value, err := strconv.ParseFloat(parts[len(parts)-1], 64)
	if err != nil {
		log.Printf("Error al convertir el valor: %v", err)
		return 0
	}
	return value
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

	// Procesar cada archivo JSON
	for _, file := range files {
		// Leer el archivo JSON
		data, err := os.ReadFile(file)
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
				SetUpHostAndDisk(host)
			}
		}()
	}
}

func SetUpHostAndDisk(host models.Host) {

	var diskNames []string = []string{"Alpine", "fedora", "Ubuntu", "Debian"}
	var auxDisk models.Disco = models.Disco{}
	const diskPath = "c:\\uqcloud"

	// Guardar el host en la base de datos usando GORM
	if err := database.DATABASE.Create(&host).Error; err != nil {
		log.Printf("Error al guardar el host en la base de datos: %v", err)
	} else {
		log.Printf("Host pre-registrado exitosamente: %s", host.Nombre)
		auxDisk.Arquitectura = 64
		auxDisk.Sistema_operativo = "linux"
		auxDisk.Host_id = host.Id

		for _, diskName := range diskNames {
			auxDisk.Nombre = diskName
			auxDisk.Distribucion_sistema_operativo = "" + diskName + "-" + strconv.Itoa(auxDisk.Arquitectura)
			auxDisk.Ruta_ubicacion = diskPath + "\\" + diskName + ".vdi"
			err := database.CreateDisck(auxDisk)
			if err != nil {
				log.Printf("Error al pre-registrar el disco en el host: %v", err)
			} else {
				log.Printf("Disco pre-registrado exitosamente: %s", auxDisk.Nombre)
			}
		}
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
