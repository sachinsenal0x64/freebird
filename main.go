package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

type constantHandler struct {
	Val string
}

type env struct {
	apiKey string
	apiUrl string
}

func NewEnv() (*env, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	apiKey := os.Getenv("API_KEY")

	apiUrl := os.Getenv("SERVER_HOST")

	return &env{
		apiKey: apiKey,
		apiUrl: apiUrl,
	}, nil
}

func (ch constantHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_, _ = rw.Write([]byte(ch.Val))
}

func createMux(version string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/login", constantHandler{version + " login"})
	mux.Handle("/submit", constantHandler{version + " submit"})
	mux.Handle("/read", constantHandler{version + " read"})
	return mux
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {

	// Create a new Env instance
	env, err := NewEnv()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Use the Env instance
	log.Println("API_KEY:", env.apiKey)

	ip := GetOutboundIP()

	addressFormat := fmt.Sprintf("%s:%d", ip, 9000)

	w := os.Stderr

	slog.SetDefault(slog.New(tint.NewHandler(colorable.NewColorable(w), &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})))

	topMux := http.NewServeMux()
	topMux.Handle("/v1/", http.StripPrefix("/v1", createMux("v1")))
	topMux.Handle("/v2/", http.StripPrefix("/v2", createMux("v2")))

	slog.Info("Starting server", "addr", fmt.Sprintf("http://%s", addressFormat))

	http.ListenAndServe(addressFormat, topMux)

}
