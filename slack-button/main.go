package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"bytes"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"os"
	"time"
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
	actions          []*SlackAction `json:"actions"`
	callback_id      string         `json:"callback_id"`
	team             *SlackTeam     `json:"team"`
	channel          *SlackKV       `json:"channel"`
	user             *SlackKV       `json:"user"`
	token            string         `json:"token"`
	original_message string         `json:"original_message"`
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

type SlackClient struct {
	Id         string
	Secret     string
	Endpoint   string
	httpClient *http.Client
}

type SlackRequestData map[string]string

func (s *SlackClient) request(path string, data SlackRequestData) []byte {
	apiUrl := fmt.Sprintf("%s/%s", s.Endpoint, path)

	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(form.Encode())))
	if err != nil {
		return nil
	}

	resp, err := s.httpClient.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return body
}

func initSlackClient(id, secret string) *SlackClient {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	c := SlackClient{
		Id:         id,
		Secret:     secret,
		Endpoint:   "https://slack.com/api",
		httpClient: netClient,
	}
	return &c
}

var (
	Store       *Storage
	slackClient *SlackClient
)

func main() {
	Store = initStorage()
	slackClient = initSlackClient(os.Getenv("client_id"), os.Getenv("client_secret"))

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/button", SlackButtonHandle)
	router.HandleFunc("/message", MessageHandle)
	router.HandleFunc("/incoming", OAuthHandle)

	bind := "127.0.0.1:2830"
	fmt.Println("Bind to http://" + bind)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Fatal(http.ListenAndServe(bind, loggedRouter))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func OAuthHandle(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		fmt.Fprintln(w, "Invalid code")
		return
	}

	fmt.Fprintln(w, "OK")
}

func SlackButtonHandle(w http.ResponseWriter, r *http.Request) {
	var message SlackMessageAction

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	log.Println(string(body))

	if err != nil {
		fmt.Fprintf(w, "Cannot read body JSON %v", err)
		return
	}

	err = json.Unmarshal(body, &message)
	log.Println(message)

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
