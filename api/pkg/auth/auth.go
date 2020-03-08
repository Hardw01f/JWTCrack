package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"

	"../../pkg/jwtmanage"
)

type Users struct {
	Id         int    `json:id`
	Name       string `json:name`
	Password   string `json:password`
	Email      string `json:email`
	Secretword string `json:secretword`
}

type Status struct {
	Uid    string `json:uid`
	Status string `json:Status`
	Token  string `json:Status`
}

type OutputDetails struct {
	Uid        string `json:uid`
	Name       string `json:name`
	Email      string `json:email`
	Secretword string `json:secretword`
}

func GormInit() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:gopher@(jwt_mysql)/jwtapp?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		err = xerrors.Errorf(": %w", err)
		return nil, err
	}
	return db, nil
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		userName := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		db, err := GormInit()
		if err != nil {
			fmt.Println(err)
		}

		user := Users{}
		user.Name = userName
		user.Email = email
		user.Password = password

		if err := db.Create(&user).Error; err != nil {
			fmt.Println(err)
		}

		if err := db.Where("name=? AND password=?", user.Name, user.Password).Find(&user).Error; err != nil {
			fmt.Println(err)
		}

		signinStatus := Status{}

		signinStatus.Uid = strconv.Itoa(user.Id)
		signinStatus.Status = "true"
		signinStatus.Token, err = jwtmanage.GetJwt(user.Name, signinStatus.Uid)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(signinStatus)
		bytes, _ := json.Marshal(signinStatus)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
		return

	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "POST" {
		email := r.Form["email"]
		password := r.Form["password"]
		fmt.Println(email)
		fmt.Println(password)

		db, err := GormInit()
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			fmt.Printf("%+v\n", err)
		}
		defer db.Close()

		signinStatus := Status{}

		user := Users{}
		if err := db.Where("email = ? AND password = ?", email, password).Find(&user).Error; err != nil {
			fmt.Println(err)
			signinStatus.Uid = ""
			signinStatus.Status = "false"
			fmt.Println(signinStatus)
			bytes, _ := json.Marshal(signinStatus)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(bytes))
			return
		}

		signinStatus.Uid = strconv.Itoa(user.Id)
		signinStatus.Status = "true"
		signinStatus.Token, err = jwtmanage.GetJwt(user.Name, signinStatus.Uid)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(signinStatus)
		bytes, _ := json.Marshal(signinStatus)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
		return
	} else if r.Method == "GET" {
		http.NotFound(w, nil)
		return
	} else {
		http.NotFound(w, nil)
		return
	}
	http.NotFound(w, nil)
	return
}

func Secret(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tokenValues, err := jwtmanage.VerifyJWT(r)
		if err != nil {
			//http.NotFound(w, nil)
			http.Redirect(w, r, "http://localhost:3030/signin", 301)
			return
		}

		if tokenValues.Status == true {
			fmt.Println("verify success")
			fmt.Println(tokenValues.User)

			db, err := GormInit()
			if err != nil {
				fmt.Println(err)
			}
			defer db.Close()

			userDetails := Users{}
			if err := db.Where("name=? AND id=?", tokenValues.User, tokenValues.Uid).Find(&userDetails).Error; err != nil {
				fmt.Println(err)
			}

			fmt.Println(userDetails.Secretword)

			DetailsValue := OutputDetails{}

			DetailsValue.Email = userDetails.Email
			DetailsValue.Name = userDetails.Name
			DetailsValue.Uid = strconv.Itoa(userDetails.Id)
			DetailsValue.Secretword = userDetails.Secretword

			bytes, _ := json.Marshal(DetailsValue)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(bytes))
			return
		} else {
			//http.NotFound(w, nil)
			http.Redirect(w, r, "http://localhost:3030/signin", 301)
			return
		}

	} else {

	}

}
