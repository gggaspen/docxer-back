# Usar una imagen base de Go
FROM golang:1.23-alpine

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos del proyecto al contenedor
COPY . .

# Descargar dependencias del proyecto (si usas mod)
RUN go mod tidy

# Compilar la aplicación Go
RUN go build -o main .

# Exponer el puerto en el que la app va a correr
EXPOSE 8080

# Comando para ejecutar la app cuando el contenedor inicie
CMD ["./main"]
