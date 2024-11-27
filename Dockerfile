# Etapa de construcción
FROM golang:1.21-alpine AS build

WORKDIR /server

# Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del proyecto
COPY . .

# Compilar la aplicación
RUN go build -o blank .

# Etapa de ejecución
FROM alpine:latest

WORKDIR /root/

# Copiar binario y recursos necesarios
COPY --from=build /server/blank .
COPY --from=build /server/config ./config

# Añadir permisos al binario
RUN chmod +x blank

# Exponer el puerto
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./blank"]
