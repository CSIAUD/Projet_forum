package main

import (
	bdd "Forum/static/go/bdd"
	cookie "Forum/static/go/cookies"
	session "Forum/static/go/session"
	template "html/template"

	"fmt"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	// "io/iutil"
	// "log"
	// "time
	// guuid gitub.com/google/uuid"
	// "net/ul"
	"database/sql"
	"strings"
)

var tmplCache map[string]*template.Template
var db bdd.MyDB

var location string

func main() {
	location = "tickets"

	var err error
	// Charger les fichiers du dossier 'static' ur le serveur :
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", loadPage)

	db.DB, err = sql.Open("sqlite3", "./SQLite/mlcData.db")
	if err != nil {
		panic(err)
	}
	tests()
	tmplCache, err = newTemplateCache("./static/html/")

	if err != nil {
		panic(err)
	}
	fmt.Println("===========================")
	fmt.Println(tmplCache)
	fmt.Println("============================")

	fmt.Println("Listening server at port 8000")
	http.ListenAndServe("localhost:8000", nil)
	db.DB.Close()
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cahe.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essetially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}
	// Loop through the pages on-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.pge.tmpl') from the full file path
		// and assign it to the name variable.

		name := filepath.Base(page)
		// Parse the page template file in to a templae set.
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'layout' templates to the
		// templateset (in our case, it's just the 'base' layout at the
		// moment).
		ts, err = ts.ParseGlob(dir + "layout.html")
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'partial' templates to the
		// templateset (in our case, it's just the 'footer' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}
		// Add the template set to the cache, sing the name of the page
		// (like 'home.pge.tmpl') as the key.
		cache[name] = ts
	}
	// Return the map
	return cache, nil
}

func loadPage(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	var temp string
	for i := len(url) - 1; i >= 0; i-- {
		temp += string(url[i])
	}
	position_slash := strings.Index(temp, "/")
	location = url[len(url)-position_slash:]

	fmt.Println("page : " + location)

	if r.URL.Path != "/"+location {
		http.Error(w, "404 - page not found", http.StatusNotFound)
	}

	_, err := template.ParseFiles("./static/html/layout.html")

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
	session.GetUserByCookie(w, r)

	fmt.Println("----------------")
	fmt.Println(location)
	fmt.Println("+++++++++++++++++")

	page := location + ".page.html"

	err = tmplCache[page].Execute(w, nil)
	fmt.Println("****************")
	if err != nil {
		panic(err)
	}
	// tmpl.Execute(w, nil)
}
func tests() {
	username := "cyph"
	mail := "azerty@azertyr.fr"
	mdp := "coucou"
	avatar := "lechat.png"
	sessionToken := "zryzrytu-5z6sj6hg4"
	err := db.CreateUser(username, mail, mdp, avatar, sessionToken)

	if err != nil {
		fmt.Println("Username / Mail déjà utilisé")
	} else {
		fmt.Println("Bienvenue dans la secte")
	}

}
