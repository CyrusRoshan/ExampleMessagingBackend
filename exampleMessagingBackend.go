package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Message has sender, recipient, content, and epoch time
type Message struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"message"`
	At      int64  `json:"at"`
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
	fmt.Println("Starting...")

	// Create router
	r := mux.NewRouter()

	// POST request for sending messages from user to user
	r.HandleFunc("/", setMessage).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(text|json)")

	// GET request for getting messages between two users after FromTimeStamp
	r.HandleFunc("/{UserNameA}/{UserNameB}/{FromTimeStamp:[0-9]+}", getMessage).
		Methods("GET")

	// Listen to port, so we access this at http://localhost:3001
	http.ListenAndServe(":3001", r)
}

func setMessage(w http.ResponseWriter, r *http.Request) {
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

	Messages = append(Messages, newMessage) // Add the new message to the end of the Messages array

	// Add the message contents to both convos.
	// This is NOT REQUIRED, we can get by just storing the pointers in one array where the first user
	// is the alphabetically first one, however this method makes it easier to add other features on later
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
	}

	w.Write([]byte("LENGTH: " + strconv.Itoa(len(Messages)) + ", DATA: " + newMessage.To + ", " + newMessage.From + ", " + newMessage.Message + ", " + strconv.FormatInt(newMessage.At, 10) + "\n"))
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	UserNameA := vars["UserNameA"]
	UserNameB := vars["UserNameB"]
	FromTimeStamp, err := strconv.ParseInt(vars["FromTimeStamp"], 10, 64)
	jsonString := "["
	if err != nil {
		panic(err.Error())
	}

	for _, currMessage := range Users[UserNameA].Convo[UserNameB] {
		if (*currMessage).At > FromTimeStamp {
			data, err := json.Marshal(*currMessage)
			if err != nil {
				panic(err.Error())
			}
			jsonString += string(data)
			jsonString += ","
		}
	}
	if jsonString != "[" {
		jsonString = jsonString[:len(jsonString)-1]
	}
	jsonString += "]"

	w.Write([]byte(jsonString + "\n"))
}
