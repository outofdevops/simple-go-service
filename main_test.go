package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"event-service/events"
)

var a App

func TestMain(m *testing.M) {
	log.Println("Calling init")
	a.Init("localhost",
		"5432",
		"testdb",
		"testdb",
		"testdb")

	log.Println("Table created")
	code := m.Run()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/events", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array - Actual %s", body)
	}
}

func TestEventNotFound(t *testing.T) {
	//Given
	req, _ := http.NewRequest("GET", "/events/A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A11", nil)

	//When
	response := executeRequest(req)

	//Then
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Event not found" {
		t.Errorf("Expected 'error' to be set to 'Event not found' - Actual '%s'", m["error"])
	}
}

func TestCreateEvent(t *testing.T) {
	//Given
	var payload = []byte(`{"name": "registration", "state": "Accepted"}`)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	//When
	response := executeRequest(req)

	//Then
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "registration" {
		t.Errorf("Expected name to be 'registration' - Actual '%v'", m["type"])
	}

	if m["state"] != "Accepted" {
		t.Errorf("Expected state to be 'Accepted' - Actual '%v'", m["value"])
	}
}

func TestGetEvent(t *testing.T) {
	//Given
	id := addEvent()

	//When
	req, _ := http.NewRequest("GET", "/events/"+id, nil)
	response := executeRequest(req)

	//Then
	checkResponseCode(t, http.StatusOK, response.Code)
}

func addEvent() (id string) {
	e, _ := a.eventRepo.CreateEvent(events.Event{Name: "Registration", State: "Accepted"})
	return e.ID
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected %d - Actual %d\n", expected, actual)
	}
}
