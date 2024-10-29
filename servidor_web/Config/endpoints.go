package Config

//CLASE EN DONDE SE ESTABLECEN LAS URL, CADA URL PUEDE TENER ASOCIADAS
//DIFERENTES FUNCIONES DEACUERDO AL METODO CON EL QUE SEA LLAMADO

var PUERTO = "8081"

var API_PREFIX = "/api/v1/"

// URL LOGIN
var LOGIN_URL = API_PREFIX + "login"
var TEMP_USER_ACCOUNT = API_PREFIX + "temp-user"

// URL asociadas a las maquinas virtuales
var VIRTUAL_MACHINE_URL = API_PREFIX + "virtual-machine"
var START_VM_URL = API_PREFIX + "start-virtual-machine"
var STOP_VM_URL = API_PREFIX + "stop-virtual-machine"
var QUICK_VM = API_PREFIX + "quick-virtual-machine"

// URL asociadas a los host
var HOSTS_URL = API_PREFIX + "hosts"
var HOST_URL = API_PREFIX + "host"

// URL asociadas a los discos
var DISK_VM_URL = API_PREFIX + "disk"
var DISK_NAMES_URL = API_PREFIX + "disks"

// URL asociadas a las metricas
var METRICS_URL = API_PREFIX + "metrics"
