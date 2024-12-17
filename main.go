package main

import (
	"log"
	"net/http"
	"notificaciones/handlers"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// Cargar las variables de entorno una sola vez
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	r := mux.NewRouter()

	// Ruta para enviar correos
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/send-email", handlers.SendEmailHandler).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡Servidor funcionando correctamente!"))
	})

	// Iniciar servidor
	port := os.Getenv("PORT")

	log.Println("Servidor escuchando en el puerto ", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, r); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
