package main

import (
	bdd "Forum/static/go/bdd"
	cookie "Forum/static/go/cookies"
	"database/sql"
	"fmt"
	textTemplate "html/template"
	session "Forum/static/go/session"
	// structs "Forum/static/go/structs"

	_ "github.com/mattn/go-sqlite3"

	// "io/ioutil"
	// "log"
	"net/http"
	// "time"
	// "golang.org/x/crypto/bcrypt"
	// guuid "github.com/google/uuid"
)

	var db bdd.MyDB

func main() {
	var err error

	// Charger les fichiers du dossier 'static' sur le serveur :
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", index)

	db.DB, err = sql.Open("sqlite3", "./SQLite/mlcData.db")
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening server at port 8000.")
	http.ListenAndServe("localhost:8000", nil)
	db.DB.Close()
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 - Page not found", http.StatusNotFound)
	}

	tmpl, err := textTemplate.ParseFiles("./static/html/index.html")

	if err != nil {
		http.Error(w, "Error 400 - Bad Request!", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "POST":
		if err != nil {
			http.Error(w, "Error 500 - Internal Server Error", http.StatusInternalServerError)
			return
		}

	}

	cookie.SetCookie(w, r)
	cookies, err := r.Cookie("SessionToken")
	cookie.LogInCookie(w, cookies, r, err)
	session.GetUserByCookie(w, r)
	
	db.GetNbPost(5,0)

	tmpl.Execute(w, nil)

	// cookie, err := r.Cookie("session-id")
	// if err != nil {
	// 	id := guuid.New()
	// 	expiration := time.Now().Add(365 * 24 * time.Hour)
	// 	cookie = &http.Cookie{Name: "session-id", Value: id.String(), Expires: expiration, HttpOnly: true, Path:"/"}
	// 	http.SetCookie(w, cookie)
	// }
	// fmt.Println(cookie)
}