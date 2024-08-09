# Usa una imagen base de Go
FROM golang:1.21 AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de dependencias
COPY go.mod go.sum ./

# Descarga las dependencias con ayuda de la imagen de golang
RUN go mod download

# Copia el código fuente
COPY . .

# Compila el binario
RUN set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o /binario .

# Crea una imagen mínima para ejecutar el binario
FROM debian:sid-slim

# Copia el binario compilado desde la fase de construcción
COPY --from=builder ./binario ./binario

# Asegúrate de que el binario tiene permisos de ejecución
RUN chmod +x ./binario

# Ejecuta el binario
CMD ["./binario", "-key=/home/uqcloud/.ssh/id_rsa"] 