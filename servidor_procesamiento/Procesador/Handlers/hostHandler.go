package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models"
	utilities "servidor_procesamiento/Procesador/Utilities"
	"strconv"
)

/*
Clase encargada de contener los handlers que responden a los endpoints que
atienden las consultas sobre estos
*/

// Funcion que responde al endpoint encargado de consultar las maquinas virtuales
func ConsultMachineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	hosts, err := database.ConsultHosts()
	if err != nil && err.Error() != "no Hosts encontrados" {
		fmt.Println(err)
		log.Println("Error al consultar los Host")
		return
	}

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hosts)

}

// Funcion que responde al endpoint encargado de verificar los diferentes host registrados en la base de datos
func CheckHostHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica que la solicitud sea del método POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica el JSON recibido en la solicitud en una estructura Specifications.
	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Error al decodificar JSON de la solicitud", http.StatusBadRequest)
		return
	}

	//Se capturan los datos de la maquina virtual a crear
	mv := payload["specifications"].(map[string]interface{})

	/*En la base de datos los indices de los host empiezan desde el indice 1
	si el valor es cero se utiliza para disparar el algoritmo aleatorio
	*/
	id := int(mv["Host_id"].(float64))
	switch {
	case id == 0:
		//Se encola la maquina virtual a crear
		config.GetMu().Lock()
		config.GetMaquina_virtualQueue().Queue.PushBack(payload)
		config.GetMu().Unlock()

		//Se imprime el estado actual de la cola
		fmt.Println("Cantidad de Solicitudes de Especificaciones en la Cola: " + strconv.Itoa(config.GetMaquina_virtualQueue().Queue.Len()))

		// Envía una respuesta al servidor web
		response := map[string]string{"mensaje": "Mensaje JSON de crear MV recibido correctamente", "centinela": "true"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	case id > 0:
		
		mihost, _ := database.GetHost(int(mv["Host_id"].(float64)))
		estadossh := utilities.Marcapasos(config.GetPrivateKeyPath(), mihost.Hostname, mihost.Ip)
		if estadossh {
			//Se encola la maquina virtual a crear
			config.GetMu().Lock()
			config.GetMaquina_virtualQueue().Queue.PushBack(payload)
			config.GetMu().Unlock()

			//Se imprime el estado actual de la cola
			fmt.Println("Cantidad de Solicitudes de Especificaciones en la Cola: " + strconv.Itoa(config.GetMaquina_virtualQueue().Queue.Len()))

			// Envía una respuesta al servidor web
			response := map[string]string{"mensaje": "Mensaje JSON de crear MV recibido correctamente", "centinela": "true"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			// Envía una respuesta al servidor web
			response := map[string]string{"mensaje": "Esta maquina tiene problemas :(", "centinela": "false", "hostmalo:": strconv.Itoa(id)}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}

	case id < 0:
		// Envía una respuesta al cliente.
		response := map[string]string{"mensaje": "Error de seleccion de maquina", "centinela": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

// Funcion que responde al endpoint encargado de consultar los host registrados en la base de datos
func ConsultHostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	var persona models.Persona
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&persona); err != nil { //Solo llega el email
		http.Error(w, "Error al decodificar JSON de inicio de sesión", http.StatusBadRequest)
		return
	}

	persona, error := database.GetUser(persona.Email)
	if error != nil {
		return
	}

	hosts, err := database.ConsultHosts()
	if err != nil && err.Error() != "no Hosts encontrados" {
		fmt.Println(err)
		log.Println("Error al consultar los Host")
		return
	}

	// Respondemos con la lista de máquinas virtuales en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hosts)

}

// Funcion que responde al endpoint encargado de agregar un host a la base de datos
func AddHostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se requiere una solicitud POST", http.StatusMethodNotAllowed)
		return
	}

	var host models.Host

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&host); err != nil {
		http.Error(w, "Error al decodificar JSON de especificaciones", http.StatusBadRequest)
		return
	}

	query := "insert into host (nombre, mac, ip, hostname, ram_total, cpu_total, almacenamiento_total, adaptador_red, estado, ruta_llave_ssh_pub, sistema_operativo, distribucion_sistema_operativo) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	//Registra el usuario en la base de datos
	_, err := database.DB.Exec(query, host.Nombre, host.Mac, host.Ip, host.Hostname, host.Ram_total, host.Cpu_total, host.Almacenamiento_total, host.Adaptador_red, "Activo", host.Ruta_llave_ssh_pub, host.Sistema_operativo, host.Distribucion_sistema_operativo)
	if err != nil {
		fmt.Println("Error al registrar el host.")

	} else if err != nil {
		panic(err.Error())
	}

	fmt.Println("Registro del host exitoso")
	response := map[string]bool{"registroCorrecto": true}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}