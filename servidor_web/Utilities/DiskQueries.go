package Utilities

import (
	"AppWeb/Config"
	"AppWeb/DTO"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetDiskNamesFromServer() (DTO.DiskResponseDTO, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.DISK_NAMES_URL)
	log.Println(serverURL)
	var diskResponse DTO.DiskResponseDTO

	resp, err := http.Get(serverURL)
	if err != nil {
		return diskResponse, fmt.Errorf("Error al realizar la solicutud: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diskResponse, fmt.Errorf("Error al obtener la solicitud")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return diskResponse, fmt.Errorf("Error al obtener la solicitud")
	}

	if err := json.Unmarshal(body, &diskResponse); err != nil {
		return diskResponse, fmt.Errorf("Error al decodificar el JSON: %w", err)
	}

	return diskResponse, nil
}

func GetHostsOfDiskFromServer(DiskName string) (DTO.HostsOfDisksResponseDTO, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s/%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.DISK_VM_URL, DiskName)
	log.Println(serverURL)
	var hostDisk DTO.HostsOfDisksResponseDTO

	resp, err := http.Get(serverURL)
	if err != nil {
		return hostDisk, fmt.Errorf("Error al realizar la solicutud: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return hostDisk, fmt.Errorf("Error al obtener la solicitud")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return hostDisk, fmt.Errorf("Error al obtener la solicitud")
	}

	if err := json.Unmarshal(body, &hostDisk); err != nil {
		return hostDisk, fmt.Errorf("Error al decodificar el JSON: %w", err)
	}

	return hostDisk, nil

}

func DeleteHostOfDiskFromServer(DiskName string, HostID int) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s/%s?host_id=%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.DISK_VM_URL, DiskName, HostID)
	log.Println(serverURL)

	req, err := http.NewRequest("DELETE", serverURL, bytes.NewBuffer(nil))
	if err != nil {
		return false, fmt.Errorf("Error al realizar la solicutud: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("Error al realizar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Error al obtener la respuesta: %w", err)
	}

	return true, nil
}
