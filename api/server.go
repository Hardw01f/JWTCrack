package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"./pkg/auth"
)

func main() {
	var portNum = flag.String("p", "80", "Specify application server listening port")
	flag.Parse()
	fmt.Println("API server listening : " + *portNum)

	//http.HandleFunc("/getjwt", jwtmanage.GetJwt)
	http.HandleFunc("/secret", auth.Secret)
	http.HandleFunc("/signin", auth.SignIn)
	err := http.ListenAndServe(":"+*portNum, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
