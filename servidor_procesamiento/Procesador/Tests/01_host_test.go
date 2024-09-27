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

func TestAddHostOK(t *testing.T) {
	resp, err := http.Post(RootEndpointURL+"host", "application/json", strings.NewReader(`{
		"hst_name": "TestHost",
		"hst_mac": "0A-00-27-00-00-0A",
		"hst_ip":"192.168.1.5",
		"hst_hostname": "Test",
		"hst_ram" : 4,
		"hst_cpu" : 2,
		"hst_storage": 50,
		"hst_used_ram": 2,
		"hst_used_cpu": 1,  
		"hst_used_storage" : 1,
		"hst_network": "Realtek 8821CE Wireless LAN 802.11ac PCI-E NIC",
		"hst_state": "Apagado",
		"hst_sshroute": "c:\\users\\test\\.ssh\\id_rsa.pub",
		"hst_so" : "windows",
		"hst_so_distro": "64"
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestConsultHostBadRequest(t *testing.T) {
	resp, err := http.Get(RootEndpointURL + "host/noexiste")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Usuario no encontrado en la base de datos\n", ReadBody(resp))
}
