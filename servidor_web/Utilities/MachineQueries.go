package Utilities

import (
	"AppWeb/Config"
	"AppWeb/Models"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

// Consultat maquinas virtuales asociadas a un email
func ConsultMachineFromServer(email string) ([]Models.VirtualMachine, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.VIRTUAL_MACHINE_URL+"/"+email)

	// persona := Models.Person{Email: email}
	// jsonData, err := json.Marshal(persona)
	// if err != nil {
	// 	log.Println("Error al deccodificar la estrcutura peronsa como json", err.Error())
	// 	return nil, err
	// }

	// Crea una solicitud HTTP POST con el JSON como cuerpo
	// req, err := http.NewRequest("GET", serverURL, bytes.NewBuffer(jsonData)) ESTO MANDA EL EMAIL POR BODY NO POR URL
	req, err := http.NewRequest("GET", serverURL, nil)
	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return nil, err
	}

	// Establece el encabezado de tipo de contenido
	req.Header.Set("Content-Type", "application/json")

	// Realiza la solicitud HTTP
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error al realizar la solicutad HTTP", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	// defer func() { // De aquí saltaba un error REVISAR
	// 	if resp.Body != nil {
	// 		resp.Body.Close()
	// 		log.Println("Se ha cerrado el cuerpo de la solicitud")
	// 	}
	// }()

	// Verifica la respuesta del servidor (resp.StatusCode) aquí si es necesario
	log.Println("Respuesta del servidor: ", resp.Status)
	if resp.StatusCode == http.StatusNoContent {
		log.Println("Informacion: No fue encontrada ninguna máquina virtual para este usuario, ", err)
		return nil, errors.New("no fue encontrada ninguna máquina virtual")
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: La solicitud no fue exitosa, ", err)
		return nil, errors.New("la solicitud al servidor no fue exitosa")
	}

	// Lee la respuesta del cuerpo de la respuesta HTTP
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error al leer el cuerpo de la respuesta", err.Error())
		return nil, err
	}

	// Utilizar un mapa genérico para decodificar el cuerpo JSON
	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("error al decodificar el JSON: %w", err)
	}

	// Extraer el campo "data" del mapa
	data, ok := result["data"]
	if !ok {
		return nil, errors.New("el campo 'data' no se encuentra en la respuesta")
	}

	machinesData, err := json.Marshal(data) // Convertir "data" a JSON para luego deserializarlo

	if err != nil {
		return nil, fmt.Errorf("error al procesar el campo 'data': %w", err)
	}

	var machines []Models.VirtualMachine

	// Decodifica los datos de respuesta en la variable machines.
	if err := json.Unmarshal(machinesData, &machines); err != nil {
		// Maneja el error de decodificación aquí
		log.Println("error al decodifcar el JSON", err.Error())
		return nil, fmt.Errorf("error al decodificar JSON de respuesta: %w", err)
	}

	return machines, nil
}

// Enviar creación de la VM al servidor
func CreateMachineFromServer(VM Models.VirtualMachineTemp, clienteIp string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.VIRTUAL_MACHINE_URL)

	payload := map[string]interface{}{
		"specifications": VM,
		"clientIP":       clienteIp,
	}

	confirmacion, err := SendRequest("POST", serverURL, payload)
	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	// Esta confirmacion es digamos, del servidor web. Todo fue bien en la creación, pero no sabemos
	// si todo fue bien en el servidor de procesamiento, por lo que se debe verificar la existencia de la VM
	if confirmacion {
		maquinaCreada, err := VerifyMachineCreated(VM.Name, VM.Person_Email)
		if err != nil {
			log.Println("Error al consultar si la maquina fue creada: ", err.Error())
			return false, err
		}

		if !maquinaCreada {
			log.Println("La máquina no fue creada")
			return false, err
		}

		// Y este ya verifica por completo de que la vm efectivamente si existe :y:
		return maquinaCreada, nil
	}

	// Claro, si el web dice que algo ocurrió mal, mande false
	return confirmacion, nil
}

func VerifyMachineCreated(vmName, email string) (bool, error) {
	// TODO: Cambiar esto, porque acá se obtienen todas las mquinas y se verifica si existe para este user
	//		 pero deberia ser (en el servidor de procesamiento) que se obtenga solo un bool, de si existe
	//       un registro en la bd donde clienteEmail 'x' tiene una vm con nombre 'y'

	// Esto puede ser un poco lento, hablando de rendimiento y consultas y llamados API, porque se
	// podria obtener solo un bool (EXISTE VM) que todo un JSON con todas las vm del user. y verificarlos aquí

	const intentos = 3
	const tiempoEspera = 4 * time.Second

	// Como no se puede verificar el nombre exacto (por que el servidor de procesamiento le agrega chars aleatorios)
	// entonces se verifica si la vm con el nombre (sin las chars aleatorios) fue creada en un intervalo de unos 10 segundos
	tiempoActual := time.Now()

	for intento := 1; intento <= intentos; intento++ {

		// Se trae todas las vm del user y se verifica si ya se le creó la vm en la BD
		machines, err := ConsultMachineFromServer(email)
		fmt.Println("Machines: ", machines)
		if err != nil {
			log.Printf("Intento %d: Error al consultar si la máquina fue creada: %v", intento, err)

			if intento == intentos {
				return false, err
			}

			// Como hubo un error entonces esperemos un momento, por si el servidor de procesamiento se está demorando
			time.Sleep(tiempoEspera)
		}

		for _, machine := range machines {
			nombreSinAleatorios := machine.Name[:len(machine.Name)-5]

			// Esto se hace por si el usuario, le da por crear otra vm con el mismo nombre y pues eso lo debe
			// dejar hacer el servidor. Para eso se verifica las fechas de creacion a la hora de crear una vm
			if vmName == nombreSinAleatorios {

				// TODO: Si aparecen errores raros, implementar los log.print del commit de [web 508e559].
				// tambien pueden darse errores si la duracion es negativa.
				duracion := machine.Fecha_creacion.Sub(tiempoActual)

				intervaloMaximo := 20 * time.Second //TODO: AJUSTAR TIEMPO EXACTO-APROX | primero 30 segundos pa probar

				if duracion < intervaloMaximo {
					log.Println("Máquina encontrada | vmName: ", vmName, " machine.Name: ", machine.Name)
					return true, nil
				}
			}
		}

		log.Printf("Máquina no encontrada en el intento %d", intento)
		if intento < intentos {
			time.Sleep(tiempoEspera) // Se espera por si el servidor de procesamiento se está demorando
		}
	}

	return false, nil
}

// Encender Maquina virtual
func StartMachineFromServer(nombre string, clientIP string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.START_VM_URL)
	payload := map[string]interface{}{
		"tipo_solicitud": "start",
		"nombreVM":       nombre,
		"clientIP":       clientIP,
	}

	confirmacion, err := SendRequest("POST", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil

}

// Apagar Maquina virtual
func StopMachineFromServer(nombre string, clientIP string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.STOP_VM_URL)
	payload := map[string]interface{}{
		"tipo_solicitud": "stop",
		"nombreVM":       nombre,
		"clientIP":       clientIP,
	}

	confirmacion, err := SendRequest("POST", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil

}

// Eliminar una Maquina virtual
func DeleteMachineFromServer(nombre string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s/%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.VIRTUAL_MACHINE_URL, nombre)

	// ESTO NO SE DEBE MANDAR, ES UN METODO DELETE NO PERMITE PAYLOAD
	payload := map[string]interface{}{
		"tipo_solicitud": "delete",
		"nombreVM":       nombre,
	}

	confirmacion, err := SendRequest("DELETE", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil
}

// Consultar estado de la maquina virtual
func CheckStatusMachineFromServer(VM Models.VirtualMachine, clienteIp string) (bool, error) {
	serverURL := fmt.Sprintf("http://%s:%s%s", Config.ServidorProcesamientoRoute, Config.PUERTO, Config.CHECK_HOST_URL)

	// ESTO NO SE DEBE MANDAR, ES UN METODO GET NO PERMITE PAYLOAD
	// Crear el objeto JSON con los datos del cliente
	payload := map[string]interface{}{
		"clientIP":       clienteIp,
		"ubicacion":      VM.Host_id,
		"specifications": VM,
	}

	confirmacion, err := SendRequest("GET", serverURL, payload)

	if err != nil {
		log.Println("Error al crear la solicitud HTTP", err.Error())
		return false, err
	}

	return confirmacion, nil

}
