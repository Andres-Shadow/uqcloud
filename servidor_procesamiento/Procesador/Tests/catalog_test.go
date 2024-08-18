package tests

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsultCatalog(t *testing.T) {
	resp, err := http.Get(RootEndpointURL + "catalog")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "Error al consultar el catálogo: el catalogo se encuentra vacío\n", ReadBody(resp))
}

// readBody lee el cuerpo de la respuesta HTTP y lo devuelve como una cadena.
func ReadBody(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
