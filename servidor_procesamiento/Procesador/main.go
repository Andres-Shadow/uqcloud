package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	config "servidor_procesamiento/Procesador/Config"
	database "servidor_procesamiento/Procesador/Database"
	handlers "servidor_procesamiento/Procesador/Handlers"
	jobs "servidor_procesamiento/Procesador/Jobs"
	models "servidor_procesamiento/Procesador/Models/Entities"
	utilities "servidor_procesamiento/Procesador/Utilities"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Variable que almacena la ruta de la llave privada ingresada por paametro cuando de ejecuta el programa
var privateKeyPath = flag.String("key", "./keys/id_rsa", "")
var preregisteredHosts = flag.Bool("h", false, "")

func main() {

	flag.Parse()

	// Carga las variables de entorno del archivo .env
	err := godotenv.Load("Environment/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Inicializa la conexión a la base de datos y precarga de datos
	setDatabase()

	//Verifica que el paràmetro de la ruta de la llave privada no estè vacìo
	if *privateKeyPath == "" {
		fmt.Println("Error al encontrar la llave privada externa especificada")
		return
	} else {
		fmt.Println("Iniciando servidor con llave privada: ", *privateKeyPath)
	}

	// Inicializa el RoundRobinManager
	setRoundRobinManager()

	r := mux.NewRouter()

	// Inicializa la ruta de la llave privada SSH
	config.InitPrivateKeyPath(*privateKeyPath)

	// Verifica si se ingresò el paràmetro para precargar los hosts
	if *preregisteredHosts {
		registerHostData()
	}

	// Recargar configuración de Prometheus
	config.ReloadPrometheusConfig()

	// Configura un manejador de solicitud para la ruta "/json".
	manageServer(r)
	// Función que verifica la cola de especificaciones constantemente.
	go jobs.CheckVirtualMachinesQueueChanges()
	//Funciòn que verifica el tiempo de creaciòn de una MV
	//go checkTime(privateKeyPath)
	// Función que verifica la cola de cuentas constantemente.
	go jobs.CheckManagementQueueChanges()
	// Inicia el servidor HTTP en el puerto 8081.
	fmt.Println("Servidor escuchando en el puerto 8081...")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Println("Error al iniciar el servidor:", err)
	}

}

/*
Funcion que se encarga de realizar la conexión a la base de datos, cargar los modelos a la base de datos y
precargar el usuario administrador
*/
func setDatabase() {

	// Migración de modelos para la creación de las tablas en la base de datos.
	database.DBConnection()
	database.DATABASE.AutoMigrate(
		&models.Persona{},
		&models.Maquina_virtual{},
		&models.Host{},
		&models.Catalogo{},
		&models.Disco{},
		&models.Imagen{},
		&models.Contenedor{},
		&models.CatalogoDisco{})

	// Precarga del usuario administrador
	database.CreateAdmin()

	// Actualizar las máquinas virtuales que estén disponibles realmente en los hosts
	// Esto se hace para que haya congruencia entre la BD y las VM existentes realmente
	utilities.UpdateVirtualMachinesActualStatus()
}

// Función para precargar los datos de los hosts de la sala B y C (No cambian)
func registerHostData() {
	// Obtener la cantidad de hosts registrados en la BD
	count, err := database.CountRegisteredHosts()
	if err != nil {
		log.Println("Error al contar los hosts registrados:", err)
		return
	}

	// Verificar que no hayan hosts registrados
	if count == 0 {
		fmt.Println("Preregistrando datos de hosts...")
		utilities.PreregisterHostJsonData()
	} else {
		fmt.Println("Ya existen hosts registrados")
	}
}

func setRoundRobinManager() {
	registeredHosts := database.GetHosts()
	config.RoundRobinManager = config.NewRoundRobin(registeredHosts)
}

/*
Funciòn que se encarga de configurar los endpoints, realizar las validaciones correspondientes a los JSON que llegan
por solicitudes HTTP. Se encarga tambièn de ingresar las peticiones para gestiòn de MV a la cola.
Si la peticiòn es de inicio de sesiòn, la gestiona inmediatamente.
*/
func manageServer(r *mux.Router) {

	config.InitQueues()

	var apiPrefix string = "/api/v1/"

	/*
		--------------------------------
		| VIRTUAL MACHINE ENDPOINTS    |
		-------------------------------
	*/

	//Endpoint para las peticiones de creaciòn de màquinas virtuales
	r.HandleFunc(apiPrefix+"virtual-machine", handlers.CreateVirtualMachineHandler).Methods("POST")

	//Endpoint para consultar las màquinas virtuales de un usuario
	r.HandleFunc(apiPrefix+"virtual-machine/{email}", handlers.ConsultVirtualMachineHandler).Methods("GET")

	//End point para eliminar màquinas virtuales
	r.HandleFunc(apiPrefix+"virtual-machine/{name}", handlers.DeleteVirtualMachineHandler).Methods("DELETE")

	//End point para encender màquinas virtuales
	r.HandleFunc(apiPrefix+"start-virtual-machine", handlers.StartVirtualMachineHandler).Methods("POST")

	//End point para apagar màquinas virtuales
	r.HandleFunc(apiPrefix+"stop-virtual-machine", handlers.StopVirtualMachineHandler).Methods("POST")

	//End point para crear una máquina rápida
	r.HandleFunc(apiPrefix+"quick-virtual-machine", handlers.CreateQuickVirtualMachineHandler).Methods("POST")

	/*
		--------------------
		| HOST ENDPOINTS   |
		-------------------
	*/

	//Endpoint para consultar los Host
	r.HandleFunc(apiPrefix+"hosts", handlers.ConsultHostsHandler).Methods("GET")

	//Endpoint para consultar los Host
	r.HandleFunc(apiPrefix+"host/{email}", handlers.ConsultHostHandler).Methods("GET")

	//Endpoint para agregar un host
	r.HandleFunc(apiPrefix+"host", handlers.AddHostHandler).Methods("POST")

	//Endpoint para eliminar un host
	r.HandleFunc(apiPrefix+"host/{id}", handlers.DeleteHostHandler).Methods("DELETE")

	//Endpoint para registro rapido de host
	r.HandleFunc(apiPrefix+"host-fast-register", handlers.FastRegisterHostsHandler).Methods("POST")

	/*
		--------------------
		| DISK ENDPOINTS   |
		-------------------
	*/

	//Endpoint para agregar un disco
	r.HandleFunc(apiPrefix+"disk", handlers.AddDiskHandler).Methods("POST")

	//Endpoint para consultar los sistemas operativos de los discos de forma única
	r.HandleFunc(apiPrefix+"disks", handlers.GetDisksHandler).Methods("GET")

	//Enpoint para consultar los host que tienen un disco en específico
	r.HandleFunc(apiPrefix+"disk/{name}", handlers.GetHostsWithDiskHandler).Methods("GET")

	//Endpoint para eliminar un disco enviando la distribucion del disco y el host donde se encuentra
	r.HandleFunc(apiPrefix+"disk/{name}", handlers.DeleteDiskHandler).Methods("DELETE")

	/*
		------------------
		| USER ENDPOINTS |
		-----------------
	*/

	// Endpoint para peticiones de inicio de sesiòn
	r.HandleFunc(apiPrefix+"login", handlers.UserLoginHandler)

	// Endpoint para crear un usuario nuevo tempral
	r.HandleFunc(apiPrefix+"temp-user", handlers.CreateTempUserHandler).Methods("POST")

	/*
		-----------------------
		|  CATALOG ENDPOINTS  |
		----------------------
	*/

	//Endpoint para consultar el catàlogo
	r.HandleFunc(apiPrefix+"catalog", handlers.ConsultCatalogHandler).Methods("GET")

	/*
		-----------------------
		| METRICS ENDPOINTS   |
		----------------------
	*/

	//Endpoint para consultar las metricas
	r.HandleFunc(apiPrefix+"metrics", handlers.ConsultMetricsHandler).Methods("GET")

}
