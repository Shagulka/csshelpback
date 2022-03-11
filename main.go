package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type user struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type userwithtype struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

var users = make(map[string]string)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var checkuser user
		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &checkuser)
		logpass := checkuser.Login + ":" + checkuser.Password
		typeofuser := users[logpass]
		var resp userwithtype
		resp.Login = checkuser.Login
		resp.Password = checkuser.Password
		if typeofuser == "" {
			typeofuser = "not_a_user"
		}
		resp.Type = typeofuser
		file, _ := json.MarshalIndent(resp, "", " ")
		_, _ = fmt.Fprintf(w, "%s", file)
	case "POST":
		var newuser userwithtype
		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &newuser)
		logpass := newuser.Login + ":" + newuser.Password
		users[logpass] = newuser.Type
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	http.HandleFunc("/", handleUsers)
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}
