package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func indexController(w http.ResponseWriter, r *http.Request) {

	type indexContextData struct {
		Events []Event
		Today  time.Time
	}

	theEvents, err := getAllEvents()
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	contextData := indexContextData{
		Events: theEvents,
		Today:  time.Now(),
	}

	tmpl["index"].Execute(w, contextData)
}

func createEventController(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data from the POST request
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}

		// Extract form values
		title := r.FormValue("title")
		location := r.FormValue("location")
		image := r.FormValue("image")
		dateStr := r.FormValue("date")

		// Parse the date
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format. Please use YYYY-MM-DD.", http.StatusBadRequest)
			return
		}

		// Create new event
		newEvent := Event{
			Title:    title,
			Location: location,
			Image:    image,
			Date:     date,
		}

		// Add the event to the list of all events
		addEvent(newEvent)

		// Redirect or render success page
		http.Redirect(w, r, "/events", http.StatusSeeOther)
	} else {
		// Render the form if the request is a GET request
		tmpl["create"].Execute(w, nil)
	}
}

func accessEventController(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/events/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	contextEvent, exists := getEventByID(id)
	if !exists {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	tmpl["access"].Execute(w, contextEvent)
}