package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var RootEndpointURL = "http://localhost:8081/api/v1/"

// OJO, ESTE TEST PUEDE ROMPER TU CONTENEDOR DE DOCKER, USALO CON CUIDADO
func TestCreateVirtualMachine(t *testing.T) {
	resp, err := http.Post(RootEndpointURL+"virtual_machine", "application/json", strings.NewReader(`{
		"specifications": {
			"Nombre": "Servidor_Principal",
			"Sistema_operativo": "Linux",
			"Distribucion_sistema_operativo": "Ubuntu 22.04 LTS",
			"Ram": 8192,
			"Cpu": 4,
			"Persona_email": "admin@uqcloud.co",
			"Host_id": 2
		},
		"clientIP" : "192.168.1.2"
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetVirtualMachine(t *testing.T) {
	resp, err := http.Get(RootEndpointURL + "virtual_machine/admin@uqcloud.co")
	// assert error 400 and message "No se encontraron màquinas virtuales para el usuario\n"
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "No se encontraron màquinas virtuales para el usuario\n", ReadBody(resp))

}

func TestModifyVirtualMachine(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, RootEndpointURL+"virtual_machine", strings.NewReader(`{
		"specifications": {
			"Nombre": "Servidor_Principal",
			"Sistema_operativo": "Linux",
			"Distribucion_sistema_operativo": "Ubuntu 22.04 LTS",
			"Ram": 8192,
			"Cpu": 4,
			"Persona_email": "admin@uqcloud.co",
			"Host_id": 2
		},
		"clientIP" : "localhost"
	}`))
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteVirtualMachine(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, RootEndpointURL+"virtual_machine", strings.NewReader(`{
		"specifications": {
			"Nombre": "Servidor_Principal",
			"Sistema_operativo": "Linux",
			"Distribucion_sistema_operativo": "Ubuntu 22.04 LTS",
			"Ram": 8192,
			"Cpu": 4,
			"Persona_email": "admin@uqcloud.co",
			"Host_id": 2
		},
		"clientIP" : "localhost",
		"nombreVM" : "test"
	}`))
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
