/*
Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	cmdPort = "9876"
)

const (
	minUsernameLength = 8
	minPasswordLength = 8
)

type errorResponse struct {
	Error string `json:"error"`
}

func init() {
	flag.StringVar(&cmdPort, "port", cmdPort, "port to listen at")
}

func respondWithError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Println(err)

	errRespJSON, err := json.Marshal(&errorResponse{
		Error: err.Error(),
	})
	if err != nil {
		return
	}
	w.Write(errRespJSON)
}

func validatePassword(pass string) error {
	if len(pass) < minPasswordLength {
		return errors.New("problem registering user with small password")
	}
	return nil
}

func validateUsername(user string) error {
	if len(user) < minUsernameLength {
		return errors.New("problem registering user with small username")
	}
	return nil
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	type status struct {
		Message string `json:"status"`
	}

	ret := status{Message: "ok"}
	data, err := json.Marshal(&ret)
	if err != nil {
		log.Println("could not encode message: ", err)
		fmt.Fprintf(w, "Error")
		return
	}

	fmt.Fprintf(w, "%s", data)
}

var mutex sync.Mutex

// POST
func handleRegister(w http.ResponseWriter, r *http.Request) {
	type register struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var regReq register
	error := json.NewDecoder(r.Body).Decode(&regReq)
	if error != nil {
		w.WriteHeader(400)
		log.Println("problem parsing registration request")
		errorResponse := errorResponse{Error: "problem parsing registration request"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}
		w.Write(errRespJSON)

		return
	}

	if len(regReq.Password) < minPasswordLength {
		w.WriteHeader(400)
		log.Println("problem registering user with small password")
		errorResponse := errorResponse{Error: "passwords must be larger than 8 characters"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}
		w.Write(errRespJSON)

		return
	}

	if len(regReq.Username) < minUsernameLength {
		w.WriteHeader(400)
		log.Println("problem registering user with small username")
		errorResponse := errorResponse{Error: "problem registering user with small username"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}
		w.Write(errRespJSON)

		return
	}

	registerError := registerUser(regReq.Username, regReq.Password, &mutex)
	if registerError != nil {
		w.WriteHeader(400)
		errorResponse := errorResponse{Error: registerError.Error()}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}
		w.Write(errRespJSON)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginReq loginRequest
	error := json.NewDecoder(r.Body).Decode(&loginReq)
	if error != nil {
		respondWithError(w, error)
		return
	}

	if err := validateUsername(loginReq.Username); err != nil {
		respondWithError(w, err)
		return
	}

	if err := validatePassword(loginReq.Password); err != nil {
		respondWithError(w, err)
		return
	}

	token, err := login(loginReq.Username, loginReq.Password)
	if err != nil {
		respondWithError(w, err)
		return
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}

	tokenJSON, err := json.Marshal(&tokenResponse{token})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(tokenJSON)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		err := errors.New("only post supported")
		respondWithError(w, err)
		return
	}

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if auth[0] != "Token" {
		respondWithError(w, errors.New("token authorization only supported"))
		return
	}
	token := auth[1]

	username, ok := srvState.session[token]
	if !ok {
		respondWithError(w, errors.New("please login"))
		return
	}

	err := upload(r.URL.Path, username, r.Body)
	if err != nil {
		// 500 should be issued here instead
		respondWithError(w, err)
	}
}

func main() {
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/upload/", handleUpload)

	port := fmt.Sprintf(":%s", cmdPort)
	log.Fatal(http.ListenAndServe(port, nil))
}
