package Config

//CLASE EN DONDE SE ESTABLECEN LAS URL, CADA URL PUEDE TENER ASOCIADAS
//DIFERENTES FUNCIONES DEACUERDO AL METODO CON EL QUE SEA LLAMADO

var PUERTO = "9091"

var API_PREFIX = "/api/v1/"

//URL LOGIN
var LOGIN_URL = API_PREFIX + "login"
var TEMP_USER_ACCOUNT = API_PREFIX + "temp-user"

// URL asociadas a las maquinas virtuales
var VIRTUAL_MACHINE_URL = API_PREFIX + "virtual-machine"
var START_VM_URL = API_PREFIX + "start-virtual-machine"
var STOP_VM_URL = API_PREFIX + "stop-virtual-machine"
var QUICK_VM = API_PREFIX + "quick-virtual-machine"

// URL asociadas a los host
var HOSTS_URL = API_PREFIX + "hosts"
var CHECK_HOST_URL = API_PREFIX + "check-host"
var HOST_URL = API_PREFIX + "host"

// URL asociadas a los discos
var DISK_VM_URL = API_PREFIX + "disk"

// URL asociadas a las metricas
var METRICS_URL = API_PREFIX + "metrics"

// URL asociadas a la imagenes
// var IMAGEN_HUB_URL = API_PREFIX + "dockerhub_image"
// var IMAGEN_TAR_URL = API_PREFIX + "tar_image"
// var IMAGEN_DOCKER_FILE_URL = "/api/v1/dockerfile_image"
// var DELETE_IMAGEN_URL = "/api/v1/docker_image" //el nombre de las variables podría cambiar no se me ocurrió ninguno más
// var VM_IMAGES_URL = "/api/v1/virtual_machine_image"

// // URL CONTENEDORES
// var CONTAINER_URL = "/api/v1/docker"
// var VM_DOCKER_URL = "/api/v1/virtual_machine_docker"
