package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"gopkg.in/gomail.v2"
)

type TKey struct {
	Key string `json:"key"`
}

type TResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var key TKey
	json.NewDecoder(r.Body).Decode(&key)
	var serverKey = os.Getenv("SERVER_KEY")
	if key.Key == serverKey {
		tokenString := createToken("Authorized")
		w.WriteHeader(http.StatusOK)
		response := TResponse{Token: tokenString}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid credentials")
	}
}

type EmailRequest struct {
	Destinatario string `json:"destinatario"`
	Cuerpo       string `json:"cuerpo"`
	Asunto       string `json:"asunto"`
}

// WaitGroup para esperar a que todas las goroutines terminen (global si es necesario)
var wg sync.WaitGroup

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verificar autorización
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	// Leer el cuerpo de la solicitud
	var emailReq EmailRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo solicitud", http.StatusBadRequest)
		return
	}

	// Parsear el JSON
	if err := json.Unmarshal(body, &emailReq); err != nil {
		http.Error(w, "Error procesando JSON", http.StatusBadRequest)
		return
	}

	// Incrementar el contador y lanzar la goroutine
	wg.Add(1)
	go sendEmailAsync(emailReq)

	// Respuesta inmediata al cliente
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Correo recibido, se enviará pronto"))
}

func createToken(permission string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"permission": permission,
			"exp":        time.Now().Add(time.Hour * 24).Unix(),
		})

	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error al crear el token:", err)
		return ""
	}

	return tokenString
}

func verifyToken(tokenString string) error {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func sendEmailAsync(emailReq EmailRequest) {
	defer wg.Done() // Decrementar el contador al finalizar la goroutine

	// Obtener las variables de entorno localmente
	sender := os.Getenv("MAINGUN_SENDER")
	host := os.Getenv("MAILGUN_HOST")
	password := os.Getenv("MAILGUN_PASSWORD")
	var port = 587

	// Configuración del correo
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", emailReq.Destinatario)
	m.SetHeader("Subject", emailReq.Asunto)
	m.SetBody("text/html", emailReq.Cuerpo)

	// Configuración del servidor SMTP de Mailgun
	d := gomail.NewDialer(host, port, sender, password)

	// Envío del correo
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error al enviar el correo:", err)
		return
	}

	fmt.Println("Correo enviado exitosamente con Mailgun y Gomail.")
}
