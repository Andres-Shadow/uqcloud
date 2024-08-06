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

	// Inicializa la ruta de la llave privada SSH
	config.InitPrivateKeyPath(*privateKeyPath)

	// Conexión a SQL
	database.ManageSqlConecction()

	// Migración de modelos para la creación de las tablas en la base de datos.
	database.DBConnection()
	database.DATABASE.AutoMigrate(models.Persona{})
	database.DATABASE.AutoMigrate(models.Maquina_virtual{})
	database.DATABASE.AutoMigrate(models.Host{})
	database.DATABASE.AutoMigrate(models.Catalogo{})
	database.DATABASE.AutoMigrate(models.Disco{})
	database.DATABASE.AutoMigrate(models.Imagen{})
	database.DATABASE.AutoMigrate(models.Conetendor{})

	// Configura un manejador de solicitud para la ruta "/json".
	manageServer()
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
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Println("Error al iniciar el servidor:", err)
	}

}

/*
Funciòn que se encarga de configurar los endpoints, realizar las validaciones correspondientes a los JSON que llegan
por solicitudes HTTP. Se encarga tambièn de ingresar las peticiones para gestiòn de MV a la cola.
Si la peticiòn es de inicio de sesiòn, la gestiona inmediatamente.
*/
func manageServer() {

	config.InitQueues()

	//Endpoint para las peticiones de creaciòn de màquinas virtuales
	http.HandleFunc("/json/createVirtualMachine", handlers.CreateVirtualMachineHandler)

	//Endpoint para consultar los Host
	http.HandleFunc("/json/consultHosts", handlers.ConsultHostHandler)

	//Endpoint para checkear el Host seleccionado por el usuario caso de uso asignacioon de recursos
	http.HandleFunc("/json/checkhost", handlers.CheckHostHandler)

	//Endpoint para peticiones de inicio de sesiòn
	http.HandleFunc("/json/login", handlers.UserLoginHandler)

	//Endpoint para peticiones de inicio de sesiòn
	http.HandleFunc("/json/signin", handlers.UserSignInHandler)

	//Endpoint para consultar las màquinas virtuales de un usuario
	http.HandleFunc("/json/consultMachine", handlers.ConsultMachineHandler)

	//Endpoint para consultar los Host
	http.HandleFunc("/json/consultHost", handlers.ConsultHostHandler)

	//Endpoint para consultar el catàlogo
	http.HandleFunc("/json/consultCatalog", handlers.ConsultCatalogHandler)

	//End point para modificar màquinas virtuales
	http.HandleFunc("/json/modifyVM", handlers.ModifyVirtualMachineHandler)

	//End point para eliminar màquinas virtuales
	http.HandleFunc("/json/deleteVM", handlers.DeleteVirtualMachineHandler)

	//End point para encender màquinas virtuales
	http.HandleFunc("/json/startVM", handlers.StartVirtualMachineHandler)

	//End point para apagar màquinas virtuales
	http.HandleFunc("/json/stopVM", handlers.StopVirtualMachineHandler)

	//End point para crear maquinas virtuales para invitados
	http.HandleFunc("/json/createGuestMachine", handlers.CreateGuestVirtualMachineHandler)

	//Endpoint para agregar un host
	http.HandleFunc("/json/addHost", handlers.AddHostHandler)

	//Endpoint para agregar un disco
	http.HandleFunc("/json/addDisk", handlers.AddDiskHandler)

	//Endpoint para consultar las metricas
	http.HandleFunc("/json/consultMetrics", handlers.ConsultMetricsHandler)

	//Endpoint para crear imagen docker desde dockerhub
	http.HandleFunc("/json/imagenHub", handlers.CreateImageDockerHubHandler)

	//Endpoint para crear imagen docker desde archivo tar
	http.HandleFunc("/json/imagenTar", handlers.CreateImageDockerTarHandler)

	//Endpoint para crear imagen docker desde archivo Dockerfile
	http.HandleFunc("/json/imagenDockerFile", handlers.CreateImageDockerfileHandler)

	//Endpoint para eliminar imagen docker
	http.HandleFunc("/json/eliminarImagen", handlers.DeleteDockerImageHandler)

	//Endpoint para consultar las imagenes de docker en una maquina virtual
	http.HandleFunc("/json/imagenesVM", handlers.CheckVirtualMachineDockerImagesHandler)

	//Endpoint para crear contenedor
	http.HandleFunc("/json/crearContenedor", handlers.CreateDockerHandler)

	//Endpoint para administrar el listado de contenedores en una maquina virtual
	http.HandleFunc("/json/gestionContenedor", handlers.ManageDockerImagesHandler)

	//Endpoint para consultar los contenedores de una maquina virtual
	http.HandleFunc("/json/ContenedoresVM", handlers.CheckContainersHandler)
}
