# Desktop-Cloud-Web

Antes de ejecutar los servidores recordar:
    
    - Si se va a ejecutar con compose crear la red con: "docker network create uqcloud"
    - Obtener la llave privada rsa y ubicarla en servidor_procesamiento/Procesador/Keys
    - Crear un .env en el mismo nivel de Docker-Compose.yml con la variable de entorno      
        DB_PASSWORD
    - Crear un .env en el directorio servidor_procesamiento/Procesador/Environment con
    las variables de entorno:
        DB_PASSWORD (la misma que el .env anterior)
        DEFAULT_QUICK_VM_DISTRO
    
Para los archivos .env se brinda un archivo example en cada lugar donde debe ir y con las
variables de ejemplo.
