package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const alwaysMonkeyingAround = 1.0
const neverMonkeyAround = 0.0
const cannedResponse = "hello, world"

func TestItMonkeysWithStatusCodesAndBodies(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = alwaysMonkeyingAround
	monkeyBehaviour.Status = http.StatusNotFound
	monkeyBehaviour.Body = "hello, monkey"

	testServer, request := makeTestServerAndRequest()

	monkeyServer := NewMonkeyServer(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	if w.Code != monkeyBehaviour.Status {
		t.Error("Server shouldve returned a 404 because of monkey override")
	}

	if w.Body.String() != monkeyBehaviour.Body {
		t.Error("Server should've returned a different body because of monkey override")
	}
}

func TestItReturnsGarbage(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = alwaysMonkeyingAround
	monkeyBehaviour.Garbage = 1984

	testServer, request := makeTestServerAndRequest()

	monkeyServer := NewMonkeyServer(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	if len(w.Body.Bytes()) != monkeyBehaviour.Garbage {
		t.Error("Server shouldve returned garbage")
	}
}

func TestItDoesntMonkeyAroundWhenFrequencyIsNothing(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = neverMonkeyAround
	monkeyBehaviour.Body = "blah blah"

	testServer, request := makeTestServerAndRequest()

	monkeyServer := NewMonkeyServer(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	if w.Body.String() != cannedResponse {
		t.Error("Server shouldn't have been monkeyed with ")
	}
}

func makeTestServerAndRequest() (*httptest.Server, *http.Request) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cannedResponse))
	}))
	request, _ := http.NewRequest("GET", server.URL, nil)

	return server, request
}