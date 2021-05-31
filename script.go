package main

import (
	bdd "Forum/static/go/bdd"
	cookie "Forum/static/go/cookies"

	"path/filepath"

	"fmt"
	template "html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	// "io/ioutil"
	// "log"
	// "time
	// guuid github.com/google/uuid"
	// "net/url"
	"strings"
)

var cache map[string]*template.Template
var db bdd.MyDB

var page string

func main() {
	page = "tickets"

	var err error
	fs := http.FileServer(http.Dir("./static/"))
	http.HandleFunc("/", index)
	// Charger les fchiers du dossier 'static' sur le serveur :
	cache, err := newTemplateCache("./static/html/")

	if err != nil {
		panic(err)
	}
	fmt.Println("============================")
	fmt.Println(cache)
	fmt.Println("============================")

	fmt.Println("Listening server at port 8000.")
	http.ListenAndServe("localhost:8000", nil)
	db.DB.Close()
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache = map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}
	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable.

		name := filepath.Base(page)
		// Parse the page template file in to a template set.
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout at the
		// moment).
		ts, err = ts.ParseGlob(dir + "layout.html")
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}
		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}

func index(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	var temp string
	for i := len(url) - 1; i >= 0; i-- {
		temp += string(url[i])
	}
	slash := strings.Index(temp, "/")
	page = url[len(url)-slash:]

	if r.URL.Path != "/"+page {
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
	// cookies, err := r.Cookie("SessionToken")
	// cookie.LogInCookie(w, cookies, r, err)
	// session.GetUserByCookie(w, r)

	db.GetNbPost(5, 0)

	fmt.Println("=================")
	fmt.Println(r.FormValue("test"))
	fmt.Println("=================")

	fmt.Println(page)

	err = cache[page+".page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
	// tmpl.Execute(w, nil)
}
