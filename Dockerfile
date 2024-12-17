# Etapa de construcción
FROM golang:1.23.4 AS builder

# Configurar el directorio de trabajo
WORKDIR /app

# Copiar los archivos del proyecto al contenedor
COPY . .

# Descargar dependencias y construir la aplicación de forma estática
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Etapa de ejecución con Alpine
FROM alpine:latest

# Crear directorio de trabajo en la imagen final
WORKDIR /root/

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/main .

# Copiar el archivo .env al contenedor
COPY .env /root/.env

# Asegurarse de que el binario tiene permisos de ejecución
RUN chmod +x main

# Definir el puerto en el que la app escucha
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
