package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsultMetrics(t *testing.T) {
	resp, err := http.Get(RootEndpointURL + "metrics")
	if err != nil {
		t.Fatalf("Error al hacer la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Verifica que el estado de respuesta sea 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Lee el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error al leer el cuerpo de la respuesta: %v", err)
	}

	// Define la estructura esperada
	expectedResponse := map[string]interface{}{
		"total_CPU":                 float64(12),
		"total_CPU_usada":           float64(0),
		"total_RAM":                 float64(0),
		"total_RAM_usada":           float64(0),
		"total_estudiantes":         float64(0),
		"total_invitados":           float64(0),
		"total_maquinas_creadas":    float64(0),
		"total_maquinas_encendidas": float64(0),
		"total_usuarios":            float64(1),
	}

	// Decodifica la respuesta JSON
	var actualResponse map[string]interface{}
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Error al decodificar el JSON de la respuesta: %v", err)
	}

	// Verifica que la respuesta coincida con la esperada
	assert.Equal(t, expectedResponse, actualResponse)
}
