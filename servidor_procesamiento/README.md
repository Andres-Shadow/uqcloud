# Servidor de Procesamiento

El **Servidor de Procesamiento** es una solución escrita en Golang para gestionar máquinas virtuales y sus recursos asociados. Utiliza **Gorilla Mux** para la gestión de peticiones HTTP y **GORM** para la interacción con la base de datos. El servidor permite realizar operaciones de creación, eliminación, encendido y apagado de máquinas virtuales, así como registrar hosts y discos duros, y gestionar sesiones de usuarios.

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


