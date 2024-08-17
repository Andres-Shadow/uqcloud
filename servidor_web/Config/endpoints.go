package Config

//CLASE EN DONDE SE ESTABLECEN LAS URL, CADA URL PUEDE TENER ASOCIADAS
//DIFERENTES FUNCIONES DEACUERDO AL METODO CON EL QUE SEA LLAMADO

var PUERTO = "8081"

//URL LOGIN
var LOGIN_URL = "/api/v1/login"

// URL asociadas a las maquinas virtuales
var VIRTUAL_MACHINE_URL = "/api/v1/virtual_machine"
var START_VM_URL = "/api/v1/start_virtual_machine"
var CREATE_GUEST_VM_URL = "/api/v1/guest_virtual_machine"

// URL asociadas a los host
var HOSTS_URL = "/api/v1/hosts"
var CHECK_HOST_URL = "api/v1/check_host"
var HOST_URL = "api/v1/host"

// URL asociadas a los discos
var DISK_VM_URL = "/api/v1/disk"

// URL asociadas a las metricas
var METRICS_URL = "/api/v1/metrics"

// URL asociadas a la imagenes
var IMAGEN_HUB_URL = " /api/v1/dockerhub_image"
var IMAGEN_TAR_URL = "/api/v1/tar_image"
var IMAGEN_DOCKER_FILE_URL = "/api/v1/dockerfile_image"
var DELETE_IMAGEN_URL = "/api/v1/docker_image" //el nombre de las variables podría cambiar no se me ocurrió ninguno más
var VM_IMAGES_URL = "/api/v1/virtual_machine_image"

// URL CONTENEDORES
var CONTAINER_URL = "/api/v1/docker"
var VM_DOCKER_URL = "/api/v1/virtual_machine_docker"
