package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Defining all constants in-order to avoid typo error
const (
	HostPort        = ":8081"
	ContentType     = "content-type"
	ApplicationJson = "application/json"
	URL             = "http://localhost:8081/"
	Msg             = "message"
	Msgs            = "messages"
	Users           = "users"
)

// Message has user, text, timestamp information
type Message struct {
	MessageReceivedAt string `json:"timestamp"`
	User              string `json:"user"`
	Text              string `json:"text"`
}

type messageHandlers struct {
	sync.Mutex
	messages []Message
	users    map[string]bool
}

func (m *messageHandlers) getMessages(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	messages := make([]Message, len(m.messages))

	m.Lock()
	i := 0
	for _, message := range m.messages {
		messages[i] = message
		i++
	}
	m.Unlock()

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].MessageReceivedAt < messages[j].MessageReceivedAt
	})

	var msgBytes []byte
	var err error

	if len(messages) > 100 {
		msgBytes, err = json.Marshal(
			map[string][]Message{Msgs: messages[:100]},
		)

	} else {
		msgBytes, err = json.Marshal(
			map[string][]Message{Msgs: messages},
		)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(msgBytes)
}

func (m *messageHandlers) postMessage(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ct := r.Header.Get(ContentType)
	if ct != ApplicationJson {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need %q %q, but got '%s'", ContentType, ApplicationJson, ct)))
		return
	}

	var message Message
	err = json.Unmarshal(bodyBytes, &message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	message.MessageReceivedAt = fmt.Sprintf("%d", time.Now().UnixNano())
	m.Lock()
	m.messages = append(m.messages, message)
	m.users[message.User] = true
	defer m.Unlock()
}

func NewMessageHandlers() *messageHandlers {
	return &messageHandlers{
		messages: []Message{},
		users:    map[string]bool{},
	}

}

func (m *messageHandlers) getUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	users := make([]string, len(m.users))

	m.Lock()
	i := 0
	for name := range m.users {
		users[i] = name
		i++
	}
	m.Unlock()

	jsonBytes, err := json.Marshal(
		map[string][]string{Users: users},
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func main() {
	messageHandlers := NewMessageHandlers()

	srv := &http.Server{
		IdleTimeout: 5 * time.Second,
		Addr:        HostPort, //static port for the external load balancer
	}

	// Route handles & endpoints
	http.HandleFunc("/"+Msgs, messageHandlers.getMessages)
	http.HandleFunc("/"+Msg, messageHandlers.postMessage)
	http.HandleFunc("/"+Users, messageHandlers.getUsers)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

}
