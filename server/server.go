package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"panacea-data-market-validator/types"
	"panacea-data-market-validator/utils"
	"panacea-data-market-validator/validation"
	"time"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	resp := &types.CertificateResponse{}

	dealId := mux.Vars(r)[types.DealIdKey]
	resp.Certificate.DealId = dealId

	// file format check
	data, err := validation.FileReader(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteLogger(w.Write([]byte(err.Error())))
		return
	}

	fmt.Println(data)

	// get deal information from panacea

	// check if data validator is trusted or not

	// validate data (schema check)

	// encrypt and store data

	// sign certificate

	marshaledData, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteLogger(w.Write([]byte(err.Error())))
		return
	}
	w.WriteHeader(http.StatusCreated)
	utils.WriteLogger(w.Write(marshaledData))
}

func Run() {
	router := mux.NewRouter()
	router.HandleFunc("/validate-data/{dealId}", handleRequest).Methods(http.MethodPost)

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("👻 Data Validator Server Started 🎃")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
