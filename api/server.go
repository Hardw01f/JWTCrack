package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"./pkg/auth"
	"./pkg/posts"
)

func main() {
	var portNum = flag.String("p", "80", "Specify application server listening port")
	flag.Parse()
	fmt.Println("API server listening : " + *portNum)

	//http.HandleFunc("/getjwt", jwtmanage.GetJwt)
	http.HandleFunc("/", auth.Secret)
	http.HandleFunc("/posts", posts.PostLists)
	http.HandleFunc("/signin", auth.SignIn)
	http.HandleFunc("/signup", auth.SignUp)
	err := http.ListenAndServe(":"+*portNum, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
