package exampleMessagingBackend

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Message has sender, recipient, content, and epoch time
type Message struct {
	to      string
	from    string
	message string
	time    int
}

// Messages holds the messages themselves
type Messages []Message

// Convo holds pointers to the messages
type Convo struct {
	participant string
	history     []*Message
}

// User is a struct for holding each user's conversations. The conversations themselves have pointers to the messages held in a Messages struct
type User struct {
	userName string
	convos   map[string]Convo
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

func setMessage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Set message\n"))
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get message\n"))
}
