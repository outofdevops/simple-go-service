package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"event-service/events"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

/*
App struct
*/
type App struct {
	router    *mux.Router
	eventRepo *events.Repository
}

/*
Init function
*/
func (a *App) Init(host, port, user, password, dbname string) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	a.eventRepo = events.NewEventRepo(db)

	a.router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.router.HandleFunc("/events", a.getEvents).Methods("GET")
	a.router.HandleFunc("/events", a.createEvent).Methods("POST")
	a.router.HandleFunc("/events/{id}", a.getEvent).Methods("GET")
}

/*
Run function
*/
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.router))
}
