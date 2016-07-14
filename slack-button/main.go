package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
)

type SlackAction struct {
	Name   string `json:"name"`
	Action string `json:"action"`
}

type SlackKV struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SlackTeam struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SlackMessageAction struct {
	action           []*SlackAction
	callback_id      string
	team             *SlackTeam
	channel          *SlackKV
	user             *SlackKV
	original_message string
}

type Storage struct {
	store map[string]SlackMessageAction
}

func (s *Storage) Add(m SlackMessageAction) {
	s.store[m.callback_id] = m
}

func (s *Storage) GetJSON() []byte {
	out, err := json.Marshal(s.store)
	if err != nil {
		return nil
	}

	return out
}

func initStorage() *Storage {
	s := Storage{
		store: make(map[string]SlackMessageAction),
	}
	return &s
}

var (
	Store *Storage
)

func main() {
	Store = initStorage()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/button", SlackButtonHandle)
	router.HandleFunc("/message", MessageHandle)

	bind := "127.0.0.1:2830"
	fmt.Println("Bind to http://" + bind)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Fatal(http.ListenAndServe(bind, loggedRouter))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func SlackButtonHandle(w http.ResponseWriter, r *http.Request) {
	var message SlackMessageAction

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Fprintf(w, "Cannot read body JSON %v", err)
		return
	}

	err = json.Unmarshal(body, &message)

	if err != nil {
		fmt.Fprintf(w, "Cannot parse JSON %v", err)
		return
	}

	Store.Add(message)

	fmt.Fprintln(w, "OK")
}

func MessageHandle(w http.ResponseWriter, r *http.Request) {
	message := Store.GetJSON()

	fmt.Fprintf(w, "%s", message)
}
