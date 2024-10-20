package main

import (
	authentication "AppWeb/Auth"
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
	store := cookie.NewStore([]byte("ADMIN_SESSION_SUPER_SECRET"))
	router.Use(sessions.Sessions("sesion", store))

	// Configura las rutas
	router.LoadHTMLGlob("web/templates/*.html")
	router.Static("/web/static", "./web/static")

	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	router.GET("/aboutus", func(c *gin.Context) { c.HTML(http.StatusOK, "aboutUs.html", nil) })
	router.GET("/docs", func(c *gin.Context) { c.HTML(http.StatusOK, "docs.html", nil) })

	//TODO: Revisar para que siver o eliminar si se puede
	router.GET("/signin", handlers.SigninPage)

	// Explicación para julian de julian: no hay necesidad de asignarle el "adminGroup" al "router",
	// porque directamente cuando se le asocia una variable con ":=" al "router", gin los junta directamente
	// sin necesidad de escribir una funcion digamos: router.setGroups( []grupos ). ya tu sabe tu si entiendes
	// --- RUTAS DE ADMIN ---
	authGroup := router.Group("/admin")
	{
		authGroup.GET("", handlers.LoginAdminPage)
		authGroup.POST("", handlers.AdminLogin)
	}

	adminGroup := router.Group("/auth-admin")
	adminGroup.Use(authentication.AuthRequired)
	{
		adminGroup.GET("/dashboard", handlers.DashboardHandler)
		adminGroup.GET("/create-host", handlers.CreateHostPage)
		adminGroup.GET("/create-disk", handlers.CreateDiskPage)
		//adminGroup.DELETE("/deleteHosts", handlers.DeleteHost)
	}

	// --- RUTAS DE USUARIO COMUN ---
	userGroup := router.Group("/mainpage")
	{
		userGroup.GET("/control-machine", handlers.ControlMachine)
		userGroup.GET("/control-machine/create-machine", handlers.CreateMachinePage)
		userGroup.GET("/profile", handlers.ProfilePage)
		userGroup.GET("/imagenes", handlers.GestionImagenes)
		userGroup.GET("/contenedores", handlers.GestionContenedores)
		userGroup.GET("/connection-machine", handlers.ConnectionMachine)
		// router.GET("actualizaciones-maquinas", handlers.ActualizacionesMaquinas)
		// router.GET("/helpCenter", handlers.HelpCenterPage)
	}

	//router.POST("/signin", handlers.Signin)

	router.GET("/GetHost", handlers.GetHosts)
	router.POST("/deleteHosts", handlers.DeleteHost)

	//TODO: Mirar después

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

	// API ROUTES
	router.POST("/api/quick-machine", handlers.QuickMachine)

	router.GET("/api/machines", handlers.GetMachines)

	router.POST("/api/createMachine", handlers.MainSend)
	router.POST("/api/startMachine", handlers.StartMachine)
	router.POST("/api/stopMachine", handlers.StopMachine)
	router.POST("/api/deleteMachine", handlers.DeleteMachine)

	router.POST("/api/loginTemp", handlers.LoginTemp)
	router.POST("/api/contendores", handlers.GetContendores)
	router.POST("/api/imagenes", handlers.GetImages)
	router.POST("/api/checkhost", handlers.Checkhost)
	// router.POST("/api/mvtemp", handlers.Mvtemp)

	router.POST("/cambiar-contenido", Utilities.SendContent)
	router.POST("/uploadJSON", Utilities.UploadJSON)

	// Ruta para cerrar sesión
	router.GET("/logout", handlers.Logout)

	// Iniciar la aplicación
	err := router.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
