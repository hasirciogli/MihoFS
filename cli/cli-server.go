package cli

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// const socketPath = "/etc/mihofs/cli.sock"
const socketPath = "./cli.sock"

func StartCliServer(wg *sync.WaitGroup) {
	// Önceden var olan eski socket dosyasını sil
	err := os.Remove(socketPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error removing old socket:", err)
		os.Exit(1)
	}

	// Unix socket ile dinleyici başlat
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer wg.Done()
	defer listener.Close()

	fmt.Println("Listening on", socketPath)

	// HTTP router oluştur
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/hello", HelloHandler).Methods("GET")

	// HTTP sunucusunu Unix socket ile bağla
	log.Fatal(http.Serve(listener, r))
}

// Ana sayfa handler'ı
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

// Merhaba sayfası handler'ı
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
