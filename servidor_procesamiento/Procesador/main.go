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
	models "servidor_procesamiento/Procesador/Models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Variable que almacena la ruta de la llave privada ingresada por paametro cuando de ejecuta el programa
var privateKeyPath = flag.String("key", "", "Ruta de la llave privada SSH")

func main() {

	flag.Parse()

	//Verifica que el paràmetro de la ruta de la llave privada no estè vacìo
	if *privateKeyPath == "" {
		fmt.Println("Debe ingresar la ruta de la llave privada SSH")
		return
	}

	// Carga las variables de entorno del archivo .env
	err := godotenv.Load("Environment/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	r := mux.NewRouter()

	// Inicializa la ruta de la llave privada SSH
	config.InitPrivateKeyPath(*privateKeyPath)

	// Inicializa la conexión a la base de datos y precarga de datos
	setDatabase()

	// Configura un manejador de solicitud para la ruta "/json".
	manageServer(r)
	// Función que verifica la cola de especificaciones constantemente.
	go jobs.CheckVirtualMachinesQueueChanges()
	//Funciòn que verifica el tiempo de creaciòn de una MV
	//go checkTime(privateKeyPath)
	// Función que verifica la cola de cuentas constantemente.
	go jobs.CheckManagementQueueChanges()
	// Función que verifica la cola de imagenes docker constantemente.
	go jobs.CheckImagesQueueChanges()
	// Función que verifica la cola de contenedores constantemente.
	go jobs.CheckContainerQueueChanges()
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
	// Conexión a SQL
	//database.ManageSqlConecction()

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

	if !database.CountAdminsRegistered() {
		persona := models.Persona{
			Nombre:      "admin",
			Apellido:    "admin",
			Email:       "admin@uqcloud.co",
			Contrasenia: "$2y$10$JGxWitiykfO83Ep8IBab/.3fn.H/DxMjAK8dFTQCPZyJ5EHqZtfji", // Dejar este hash bcrypt para la contraseña "admin"
			Rol:         1,
		}

		database.DATABASE.Create(&persona)
		fmt.Print("Usuario administrador creado\n")
	}
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
	r.HandleFunc(apiPrefix+"virtual_machine", handlers.CreateVirtualMachineHandler).Methods("POST")

	//Endpoint para consultar las màquinas virtuales de un usuario
	r.HandleFunc(apiPrefix+"virtual_machine/{email}", handlers.ConsultVirtualMachineHandler).Methods("GET")

	//End point para modificar màquinas virtuales
	r.HandleFunc(apiPrefix+"virtual_machine", handlers.ModifyVirtualMachineHandler).Methods("PUT")

	//End point para eliminar màquinas virtuales
	r.HandleFunc(apiPrefix+"virtual_machine", handlers.DeleteVirtualMachineHandler).Methods("DELETE")

	//End point para encender màquinas virtuales
	r.HandleFunc(apiPrefix+"start_virtual_machine", handlers.StartVirtualMachineHandler).Methods("POST")

	//End point para apagar màquinas virtuales
	r.HandleFunc(apiPrefix+"stop_virtual_machine", handlers.StopVirtualMachineHandler).Methods("POST")

	//End point para crear maquinas virtuales para invitados
	r.HandleFunc(apiPrefix+"guest_virtual_machine", handlers.CreateGuestVirtualMachineHandler).Methods("POST")

	/*
		--------------------
		| HOST ENDPOINTS   |
		-------------------
	*/

	//Endpoint para consultar los Host
	r.HandleFunc(apiPrefix+"hosts", handlers.ConsultHostsHandler).Methods("GET")

	//Endpoint para checkear el Host seleccionado por el usuario caso de uso asignacioon de recursos
	r.HandleFunc(apiPrefix+"check_host", handlers.CheckHostHandler).Methods("GET")

	//Endpoint para consultar los Host
	r.HandleFunc(apiPrefix+"host/{email}", handlers.ConsultHostHandler).Methods("GET")

	//Endpoint para agregar un host
	r.HandleFunc(apiPrefix+"host", handlers.AddHostHandler).Methods("POST")

	/*
		------------------
		| USER ENDPOINTS |
		-----------------
	*/

	// Endpoint para peticiones de inicio de sesiòn
	r.HandleFunc(apiPrefix+"login", handlers.UserLoginHandler)

	//Endpoint para peticiones de inicio de sesiòn
	// r.HandleFunc("/json/signin", handlers.UserSignInHandler)

	/*
		-----------------------
		|  CATALOG ENDPOINTS  |
		----------------------
	*/

	//Endpoint para consultar el catàlogo
	r.HandleFunc(apiPrefix+"catalog", handlers.ConsultCatalogHandler).Methods("GET")

	/*
		--------------------
		| DISK ENDPOINTS   |
		-------------------
	*/

	//Endpoint para agregar un disco
	r.HandleFunc(apiPrefix+"disk", handlers.AddDiskHandler).Methods("POST")

	/*
		-----------------------
		| METRICS ENDPOINTS   |
		----------------------
	*/

	//Endpoint para consultar las metricas
	r.HandleFunc(apiPrefix+"metrics", handlers.ConsultMetricsHandler).Methods("GET")

	/*
		--------------------
		| DOCKER ENDPOINTS |
		-------------------
	*/

	//Endpoint para crear imagen docker desde dockerhub
	r.HandleFunc(apiPrefix+"dockerhub_image", handlers.CreateImageDockerHubHandler).Methods("POST")

	//Endpoint para crear imagen docker desde archivo tar
	r.HandleFunc(apiPrefix+"tar_image", handlers.CreateImageDockerTarHandler).Methods("POST")

	//Endpoint para crear imagen docker desde archivo Dockerfile
	r.HandleFunc(apiPrefix+"dockerfile_image", handlers.CreateImageDockerfileHandler).Methods("POST")

	//Endpoint para eliminar imagen docker
	r.HandleFunc(apiPrefix+"docker_image", handlers.DeleteDockerImageHandler).Methods("DELETE")

	//Endpoint para consultar las imagenes de docker en una maquina virtual
	r.HandleFunc(apiPrefix+"virtual_machine_images", handlers.CheckVirtualMachineDockerImagesHandler).Methods("GET")

	//Endpoint para crear contenedor
	r.HandleFunc(apiPrefix+"docker", handlers.CreateDockerHandler).Methods("POST")

	//Endpoint para administrar el listado de contenedores en una maquina virtual
	r.HandleFunc(apiPrefix+"docker", handlers.ManageDockerImagesHandler).Methods("PUT")

	//Endpoint para consultar los contenedores de una maquina virtual
	r.HandleFunc(apiPrefix+"virtual_machine_docker", handlers.CheckContainersHandler).Methods("GET")
}
