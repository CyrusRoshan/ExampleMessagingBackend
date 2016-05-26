package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Response struct {
	Status string `json:"status"`
}

var cases = []struct {
	userNameA, userNameB, message string
}{
	{"a", "b", "message a->b 1"},
	{"b", "a", "message a<-b 2"},
	{"a", "c", "message a->c 1"},
	{"c", "b", "message b<-c 1"},
	{"self", "self", "message self 1"},
	{"self", "self", "message self 2"},
	{"c", "self", "message c->self 1"},
	{"self", "c", "message self->c 2"},
	{"123!@#!#", "ˆ∆ˆøˆ¡™£∆∆ˆÔ¨ÓŒ„‡€⁄°‡⁄€°·™£ºªˆª˚∆", "message nonstandard characters 1"},
	{"ˆ∆ø∆ˆø√˜˜∫", "•™ª¡™ª•£¡º¡º™øˆ∆Ò˜Â˜¯Â˛˜¯˛˜ÔÅ˜˜ÓÔ", "message nonstandard characters 2"},
}

var caseTimes []int64

func TestStartup(t *testing.T) {
	t.Log("Starting up messaging server")
	go main()
}

func TestSending(t *testing.T) {
	for _, testCase := range cases {
		time.Sleep(1000 * time.Millisecond)              // sleep 1 second so we get different timestamps
		caseTimes = append(caseTimes, time.Now().Unix()) //timestamp before test case

		// convert test case to proper message format
		var caseMessage Message
		caseMessage.To = testCase.userNameA
		caseMessage.From = testCase.userNameB
		caseMessage.Message = testCase.message

		// convert test case to byte[]
		caseJSON, err := json.Marshal(caseMessage)
		if err != nil {
			t.Error(err)
		}

		// create request, using JSON string of test case
		reader := strings.NewReader(string(caseJSON))
		t.Log(string(caseJSON))
		request, err := http.NewRequest("POST", "http://localhost:3001/", reader)
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf(err.Error())
		}

		// execute request and get request response
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			t.Errorf(err.Error())
		}

		// Unmarshal request response
		var response Response
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Errorf(err.Error())
		}

		// check if response was successful, if not, this test case failed
		if response.Status != "success" {
			t.Errorf(err.Error())
		}
	}
}

/*func TestReading(t *testing.T) {
	for _, testCase := range cases {
		// curl the /UserNameA/UserNameA/FromTimeStamp, where FromTimeStamp = 0 so we can get all messages between users
		resp, err := http.Get("http://localhost:3001/" + url.QueryEscape(testCase.userNameA) + "/" + url.QueryEscape(testCase.userNameB) + "/0")
		if err != nil {
			t.Errorf(err.Error())
		}

		// Unmarshal JSON to array of messages
		var response []Message
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Log(response)

		// compare the results to those stored in the server's memory
		// linear search because there are like 10 test cases
		participants := []string{testCase.userNameA, testCase.userNameB}
		for i, thisUser := range participants {
			otherUser := participants[1-i] // turns 1 to 0 and vice versa. Kind of hacky, though

			for _, message := range Users[thisUser].Convo[otherUser] {
				messageFound := false
				for _, recievedMessage := range response {
					if recievedMessage == *message {
						messageFound = true
						break
					}
				}
				if !messageFound {
					t.Error("A message between " + (*message).From + " and " + (*message).To + " was not found.")
				}
				// could break here if thisUser == otherUser, but it doesn't matter for test case pass/fail, only very minor performance
			}
		}
	}
}*/

func TestTimedReading(t *testing.T) {
	for i, testCase := range cases {
		// curl the /UserNameA/UserNameA/FromTimeStamp, where FromTimeStamp = 0 so we can get all messages between users
		t.Log(i)
		t.Log("http://localhost:3001/" + url.QueryEscape(testCase.userNameA) + "/" + url.QueryEscape(testCase.userNameB) + "/" + strconv.FormatInt(caseTimes[i], 10))
		resp, err := http.Get("http://localhost:3001/" + url.QueryEscape(testCase.userNameA) + "/" + url.QueryEscape(testCase.userNameB) + "/" + strconv.FormatInt(caseTimes[i], 10))
		if err != nil {
			t.Errorf(err.Error())
		}

		// Unmarshal JSON to array of messages
		var response []Message
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Log(response)

		// compare the results to those stored in the server's memory
		// linear search because there are like 10 test cases
		participants := []string{testCase.userNameA, testCase.userNameB}
		for i, thisUser := range participants {
			otherUser := participants[1-i] // turns 1 to 0 and vice versa. Kind of hacky, though

			for _, message := range Users[thisUser].Convo[otherUser] {
				messageFound := false
				for _, recievedMessage := range response {
					if recievedMessage == *message {
						messageFound = true
						break
					}
				}
				if !messageFound {
					t.Error("A message between " + (*message).From + " and " + (*message).To + " was not found.")
				}
				// could break here if thisUser == otherUser, but it doesn't matter for test case pass/fail, only very minor performance
			}
		}
	}
}
