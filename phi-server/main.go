// Actual server that handles registration and uploads
// Eventually should have a TLS layer
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	cmdPort = "9876"
)

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

// POST
func handleRegister(w http.ResponseWriter, r *http.Request) {
	type register struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var regReq register
	error := json.NewDecoder(r.Body).Decode(&regReq)
	if error != nil {
		log.Println("problem parsing registration request")
		return
	}

	registerUser(regReq.Username, regReq.Password)
}

func main() {
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/register", handleRegister)

	port := fmt.Sprintf(":%s", cmdPort)
	log.Fatal(http.ListenAndServe(port, nil))
}
