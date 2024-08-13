package main

import (
	"AppWeb/handlers"
	"encoding/json"
	"html/template"
	"log"
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

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"json": func(v interface{}) string {
			jsonValue, _ := json.Marshal(v)
			return string(jsonValue)
		},
		"mod": mod,
	})

	// Carga las plantillas
	r.LoadHTMLGlob("templates/*.html")

	// Configurar la tienda de cookies para las sesiones
	store := cookie.NewStore([]byte("tu_clave_secreta"))
	r.Use(sessions.Sessions("sesion", store))

	// Configura las rutas
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	r.GET("/index", handlers.Index)

	//TODO: Revisar
	r.POST("/api/checkhost", handlers.Checkhost)

	//TODO: Cambiar ruta por /admin
	r.GET("/login", handlers.LoginPage)

	//TODO: Revisar para que siver o eliminar si se puede
	r.GET("/signin", handlers.SigninPage)

	//TODO: Adaptar a los usuarios temporales
	r.GET("/mainPage", handlers.MainPage)

	//TODO: RUTAS NO DEBERIA SER RUTAS Y VAN EN EL /mainpage
	r.GET("/profile", handlers.ProfilePage)
	//TODO: No debe ser una ruta
	r.GET("/scrollmenu", handlers.Scrollmenu)
	r.GET("actualizaciones-maquinas", handlers.ActualizacionesMaquinas)
	r.GET("/imagenes", handlers.GestionImagenes)
	r.GET("/contenedores", handlers.GestionContenedores)
	r.GET("/aboutUs", handlers.AboutUsPage)
	r.GET("/helpCenter", handlers.HelpCenterPage)

	//TODO: ELIMINAR ¿?
	r.GET("/welcome", handlers.WelcomePage)

	//TODO: DEJAR QUIETO HASTA TENER ACCESO AL ADMINISTRADOR
	r.GET("/dashboard", handlers.DashboardHandler)

	//TODO: Deberia ser un formulario dentro de dashboard y no una ruta
	r.GET("/createHost", handlers.CreateHostPage)
	r.GET("/createDisk", handlers.CreateDiskPage)

	r.GET("/api/machines", handlers.GetMachines)
	r.GET("/controlMachine", handlers.ControlMachine)

	r.POST("/login", handlers.Login)
	r.POST("/signin", handlers.Signin)

	//TODO: Mirar después
	r.POST("/api/createMachine", handlers.MainSend)
	r.POST("/powerMachine", handlers.PowerMachine)
	r.POST("/deleteMachine", handlers.DeleteMachine)
	r.POST("/configMachine", handlers.ConfigMachine)
	r.POST("/api/loginTemp", handlers.LoginTemp)
	//ToDo: Revisar la ruta
	r.POST("/createHost", handlers.CreateNewHost) //Cambio en el metodo, pero hay que revisar la ruta ToDo: Probar el nuevo método
	r.POST("/createDisk", handlers.CreateDisk)
	r.POST("/DockerHub", handlers.CrearImagen)
	r.POST("/CrearImagenTar", handlers.CrearImagenArchivoTar)
	r.POST("/CrearDockerFile", handlers.CrearImagenDockerFile)
	r.POST("/eliminarImagen", handlers.EliminarImagen)
	r.POST("/eliminarImagenes", handlers.EliminarImagenes)
	r.POST("/crearContenedor", handlers.CrearContenedor)
	r.POST("/CorrerContenedor", handlers.CorrerContenedor)
	r.POST("/PausarContenedor", handlers.PausarContenedor)
	r.POST("/ReiniciarContenedor", handlers.ReiniciarContenedor)
	r.POST("/EliminarContenedor", handlers.EliminarContenedor)
	r.POST("/eliminarContenedores", handlers.EliminarContenedores)

	r.POST("/api/contendores", handlers.GetContendores)
	r.POST("/api/imagenes", handlers.GetImages)

	r.POST("/cambiar-contenido", handlers.EnviarContenido)

	r.POST("/uploadJSON", handlers.HandleUploadJSON)
	r.POST("/api/mvtemp", handlers.Mvtemp)
	// Ruta para cerrar sesión
	r.GET("/logout", handlers.Logout)

	// Iniciar la aplicación
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
