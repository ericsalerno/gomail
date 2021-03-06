package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerCreate(t *testing.T) {
	c := Configuration{
		Port: 8081,
	}

	b := Service{
		Config: &c,
	}

	if b.Config.Port != 8081 {
		t.Fatal("Missing config value (somehow)")
	}
}

func TestServerRegexes(t *testing.T) {
	b := Service{}

	b.buildRegularExpressions()

	// [a-zA-Z0-9\.-]+
	if b.validEmailHost.MatchString("test-HOSTval124id.net") == false {
		t.Fatal("Invalid email host regex (plain)!")
	}

	if b.validEmailHost.MatchString("localhost") == false {
		t.Fatal("Invalid email host regex (local)!")
	}

	if b.validEmailHost.MatchString("87dfa@$22") == true {
		t.Fatal("Invalid email host regex (invalid characters)!")
	}

	if b.validEmailHost.MatchString("87dfa@$22.lol") == true {
		t.Fatal("Invalid email host regex (invalid characters)!")
	}

	// ^\d{1-3}\.\d{1-3\.\d{1-3}\.\d{1-3}$
	if b.validEmailHostIP.MatchString("12.14.156.255") == false {
		t.Fatal("Invalid host IP regex (plain)!")
	}

	if b.validEmailHostIP.MatchString("12.14.156") == true {
		t.Fatal("Invalid host IP regex (missing octet)!")
	}

	if b.validEmailHostIP.MatchString("12.14.156.123.554") == true {
		t.Fatal("Invalid host IP regex (extra octet)!")
	}

	// [a-zA-Z0-9!#$%&'*+/=\?^_\{\}|~\.-]+
	if b.validEmailUser.MatchString("geezer123") == false {
		t.Fatal("Invalid user regex (plain)!")
	}

	if b.validEmailUser.MatchString("") == true {
		t.Fatal("Invalid user regex (empty)!")
	}

	if b.validEmailUser.MatchString("b&l*a4{#'1_+r=gh?12}3~-.|^as%df$!") == false {
		t.Fatal("Invalid user regex (special characters)!")
	}
}

func TestGetValidResponseOutput(t *testing.T) {
	r := request{}
	r.inputEmail = "testing@thing.com"
	r.inputHost = "thing.com"
	r.inputUser = "testing"

	s := Service{}

	response := s.getResponseOutput(&r, true)

	if response.Status != 200 {
		t.Fatal("Invalid response code!")
	}

	if response.Message != "OK" {
		t.Fatal("Invalid message!")
	}

	if response.Email != r.inputEmail {
		t.Fatal("Invalid email address!")
	}

	if response.Valid != true {
		t.Fatal("Invalid response valid!")
	}

	if response.Host != r.inputHost {
		t.Fatal("Invalid response host!")
	}

	if response.User != r.inputUser {
		t.Fatal("Invalid response user!")
	}
}

func TestGetResponseError(t *testing.T) {
	r := request{}
	r.inputEmail = "testing@thing.com"
	r.inputHost = "thing.com"
	r.inputUser = "testing"

	s := Service{}

	response := s.getResponseError(&r, "Something or other")

	if response.Status != 500 {
		t.Fatal("Invalid response code!")
	}

	if response.Message != "Something or other" {
		t.Fatal("Invalid message!")
	}

	if response.Email != r.inputEmail {
		t.Fatal("Invalid email address!")
	}

	if response.Valid != false {
		t.Fatal("Invalid response valid!")
	}

	if response.Host != r.inputHost {
		t.Fatal("Invalid response host!")
	}

	if response.User != r.inputUser {
		t.Fatal("Invalid response user!")
	}
}

func TestServeHTTP(t *testing.T) {
	s := Service{
		Config: &Configuration{
			Port: 8081,
		},
	}
	s.buildRegularExpressions()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ServeHTTP)

	body := []byte(`email=test@test.com`)
	req, _ := http.NewRequest("POST", "http://localhost/", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler returned invalid status! Got %v want %v %v", status, http.StatusOK, rr.Body)
	}
}

func TestFailureServeHTTP(t *testing.T) {
	s := Service{
		Config: &Configuration{
			Port: 8081,
		},
	}
	s.buildRegularExpressions()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ServeHTTP)

	body := []byte(`email=blargh!`)
	req, _ := http.NewRequest("POST", "http://localhost/", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("Handler returned invalid status! Got %v want %v %v", status, http.StatusInternalServerError, rr.Body)
	}

	req, _ = http.NewRequest("POST", "http://localhost/", nil)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Fatalf("Handler returned invalid status! Got %v want %v %v", status, http.StatusInternalServerError, rr.Body)
	}
}
