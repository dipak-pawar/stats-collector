package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/dipak-pawar/stats-collector/db"
	"github.com/dipak-pawar/stats-collector/config"
	"github.com/dipak-pawar/stats-collector/controller"
)

//go:generate go-bindata -pkg db -o db/bindata.go -nocompress db/migrations/

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v1/status", statusHandler).Methods("GET")
	return r
}

func main() {
	conf := config.Postgres.String()

	// migrate db
	DB := db.Connect(conf)
	db.MigrateDatabase(DB)
	if err := DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}

	r := newRouter()
	DB = db.Connect(conf)
	controller.Register(r, DB)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("something went wrong", err)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok!!!\n")
}
