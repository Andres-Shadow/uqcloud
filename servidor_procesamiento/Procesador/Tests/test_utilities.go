package tests

import (
	"io/ioutil"
	"log"
	"net/http"
)

// readBody lee el cuerpo de la respuesta HTTP y lo devuelve como una cadena.
func ReadBody(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
