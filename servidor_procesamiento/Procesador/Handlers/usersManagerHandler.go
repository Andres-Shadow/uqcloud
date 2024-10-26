package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models/Entities"
	userutilities "servidor_procesamiento/Procesador/Utilities/UserUtilities"

	"golang.org/x/crypto/bcrypt"
)

/*
Clase encargada de contener los handlers que responden a las acciones a realizar con la
gestión de los usuarios
*/

// Funcion que responde al endpoint encargado de iniciar sesión a un usuario
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {

	var persona models.Persona

	// Log para ver el contenido del cuerpo
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Cuerpo de la solicitud:", string(body))

	// Restaurar el cuerpo para que json.NewDecoder pueda procesarlo
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&persona); err != nil {
		http.Error(w, "Error al decodificar JSON de inicio de sesión", http.StatusBadRequest)
		return
	}

	resultPassword, err := database.GetUserPassword(persona.Email)

	if err != nil {
		fmt.Println("Error al buscar la contraseña del usuario "+persona.Email+" en la base de datos:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	verificacion := bcrypt.CompareHashAndPassword([]byte(resultPassword), []byte(persona.Contrasenia))

	if verificacion != nil {
		fmt.Println("Contraseña incorrecta, no se concuerda con el registro en la base de datos")
		response := map[string]interface{}{
			"status":  false,
			"usuario": nil,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
	} else {
		fmt.Println("Contraseña correcta, se encontró el registro en la base de datos")
		// Consulta en la base de datos para obtener los detalles del usuario

		usuario, err := database.GetUserFromEmail(persona.Email)
		if err != nil {
			fmt.Println("Error al obtener detalles del usuario al intentar hacer login: ", err)
			response := map[string]interface{}{
				"status":  false,
				"usuario": nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := map[string]interface{}{
			"status":  true,
			"usuario": usuario,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Funcion que responde al endpoint encargado de crear maquinas virtuales para invitados
func CreateTempUserHandler(w http.ResponseWriter, r *http.Request) {
	// retorna el correo temporal
	email := userutilities.CreateTempAccount()

	// Envía una respuesta al cliente.
	response := map[string]string{"mensaje": email}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
