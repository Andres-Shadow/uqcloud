package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testea el funcionamiento de la respuesta cuando el catalogo se encuentra vacio (NOTA: CATALOGO YA NO SE USA PARA NADA)
func TestConsultEmptyCatalog(t *testing.T) {
	resp, err := http.Get(RootEndpointURL + "catalog")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "Error al consultar el catálogo: el catalogo se encuentra vacío\n", ReadBody(resp))
}
