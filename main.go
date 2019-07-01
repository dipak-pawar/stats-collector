package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//go:generate go-bindata -pkg db -o db/bindata.go -nocompress db/migrations/

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v1/status", statusHandler).Methods("GET")
	return r
}

func main() {
	r := newRouter()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("something went wrong", err)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok!!!\n")
}
