package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

func test(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./test.gtpl")
	t.Execute(w, nil)

}

type JSONValues struct {
	User   string
	Uid    string
	Status bool
}

type Ping struct {
	Status int
	Rssult string
}

var secretKey = "secret"

func getJwt(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	token.Claims = jwt.MapClaims{
		"user": "admin",
		"uid":  "1",
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

func verifyJwt(w http.ResponseWriter, r *http.Request) {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		b := []byte(secretKey)
		return b, nil
	})
	if err == nil {
		claims := token.Claims.(jwt.MapClaims)
		values := JSONValues{}
		values.User = claims["user"].(string)
		values.Uid = claims["uid"].(string)
		values.Status = true

		JSONData, _ := json.Marshal(values)
		w.Header().Set("Content-Type", "application/json")
		w.Write(JSONData)
	} else {
		values := JSONValues{}
		values.Status = false

		JSONData, _ := json.Marshal(values)
		w.Header().Set("Content-Type", "application/json")
		w.Write(JSONData)
	}
}

func main() {
	var portNum = flag.String("p", "80", "Specify application server listening port")
	flag.Parse()
	fmt.Println("API server listening : " + *portNum)

	http.HandleFunc("/getjwt", getJwt)
	http.HandleFunc("/secret", verifyJwt)
	err := http.ListenAndServe(":"+*portNum, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
