package jwtmanage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"golang.org/x/xerrors"
)

type JSONValues struct {
	User   string
	Uid    string
	Status bool
}

type Ping struct {
	Status int
	Rssult string
}

type VerifyUser struct {
	JSONValues
	jwt.StandardClaims
}

var secretKey = "secret"

func GetJwt(userName, uid string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	token.Claims = jwt.MapClaims{
		"user":   userName,
		"uid":    uid,
		"status": true,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		err = xerrors.Errorf(": %w", err)
		return "", nil
	}

	return tokenString, nil
}

func VerifyJWT(r *http.Request) (JSONValues, error) {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		errStruct := VerifyUser{}
		err := xerrors.New("Not get Authorization Header")
		return errStruct.JSONValues, err
	}
	fmt.Println(tokenHeader)

	slicedToken := strings.Split(tokenHeader, " ")
	tokenString := slicedToken[len(slicedToken)-1]

	verifyUser := VerifyUser{}

	token, err := jwt.ParseWithClaims(tokenString, &verifyUser, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	fmt.Println(token.Claims, err)

	if err != nil {
		err = xerrors.Errorf("Not Permitted : %w", err)
		errStruct := VerifyUser{}
		return errStruct.JSONValues, err
	}

	fmt.Println(verifyUser.JSONValues)

	return verifyUser.JSONValues, nil

}

func VerifyJwt(r *http.Request) string {
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
		res := string(JSONData)
		return res
	} else {
		err = xerrors.Errorf("Not Get Token : %w", err)
		fmt.Printf("%+v\n", err)
		values := JSONValues{}
		values.Status = false

		JSONData, _ := json.Marshal(values)
		res := string(JSONData)
		return res
	}
}
