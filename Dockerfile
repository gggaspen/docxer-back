# Usar una imagen base de Go
FROM golang:1.21-alpine AS build

# Establecer el directorio de trabajo
WORKDIR /server

# Copiar los archivos go.mod y go.sum
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el resto de los archivos
COPY . .

# Compilar la aplicación
RUN go build -o blank .

# Usar Alpine para la etapa de ejecución
FROM alpine:latest

WORKDIR /root/

# Copiar el binario desde la etapa de compilación
COPY --from=build /server/blank .

# Copiar el directorio config
COPY --from=build /server/config ./config

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./blank"]
