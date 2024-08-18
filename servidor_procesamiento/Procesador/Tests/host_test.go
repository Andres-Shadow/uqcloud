package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsultHosts(t *testing.T) {
	resp, err := http.Get(RootEndpointURL + "hosts")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// assert body is not empty
	assert.NotEmpty(t, ReadBody(resp))
}

func TestAddHost(t *testing.T) {
	resp, err := http.Post(RootEndpointURL+"host", "application/json", strings.NewReader(`{
		"Nombre": "pepe",
		"Mac": "0A-00-27-00-00-0A",
		"Ip":"192.168.1.4",
		"Hostname": "test",
		"Ram_total" : 12,
		"Cpu_total" : 2,
		"Almacenamiento_total": 2,
		"Ram_usada": 2,
		"Cpu_usada": 2,  
		"Almacenamiento_usado" : 12,
		"Adaptador_red": "test",
		"Estado": "apagado",
		"Ruta_llave_ssh_pub": "test",
		"Sistema_operativo" : "test",
		"Distribucion_sistema_operativo": "test"
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestConsultHost(t *testing.T) {
	resp, err := http.Get(RootEndpointURL + "host/noexiste")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Usuario no encontrado en la base de datos\n", ReadBody(resp))
}
