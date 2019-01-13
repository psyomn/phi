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
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	cmdPort = "9876"
)

type errorResponse struct {
	Error string `json:"error"`
}

func init() {
	flag.StringVar(&cmdPort, "port", cmdPort, "port to listen at")
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

	if len(regReq.Password) < 8 {
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

	if len(regReq.Username) < 8 {
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
	type login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginReq login
	error := json.NewDecoder(r.Body).Decode(&loginReq)
	if error != nil {
		log.Println("problem parsing login request")
		return
	}

}

func main() {
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)

	port := fmt.Sprintf(":%s", cmdPort)
	log.Fatal(http.ListenAndServe(port, nil))
}
