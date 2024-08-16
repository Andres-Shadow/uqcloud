package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDisk(t *testing.T) {
	resp, err := http.Post(RootEndpointURL+"disk", "application/json", strings.NewReader(`{
		"specifications": {
			"Nombre": "Nuevo Disco",
			"Capacidad": "500GB",
			"Tipo": "SSD"
		}
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
