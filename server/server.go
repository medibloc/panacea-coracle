package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func Run() {
	router := mux.NewRouter()
	router.HandleFunc("/validate-data/{dealId}", handleRequest).Methods(http.MethodPost)

	server := &http.Server{
		Handler:      router,
		Addr:         os.Getenv("HTTP_LADDR"), // listener address: <IPaddr>:<port>
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("👻 Data Validator Server Started 🎃")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
