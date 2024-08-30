# Usa una imagen base de Go
FROM golang:1.21 AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de dependencias
COPY go.mod go.sum ./

# Descarga las dependencias con ayuda de la imagen de golang
RUN go mod download

# Copia el código fuente
COPY ./Procesador/. ./Procesador

# Compila el binario
RUN cd /app/Procesador && set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o /binario .

# Crea una imagen mínima para ejecutar el binario
FROM debian:sid-slim

# Copia el binario compilado desde la fase de construcción
COPY --from=builder ./binario ./binario

# Copia el archivo de configuracion
COPY --from=builder /app/Procesador/Environment/.env /Environment/.env

# Copia la key del builder de la carpeta Keys
COPY --from=builder /app/Procesador/Keys/id_rsa /Environment/id_rsa

# Copia la key del builder de la carpeta DatosHostJson
COPY --from=builder /app/Procesador/DatosHostJson/ /DatosHostJson/

# Asegúrate de que el binario tiene permisos de ejecución
RUN chmod +x ./binario

# Instalar openssh para poder conectarse a los servidores
RUN apt-get update && apt-get install -y openssh-client

# Ejecuta el binario
CMD ["./binario", "-h", "-key=/Environment/id_rsa"] 