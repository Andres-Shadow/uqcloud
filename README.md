<h1 align="center">UQCLOUD</h1>

![Go](https://img.shields.io/badge/Go-1.21-blue)
![Docker](https://img.shields.io/badge/Docker-Enabled-blue)
![GORM](https://img.shields.io/badge/GORM-ORM-red)
![Gorilla Mux](https://img.shields.io/badge/Gorilla_Mux-Routing-green)
![SSH Client](https://img.shields.io/badge/SSH-Client-blue)
![GitHub contributors](https://img.shields.io/github/contributors/Andres-Shadow/uqcloud)

**UQCloud** es un entorno de computadoras de escritorios que permite aprovechar de forma oportunista y mediante virtualización, recursos informáticos desde cualquier dispositivo conectado a la red.

## Servidor de Procesamiento

El **Servidor de Procesamiento** es una solución escrita en Golang para gestionar máquinas virtuales y sus recursos asociados. Utiliza **Gorilla Mux** para la gestión de peticiones HTTP y **GORM** para la interacción con la base de datos. El servidor permite realizar operaciones de creación, eliminación, encendido y apagado de máquinas virtuales, así como registrar hosts y discos duros, y gestionar sesiones de usuarios.

## Servidor Web

El **Servidor Web** es una solución UI en Golang utilizando recursos Web3 para gestionar máquinas virtuales y sus funcionalidades asociadas desplegando una interfaz gráfica intuitiva y amigable con el usuario, conectandose por medio de peticiones HTTP al servidor de procesamiento.

## Tecnologías

- **Golang**: Lenguaje de programación principal.
- **Gorilla Mux**: Router y manejador de peticiones HTTP.
- **GORM**: ORM para la interacción con la base de datos.
- **SSHClient**: Para la gestión de conexiones SSH.

## Funcionalidades

- **Gestión de Máquinas Virtuales**:
  - Crear
  - Eliminar
  - Encender
  - Apagar

- **Registro de Recursos**:
  - Hosts donde se albergan las máquinas virtuales.
  - Discos duros utilizados por las máquinas virtuales.

- **Gestión de Sesiones de Usuario**.
  - Sesiones temporales que garantizan usuarios efímeros

## Flags de Configuración

- `-h`: Realiza un precargado de hosts desde archivos JSON ubicados en la carpeta `datoshost`.
- `-key`: Permite especificar una llave privada de conexión SSH diferente a la predeterminada en el paquete `keys`.

## Requisitos de Despliegue
- Si se va a desplegar el aplicativo en Docker, crear la network con el siguiente comando:
  `docker network create uqcloud`
- Ubicar la llave privada en la carpeta `servidor_procesamiento/Procesador/Keys`
- Crear los archivos `.env` en el directorio raiz (a mismo nivel del Docker-Compose.yml) y en la carpeta `servidor_procesamiento/Procesador/Environment`, siguiendo los parametros datos en los archivos `example` de cada uno


