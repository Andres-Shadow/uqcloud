package handlers

import (
	"AppWeb/Config"
	"AppWeb/DTO"
	"AppWeb/Models"
	"AppWeb/Utilities"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func CreateDiskPage(c *gin.Context) {
	session := sessions.Default(c)

	c.HTML(http.StatusOK, "createDisk.html", gin.H{
		"email":    session.Get("email").(string),
		"nombre":   session.Get("nombre").(string),
		"apellido": session.Get("apellido").(string),
		"rol":      session.Get("rol").(uint8),
		// "hosts": hosts,
	})
}

// --------- FUNCIONES NUEVA PARA LA CREACIÓN Y REGISTRO DE DISK -------- //

func CreateNewDisk(c *gin.Context) {
	//Crear el Disk a partir de la solicutud
	newDisk, err := CreateDiskFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar el JSON: " + err.Error()})
		return
	}

	//Registrar el disk
	// Definir la URL del servidor
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.DISK_VM_URL)
	log.Printf(serverURL)
	if err := Utilities.RegisterElements(serverURL, newDisk); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registro el disco"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Disco creado correctamente"})
}

/*Funcion que se encarga de decodificar los parametros para crear un nuevo disco
 */
func CreateDiskFromRequest(c *gin.Context) (Models.Disk, error) {
	var newDisk Models.Disk

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	log.Printf("Request Body: %s", string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &newDisk); err != nil {
		log.Println("Error al decodificar el JSON---: " + err.Error())
		return Models.Disk{}, err
	}

	if newDisk.Name == "" || newDisk.Ruta_Ubicacion == "" || newDisk.Sistema_Operativo == "" ||
		newDisk.Arquitectura < 0 || newDisk.Host_id <= 0 {
		log.Println("Error existen campos vacios")
		return Models.Disk{}, errors.New("error existen campos vacios")
	}

	return newDisk, nil
}

func GetDiskFromRequest(c *gin.Context) {
	diskNames, err := Utilities.GetDiskNamesFromServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se ha podido completar la consulta de los discos"})
		return
	}

	for _, diskName := range diskNames.Data {
		log.Println("DISCOS OBTENIDOS", diskName)
	}
	c.JSON(http.StatusOK, diskNames.Data)
}

func GetHostOfDiskFromRequest(c *gin.Context) {
	DiskName := c.Param("diskName")
	log.Println("DISCO A BUSCAR LOS HOSTS: ", DiskName)
	hostsOftDisk, err := Utilities.GetHostsOfDiskFromServer(DiskName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se ha podido completar la consulta"})
		return
	}
	hostsOftDiskUnique := removeDuplicate(hostsOftDisk)
	log.Println("Lista de host: ", hostsOftDiskUnique.Data)

	c.JSON(http.StatusOK, hostsOftDiskUnique.Data)
}

func DeleteHostOfDiskFrormRequest(c *gin.Context) {
	DiskName := c.Param("diskName")
	HostId, _ := strconv.Atoi(c.Query("host_id"))

	log.Println("NOMBRE DEL DISCO: ", DiskName)
	log.Println("HOST A DESASOCIAR: ", HostId)

	_, err := Utilities.DeleteHostOfDiskFromServer(DiskName, HostId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el host del disco"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted Host", "Nombre del disco": DiskName, "Host a desasociar": HostId})
}

func removeDuplicate(hostsOftDisk DTO.HostsOfDisksResponseDTO) DTO.HostsOfDisksResponseDTO {
	seen := make(map[int]bool)
	// Crear un nuevo slice para almacenar los hosts únicos
	var uniqueHosts []DTO.HostInfoDTO

	// Iterar sobre los hosts y agregar solo los que no se han visto antes
	for _, host := range hostsOftDisk.Data {
		if !seen[host.HostID] {
			// Si no hemos visto el hst_id, agregarlo a los hosts únicos
			uniqueHosts = append(uniqueHosts, host)
			// Marcar el hst_id como visto
			seen[host.HostID] = true
		}
	}

	// Actualizar la respuesta con los hosts únicos
	hostsOftDisk.Data = uniqueHosts
	return hostsOftDisk
}
