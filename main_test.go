package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
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
}

var caseTimes []int

func TestStartup(t *testing.T) {
	t.Log("Starting up messaging server")
	go main()
}

func TestSending(t *testing.T) {
	for _, testCase := range cases {
		reader := strings.NewReader(`{"to": "` + testCase.userNameA + `", "from": "` + testCase.userNameB + `", "message": "` + testCase.message + `"}`)

		request, err := http.NewRequest("POST", "http://localhost:3001/", reader)
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Errorf(err.Error())
		}

		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			t.Errorf(err.Error())
		}

		var response Response
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Errorf(err.Error())
		}

		if response.Status != "success" {
			t.Errorf(err.Error())
		}
	}
}
