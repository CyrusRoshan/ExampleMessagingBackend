package exampleMessagingBackend

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func setMessage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Set message\n"))
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get message\n"))
}

func main() {
	// Create router
	r := mux.NewRouter()

	// POST request for sending messages from user to user
	r.HandleFunc("/", YourHandler).
		Methods("GET").
		HeadersRegexp("Content-Type", "application/(text|json)")

	// GET request for getting messages between two users after fromTimeStamp
	r.HandleFunc("/{userNameA}/{userNameB}/{fromTimeStamp:[0-9]+}", ProductHandler).
		Methods("POST")

	// Listen to port, so we access this at http://localhost:3001
	http.ListenAndServe(":3001", r)
}
