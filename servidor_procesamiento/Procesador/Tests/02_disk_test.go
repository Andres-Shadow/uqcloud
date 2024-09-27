package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDiskOK(t *testing.T) {
	resp, err := http.Post(RootEndpointURL+"disk", "application/json", strings.NewReader(`{
		"dsk_name": "Disco Alpine TEST",
		"dsk_route" : "C:\\uqcloud\\Alpine.vdi",
		"dsk_so" : "Linux",
		"dsk_so_distro" : "alpine x64",
		"dsk_arch": 64,
		"dsk_host_id" : 16
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Función que verifica que responda un error 400 cuando no se mandan todos los campos obligatorios en el JSON
func TestAddDiskBadStructJson(t *testing.T) {
	// JSON con campos faltantes
	jsonBody := `{
		"dsk_name": "Disco 1",
		"dsk_route": "/path/to/disco"
	}`
	resp, err := http.Post(RootEndpointURL+"disk", "application/json", strings.NewReader(jsonBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
