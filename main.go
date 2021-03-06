package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Message has sender, recipient, content, and epoch time
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
	At      int64  `json:"at,omitempty"`
}

// User is a struct for holding each user's conversations. The conversations themselves have pointers to the messages held in a Messages struct
type User struct {
	UserName string
	Convo    map[string][]*Message
}

// Messages holds the messages themselves
var Messages []Message

// Users holds all of the users, and is the intended way to access message data for a user
var Users = make(map[string]User)

func main() {
	// Create router
	r := mux.NewRouter()

	// POST request for sending messages from user to user
	r.HandleFunc("/", setMessage).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(text|json)")

	// GET request for getting messages between two users after FromTimeStamp
	r.HandleFunc("/{UserNameA}/{UserNameB}/{FromTimeStamp:[0-9]+}", getMessage).
		Methods("GET")

	// Catch all other requests and send an error
	r.NotFoundHandler = http.HandlerFunc(errorPage)

	// Listen to port, so we access this at http://localhost:3001
	http.ListenAndServe(":3001", r)
}

// for sending messages between two users
func setMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //reply with JSON success or failure

	// Read sent data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	// Create new message and fill in current (serverside) time
	newMessage := *new(Message)
	newMessage.At = time.Now().Unix()

	// Fill newMessage with sent data
	err = json.Unmarshal(body, &newMessage)
	if err != nil {
		panic(err.Error())
	}

	if newMessage.To == "" || newMessage.From == "" { // actually have to impose this restriction on usernames
		w.Write([]byte(`{"status": "failure"}` + "\n"))
		return
	}

	Messages = append(Messages, newMessage) // Add the new message to the end of the Messages array

	// Add the message contents to both convos.
	// This is NOT REQUIRED, we can get by just storing the pointers in one array where the first user
	// is the alphabetically first one, however this method makes it easier to add other features on later
	// e.g. where on fb messenger you can delete a message and it's not visible to you but is still visible to others
	participants := []string{newMessage.To, newMessage.From}
	for i, thisUser := range participants {
		otherUser := participants[1-i] // turns 1 to 0 and vice versa. Kind of hacky, though

		// create user if user doesn't exist
		_, thisUserExists := Users[thisUser]
		if !thisUserExists {
			Users[thisUser] = User{thisUser, make(map[string][]*Message)}
		}

		// create conversation with other user if there was no prior convo
		_, otherUserExists := Users[otherUser]
		if !otherUserExists {
			Users[thisUser].Convo[otherUser] = []*Message{}
		}

		// add pointer to the message inside
		Users[thisUser].Convo[otherUser] = append(Users[thisUser].Convo[otherUser], &newMessage)

		// Similar to how on fb messenger, sending a message to yourself will only show the send message,
		// and you will never see any recieved messages from yourself
		if thisUser == otherUser {
			break
		}
	}

	w.Write([]byte(`{"status": "success"}` + "\n")) // always nice to know if you have a spotty internet connection
}

// for reading messages between two users after a particular time
func getMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // send content in JSON

	vars := mux.Vars(r)
	userNameA := vars["UserNameA"]
	userNameB := vars["UserNameB"]
	fromTimeStamp, err := strconv.ParseInt(vars["FromTimeStamp"], 10, 64)
	if err != nil {
		panic(err.Error())
	}

	// Construct the JSON that we're going to send back. We don't need the array to be sorted.
	jsonString := "["
	for _, currMessage := range Users[userNameA].Convo[userNameB] {
		if (*currMessage).At > fromTimeStamp { // If the message is after the the time sent, we can send the message
			data, err := json.Marshal(*currMessage)
			if err != nil {
				panic(err.Error())
			}
			jsonString += string(data)
			jsonString += ","
		}
	}
	// Remove the last comma (if we actually added a message and comma)
	if jsonString != "[" {
		jsonString = jsonString[:len(jsonString)-1]
	}
	jsonString += "]"

	w.Write([]byte(jsonString + "\n"))
}

// for sendiing error messages for every request type not previously defined
func errorPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // send content in JSON
	w.Write([]byte(`{"status": "nonexistent page"}` + "\n"))
}
