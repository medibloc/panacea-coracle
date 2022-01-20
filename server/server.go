package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func Run() {
	SetConfig()

	validateDataHandler, err := NewValidateDataHandler()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Handle("/validate-data/{dealId}", validateDataHandler)

	server := &http.Server{
		Handler:      router,
		Addr:         os.Getenv("HTTP_LADDR"), // listener address: <IPaddr>:<port>
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("👻 Data Validator Server Started 🎃")
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
