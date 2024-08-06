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
    err := godotenv.Load( "Enviroment/.env" )
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	r := mux.NewRouter()

	// Inicializa la ruta de la llave privada SSH
	config.InitPrivateKeyPath(*privateKeyPath)

	// Conexión a SQL
	database.ManageSqlConecction()

	// Migración de modelos para la creación de las tablas en la base de datos.
	database.DBConnection()
	database.DATABASE.AutoMigrate(
        &models.Persona{},
        &models.Maquina_virtual{},
        &models.Host{},
        &models.Catalogo{},
        &models.Disco{},
        &models.Imagen{},
        &models.Conetendor{},
    )

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
Funciòn que se encarga de configurar los endpoints, realizar las validaciones correspondientes a los JSON que llegan
por solicitudes HTTP. Se encarga tambièn de ingresar las peticiones para gestiòn de MV a la cola.
Si la peticiòn es de inicio de sesiòn, la gestiona inmediatamente.
*/
func manageServer(r *mux.Router) {

	config.InitQueues()

	//Endpoint para las peticiones de creaciòn de màquinas virtuales
	r.HandleFunc("/json/createVirtualMachine", handlers.CreateVirtualMachineHandler).Methods("POST")

	//Endpoint para consultar los Host
	r.HandleFunc("/json/consultHosts", handlers.ConsultHostsHandler).Methods("POST")

	//Endpoint para checkear el Host seleccionado por el usuario caso de uso asignacioon de recursos
	r.HandleFunc("/json/checkhost", handlers.CheckHostHandler).Methods("POST")

	//Endpoint para peticiones de inicio de sesiòn
	// r.HandleFunc("/json/login", handlers.UserLoginHandler)

	//Endpoint para peticiones de inicio de sesiòn
	// r.HandleFunc("/json/signin", handlers.UserSignInHandler)

	//Endpoint para consultar las màquinas virtuales de un usuario
	r.HandleFunc("/json/consultMachine", handlers.ConsultVirtualMachineHandler).Methods("POST")

	//Endpoint para consultar los Host
	r.HandleFunc("/json/consultHost", handlers.ConsultHostHandler).Methods("POST")

	//Endpoint para consultar el catàlogo
	r.HandleFunc("/json/consultCatalog", handlers.ConsultCatalogHandler).Methods("GET")

	//End point para modificar màquinas virtuales
	r.HandleFunc("/json/modifyVM", handlers.ModifyVirtualMachineHandler).Methods("POST")

	//End point para eliminar màquinas virtuales
	r.HandleFunc("/json/deleteVM", handlers.DeleteVirtualMachineHandler).Methods("POST")

	//End point para encender màquinas virtuales
	r.HandleFunc("/json/startVM", handlers.StartVirtualMachineHandler).Methods("POST")

	//End point para apagar màquinas virtuales
	r.HandleFunc("/json/stopVM", handlers.StopVirtualMachineHandler).Methods("POST")

	//End point para crear maquinas virtuales para invitados
	r.HandleFunc("/json/createGuestMachine", handlers.CreateGuestVirtualMachineHandler).Methods("POST")

	//Endpoint para agregar un host
	r.HandleFunc("/json/addHost", handlers.AddHostHandler).Methods("POST")

	//Endpoint para agregar un disco
	r.HandleFunc("/json/addDisk", handlers.AddDiskHandler).Methods("POST")

	//Endpoint para consultar las metricas
	r.HandleFunc("/json/consultMetrics", handlers.ConsultMetricsHandler).Methods("GET")

	//Endpoint para crear imagen docker desde dockerhub
	r.HandleFunc("/json/imagenHub", handlers.CreateImageDockerHubHandler).Methods("POST")

	//Endpoint para crear imagen docker desde archivo tar
	r.HandleFunc("/json/imagenTar", handlers.CreateImageDockerTarHandler).Methods("POST")

	//Endpoint para crear imagen docker desde archivo Dockerfile
	r.HandleFunc("/json/imagenDockerFile", handlers.CreateImageDockerfileHandler).Methods("POST")

	//Endpoint para eliminar imagen docker
	r.HandleFunc("/json/eliminarImagen", handlers.DeleteDockerImageHandler).Methods("POST")

	//Endpoint para consultar las imagenes de docker en una maquina virtual
	r.HandleFunc("/json/imagenesVM", handlers.CheckVirtualMachineDockerImagesHandler).Methods("POST")

	//Endpoint para crear contenedor
	r.HandleFunc("/json/crearContenedor", handlers.CreateDockerHandler).Methods("POST")

	//Endpoint para administrar el listado de contenedores en una maquina virtual
	r.HandleFunc("/json/gestionContenedor", handlers.ManageDockerImagesHandler).Methods("POST")

	//Endpoint para consultar los contenedores de una maquina virtual
	r.HandleFunc("/json/ContenedoresVM", handlers.CheckContainersHandler).Methods("POST")
}
