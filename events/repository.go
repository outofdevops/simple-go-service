package events

import (
	"database/sql"
	"log"
)

/*
Repository type
*/
type Repository struct {
	db *sql.DB
}

/*
NewEventRepo Given a DB returns a repository
*/
func NewEventRepo(db *sql.DB) *Repository {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	setupTable(db)

	return &Repository{db: db}
}

/*
GetEvent given an eventId gets an event
*/
func (r *Repository) GetEvent(id string) (Event, error) {
	var e Event
	e.ID = id
	err := r.db.QueryRow(selectEventQuery,
		id).Scan(&e.Name, &e.State)

	return e, err
}

/*
CreateEvent stores an event into the events table
*/
func (r *Repository) CreateEvent(e Event) (Event, error) {
	err := r.db.QueryRow(createEvent, e.Name, e.State).Scan(&e.ID)

	return e, err
}

/*
GetEvents returns a list of events
*/
func (r *Repository) GetEvents(start, count int) ([]Event, error) {
	rows, err := r.db.Query(
		selectEventsQuery,
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []Event{}

	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Name, &e.State); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func setupTable(db *sql.DB) {
	runQuery(db, createExtension)
	runQuery(db, createEventTable)
}

func runQuery(db *sql.DB, query string) {
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

/*Queries*/
const createExtension = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`

const createEventTable = `CREATE TABLE IF NOT EXISTS events
(
    id uuid DEFAULT uuid_generate_v4 (),
    name  TEXT NOT NULL,
    state TEXT NOT NULL,
    CONSTRAINT events_pkey PRIMARY KEY (id)
)`

const createEvent = `INSERT INTO events(name, state) VALUES($1, $2) RETURNING id`

const selectEventQuery = `SELECT name, state FROM events WHERE id=$1`

const selectEventsQuery = `SELECT id, name,  state FROM events LIMIT $1 OFFSET $2`
