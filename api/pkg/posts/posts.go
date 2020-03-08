package posts

import (
	"fmt"
	"net/http"

	"../../pkg/auth"
	"../../pkg/jwtmanage"
)

func main() {
	fmt.Println("vim-go")
}

type Posts struct {
	ID      string `json:postid`
	Text    string `json:text`
	User_id string `json:userid`
}

type InsertPost struct {
	Text    string `json:text`
	User_id string `json:userid`
}

func PostLists(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tokenValues, err := jwtmanage.VerifyJWT(r)
		if err != nil {
			http.NotFound(w, nil)
			return
		}

		if tokenValues.Status == true {

			db, err := auth.GormInit()
			if err != nil {
				fmt.Println(err)
			}
			defer db.Close()

			posts := []Posts{}
			if err := db.Find(&posts).Error; err != nil {
				fmt.Println(err)
			}

			fmt.Println(posts)

		}

	} else if r.Method == "POST" {
		tokenValues, err := jwtmanage.VerifyJWT(r)
		if err != nil {
			http.NotFound(w, nil)
			return
		}

		if tokenValues.Status == true {
			r.ParseForm()
			text := r.FormValue("text")
			fmt.Println(text)

			db, err := auth.GormInit()
			if err != nil {
				fmt.Println(err)
				http.NotFound(w, nil)
				return
			}

			insertPosts := Posts{}
			insertPosts.User_id = tokenValues.Uid
			insertPosts.Text = text

			if err := db.Create(&insertPosts).Error; err != nil {
				fmt.Println(err)
			}
		}

	}

}
