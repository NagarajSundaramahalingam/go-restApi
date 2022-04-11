package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetMethods(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, URL+Msgs, nil)
	req.Header.Add(ContentType, ApplicationJson)
	w := httptest.NewRecorder()

	newMH := NewMessageHandlers()

	newMH.getMessages(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d but got %d\n", http.StatusOK, res.StatusCode)
	}

	if res.Header.Get(ContentType) != ApplicationJson {
		t.Errorf("expected %q but got %q\n", ApplicationJson, res.Header.Get(ContentType))
	}

	req = httptest.NewRequest(http.MethodGet, URL+Users, nil)
	req.Header.Add(ContentType, ApplicationJson)
	w = httptest.NewRecorder()

	newMH.getUsers(w, req)

	res = w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d but got %d\n", http.StatusOK, res.StatusCode)
	}

	if res.Header.Get(ContentType) != ApplicationJson {
		t.Errorf("expected %q but got %q\n", ApplicationJson, res.Header.Get(ContentType))
	}

}

func TestPostAndGetMessages(t *testing.T) {

	payload := strings.NewReader(`
		{"user":"superman", "text":"hello"}
	`)

	newMH := NewMessageHandlers()

	req := httptest.NewRequest(http.MethodPost, URL+Msg, payload)
	req.Header.Add(ContentType, ApplicationJson)
	w := httptest.NewRecorder()

	newMH.postMessage(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d but got %d\n", http.StatusOK, res.StatusCode)
	}

	if res.Header.Get(ContentType) != ApplicationJson {
		t.Errorf("expected %q but got %q\n", ApplicationJson, res.Header.Get(ContentType))
	}


	req = httptest.NewRequest(http.MethodGet, URL+Msgs, nil)
	req.Header.Add(ContentType, ApplicationJson)
	w = httptest.NewRecorder()

	newMH.getMessages(w, req)

	res = w.Result()
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("expected error to be nil but got %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d but got %d\n", http.StatusOK, res.StatusCode)
	}

	if res.Header.Get(ContentType) != ApplicationJson {
		t.Errorf("expected %q but got %q\n", ApplicationJson, res.Header.Get(ContentType))
	}

	if body == nil {
		t.Errorf("expected %q but got %q\n", string(body), res.Header.Get(ContentType))
	}

	var messages map[string][]Message
	json.Unmarshal(body, &messages)

	for _, message := range messages {
		for i := 0; i < len(message); i++ {
			if message[i].User != "superman" || message[i].Text != "hello" {
				t.Errorf("not matched with test data %q\n", message)
			}
		}
	}

	req = httptest.NewRequest(http.MethodGet, URL+Users, nil)
	req.Header.Add(ContentType, ApplicationJson)
	w = httptest.NewRecorder()

	newMH.getUsers(w, req)

	res = w.Result()
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("expected error to be nil but got %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d but got %d\n", http.StatusOK, res.StatusCode)
	}

	if res.Header.Get(ContentType) != ApplicationJson {
		t.Errorf("expected %q but got %q\n", ApplicationJson, res.Header.Get(ContentType))
	}

	if body == nil {
		t.Errorf("expected %q but got %q\n", string(body), res.Header.Get(ContentType))
	}

	var users map[string][]string
	json.Unmarshal(body, &users)

	for _, user := range users {
		for i := 0; i < len(user); i++ {
			if user[i] != "superman" {
				t.Errorf("not matched with test data %q\n", user)
			}
		}
	}
}
