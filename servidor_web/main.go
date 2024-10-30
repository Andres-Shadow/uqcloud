package main

import (
	authentication "AppWeb/Auth"
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
		userGroup.GET("/connection-machine", handlers.ConnectionMachine)
	}

	router.GET("/GetHost", handlers.GetHosts)
	router.POST("/deleteHosts", handlers.DeleteHost)
	router.GET("/getDisks", handlers.GetDiskFromRequest)
	router.GET("/hostOfDisk/:diskName", handlers.GetHostOfDiskFromRequest)
	router.DELETE("/hostOfDisk/:diskName", handlers.DeleteHostOfDiskFrormRequest)

	router.POST("/createHost", handlers.CreateNewHost)
	router.POST("/createDisk", handlers.CreateNewDisk)

	// API ROUTES
	router.GET("/api/machines", handlers.GetMachines)
	router.GET("/api/sshKeyMachine/:vm_name", handlers.GetSSHKeyMachine)
	router.POST("/api/quick-machine", handlers.QuickMachine)
	router.POST("/api/createMachine", handlers.MainSend)
	router.POST("/api/startMachine", handlers.StartMachine)
	router.POST("/api/stopMachine", handlers.StopMachine)
	router.POST("/api/deleteMachine", handlers.DeleteMachine)
	router.POST("/api/loginTemp", handlers.LoginTemp)

	// Ruta para cerrar sesión
	router.GET("/logout", handlers.Logout)

	// Iniciar la aplicación
	err := router.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
