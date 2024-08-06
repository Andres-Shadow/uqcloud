package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
	utilities "servidor_procesamiento/Procesador/Utilities"

	"golang.org/x/crypto/bcrypt"
)

/*
Clase encargada de contener los handlers que responden a las acciones a realizar con la
gestión de los usuarios
*/

// Funcion que responde al endpoint encargado de iniciar sesión a un usuario
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}
	var persona models.Persona
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&persona); err != nil {
		http.Error(w, "Error al decodificar JSON de inicio de sesión", http.StatusBadRequest)
		return
	}

	// Si las credenciales son válidas, devuelve un JSON con "loginCorrecto" en true, de lo contrario, en false.
	query := "SELECT contrasenia FROM persona WHERE email = ?"
	var resultPassword string

	//Consulta en la base de datos si el usuario existe
	err := database.DB.QueryRow(query, persona.Email).Scan(&resultPassword)

	err2 := bcrypt.CompareHashAndPassword([]byte(resultPassword), []byte(persona.Contrasenia))
	if err2 != nil {
		fmt.Println("Contraseña incorrecta")
	} else {
		fmt.Println("Contraseña correcta")
	}
	if err2 != nil {
		response := map[string]interface{}{
			"loginCorrecto": false,
			"usuario":       nil,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
	} else if err != nil {
		panic(err.Error())
	} else {
		// Consulta en la base de datos para obtener los detalles del usuario
		queryUsuario := "SELECT * FROM persona WHERE email = ?"
		var usuario models.Persona

		errUsuario := database.DB.QueryRow(queryUsuario, persona.Email).Scan(&usuario.Email, &usuario.Nombre, &usuario.Apellido, &usuario.Contrasenia, &usuario.Rol)
		if errUsuario != nil {
			fmt.Println("Error al obtener detalles del usuario:", errUsuario)
			response := map[string]interface{}{
				"loginCorrecto": false,
				"usuario":       nil,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := map[string]interface{}{
			"loginCorrecto": true,
			"usuario":       usuario,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Funcion que responde al endpoint encargado de registrar un usuario nuevo a la base de datos
func UserSignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	var persona models.Persona
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&persona); err != nil {
		http.Error(w, "Error al decodificar JSON de inicio de sesión", http.StatusBadRequest)
		return
	}

	utilities.PrintUserAccount(persona)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(persona.Contrasenia), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error al encriptar la contraseña:", err)
		return
	}

	query := "INSERT INTO persona (nombre, apellido, email, contrasenia, rol) VALUES ( ?, ?, ?, ?, ?);"
	var resultUsername string

	//Registra el usuario en la base de datos
	_, err = database.DB.Exec(query, persona.Nombre, persona.Apellido, persona.Email, hashedPassword, "Estudiante")
	if err != nil {
		fmt.Println("Error al registrar.")
		response := map[string]bool{"loginCorrecto": false}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Registro correcto: %s\n", resultUsername)
		response := map[string]bool{"loginCorrecto": true}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}