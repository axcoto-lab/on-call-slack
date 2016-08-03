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
	//"io"
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

type SlackMessageResponse struct {
	Text            string `json:"text"`
	ResponseType    string `json:"response_type"`
	ReplaceOriginal bool   `json:"replace_original"`
	DeleteOriginal  bool   `json:"delete_original"`
}

type SlackMessageAction struct {
	Actions []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
	Team       struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTs        string `json:"action_ts"`
	MessageTs       string `json:"message_ts"`
	AttachmentID    string `json:"attachment_id"`
	Token           string `json:"token"`
	OriginalMessage struct {
		Text        string `json:"text"`
		Username    string `json:"username"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			CallbackID string `json:"callback_id"`
			Fallback   string `json:"fallback"`
			Text       string `json:"text"`
			ID         int    `json:"id"`
			Color      string `json:"color"`
			Actions    []struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Text    string `json:"text"`
				Type    string `json:"type"`
				Value   string `json:"value"`
				Style   string `json:"style"`
				Confirm struct {
					Text        string `json:"text"`
					Title       string `json:"title"`
					OkText      string `json:"ok_text"`
					DismissText string `json:"dismiss_text"`
				} `json:"confirm,omitempty"`
			} `json:"actions"`
		} `json:"attachments"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Ts      string `json:"ts"`
	} `json:"original_message"`
	ResponseURL string `json:"response_url"`
}

type hash map[string]string
type Storage struct {
	store map[string]hash
}

func (s *Storage) Add(m SlackMessageAction) {
	if s.store[m.CallbackID] == nil {
		s.store[m.CallbackID] = make(hash) // store max 200 queue item
	}

	for _, action := range m.Actions {
		s.store[m.CallbackID][action.Name] = fmt.Sprintf("%s|%s", m.Channel.ID, action.Value)
	}
}

func (s *Storage) GetJSON(provider string) []byte {
	if s.store[provider] == nil {
		return []byte{}
	}

	out, err := json.Marshal(s.store[provider])
	if err != nil {
		return []byte{}
	}
	s.store[provider] = nil
	return out
}

func initStorage() *Storage {
	s := Storage{
		store: make(map[string]hash),
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

	form.Add("client_id", s.Id)
	form.Add("client_secret", s.Secret)

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(form.Encode())))
	req.Header["Content-Type"] = []string{" application/x-www-form-urlencoded"}

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
	router.HandleFunc("/message/{provider}", MessageHandle)
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

	fmt.Fprintln(w, "Get code"+code)

	response := slackClient.request("/oauth.access", map[string]string{
		"code": code,
	})

	fmt.Fprintln(w, string(response))
}

func SlackButtonHandle(w http.ResponseWriter, r *http.Request) {
	var message SlackMessageAction

	payload := r.FormValue("payload")
	//body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	log.Println(string(payload))

	err := json.Unmarshal([]byte(payload), &message)
	if err != nil {
		fmt.Fprintf(w, "Cannot read body JSON %v", err)
		return
	}

	log.Println(message)

	if err != nil {
		fmt.Fprintf(w, "Cannot parse JSON %v", err)
		return
	}

	Store.Add(message)
	/*
		j := SlackMessageResponse{
			Text:            "Queued action. Will process",
			ResponseType:    "in_channel",
			ReplaceOriginal: false,
			DeleteOriginal:  false,
		}
		response, err := json.Marshal(&j)
		if err != nil {
			fmt.Fprintln(w, "error")
			return
		}
	*/

	fmt.Fprintf(w, "Action is queued for process")
}

func MessageHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	message := Store.GetJSON(provider)

	fmt.Fprintf(w, "%s", message)
}
