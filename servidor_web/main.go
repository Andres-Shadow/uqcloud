package main

import (
	"AppWeb/Utilities"
	"AppWeb/handlers"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type LoginPageData struct {
	ErrorMessage string
}

// Funcion Para la matriz de botones del caso de uso de asignacion de recursos
func mod(i, j int) int {
	return i % j
}

func main() {
	args := os.Args[1]
	port := ":" + args

	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"json": func(v interface{}) string {
			jsonValue, _ := json.Marshal(v)
			return string(jsonValue)
		},
		"mod": mod,
	})

	// Carga las plantillas
	router.LoadHTMLGlob("web/templates/*.html")

	// Configurar la tienda de cookies para las sesiones
	store := cookie.NewStore([]byte("tu_clave_secreta"))
	router.Use(sessions.Sessions("sesion", store))

	// Configura las rutas
	router.LoadHTMLGlob("web/templates/*.html")
	router.Static("web/static", "./web/static")

	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	router.GET("/aboutus", func(c *gin.Context) { c.HTML(http.StatusOK, "aboutUs.html", nil) })
	router.GET("/docs", func(c *gin.Context) { c.HTML(http.StatusOK, "docs.html", nil) })

	//TODO: Revisar
	router.POST("/api/checkhost", handlers.Checkhost)

	//TODO: Cambiar ruta por /admin
	router.GET("/admin", handlers.LoginAdminPage)

	//TODO: Revisar para que siver o eliminar si se puede
	router.GET("/signin", handlers.SigninPage)

	//TODO: Adaptar a los usuarios temporales
	router.GET("/mainPage", handlers.MainPage)

	//TODO: RUTAS NO DEBERIAN SER RUTAS? Y VAN EN EL /mainpage
	// Por como está hecho este proyecto, archivos .html llaman a otros .html
	// ya sea usando JQuery o iframe, por lo cual el servidor debe exponer estas
	// rutas para acceder a las templates.
	// Se debe buscar una forma de reescribir esto, como no tener el scrollmenu.html
	// en otro archivo sino en el mismo mainPage.html. O haciendo que el servidor no
	// deje acceder a estas rutas desde el navegador.

	router.GET("/profile", handlers.ProfilePage)
	router.GET("/scrollmenu", handlers.Scrollmenu)
	// router.GET("actualizaciones-maquinas", handlers.ActualizacionesMaquinas)
	// router.GET("/imagenes", handlers.GestionImagenes)
	// router.GET("/contenedores", handlers.GestionContenedores)
	// router.GET("/helpCenter", handlers.HelpCenterPage)

	//TODO: ELIMINAR ¿?
	router.GET("/welcome", handlers.WelcomePage)

	//TODO: DEJAR QUIETO HASTA TENER ACCESO AL ADMINISTRADOR
	router.GET("/dashboard", handlers.DashboardHandler)

	//TODO: Deberia ser un formulario dentro de dashboard y no una ruta
	router.GET("/createHost", handlers.CreateHostPage)
	router.GET("/createDisk", handlers.CreateDiskPage)

	// TODO: DESCOMENTAR LUEGO
	router.GET("/api/machines", handlers.GetMachines)
	router.GET("/controlMachine", handlers.ControlMachine)

	router.POST("/admin", handlers.AdminLogin)
	//router.POST("/signin", handlers.Signin)

	//TODO: Mirar después
	router.POST("/api/createMachine", handlers.MainSend)
	router.POST("/powerMachine", handlers.PowerMachine)
	router.POST("/deleteMachine", handlers.DeleteMachine)
	router.POST("/configMachine", handlers.ConfigMachine)
	router.POST("/api/loginTemp", handlers.LoginTemp)
	router.POST("/createHost", handlers.CreateNewHost)
	router.POST("/createDisk", handlers.CreateNewDisk)
	router.POST("/DockerHub", handlers.CrearImagen)
	router.POST("/CrearImagenTar", handlers.CrearImagenArchivoTar)
	router.POST("/CrearDockerFile", handlers.CrearImagenDockerFile)
	router.POST("/eliminarImagen", handlers.EliminarImagen)
	router.POST("/eliminarImagenes", handlers.EliminarImagenes)
	router.POST("/crearContenedor", handlers.CrearContenedor)
	router.POST("/CorrerContenedor", handlers.CorrerContenedor)
	router.POST("/PausarContenedor", handlers.PausarContenedor)
	router.POST("/ReiniciarContenedor", handlers.ReiniciarContenedor)
	router.POST("/EliminarContenedor", handlers.EliminarContenedor)
	router.POST("/eliminarContenedores", handlers.EliminarContenedores)

	router.POST("/api/contendores", handlers.GetContendores)
	router.POST("/api/imagenes", handlers.GetImages)

	router.POST("/cambiar-contenido", Utilities.SendContent)

	router.POST("/uploadJSON", Utilities.UploadJSON)
	//router.POST("/api/mvtemp", handlers.Mvtemp)
	// Ruta para cerrar sesión
	router.GET("/logout", handlers.Logout)

	// Iniciar la aplicación
	err := router.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
