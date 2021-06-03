package main

import (
	bdd "Forum/static/go/bdd"
	cookie "Forum/static/go/cookies"
	session "Forum/static/go/session"
	structs "Forum/static/go/structs"
	"encoding/json"
	"errors"
	"fmt"
	template "html/template"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	// "io/iutil"
	// "log"
	// "time
	// guuidgitub.com/google/uuid"
	// "net/ul"
	"database/sql"
)

var tmplCache map[string]*template.Template
var db bdd.MyDB

func main() {
	var err error
	tmplCache, _ = newTemplateCache("./static/html/")

	// Charger les fichiers du dossier 'static' ur le serveur :
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", redirectTo404)
	http.HandleFunc("/index", index)
	// http.Handle("/dashBoard", mwIsModo(http.HandlerFunc(dashBoard)))
	http.HandleFunc("/banList", banList)
	http.HandleFunc("/login", login)
	http.HandleFunc("/comment", comment)
	http.HandleFunc("/recupMdp", recupMdp)
	http.HandleFunc("/post", post)
	http.HandleFunc("/profil", profil)
	http.HandleFunc("/search", search)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/tickets", tickets)
	http.HandleFunc("/404", error404)

	db.DB, err = sql.Open("sqlite3", "./SQLite/mlcData.db")
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()
	// tests()

	if err != nil {
		panic(err)
	}
	fmt.Println("===========================")
	fmt.Println(tmplCache)
	fmt.Println("===========================")

	fmt.Println("Listening server at port 8000")
	http.ListenAndServe("localhost:8000", nil)
}

func mwIsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		fmt.Println("Executing middlewareOne again")
	})
}

// func mwIsModo(next http.Handler) http.Handler {
// 	session := "cookieSession"
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if (session.GetUserByCookie(w, r)).Role > 2 {
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func mwIsAdmin(next http.Handler) http.Handler {
// 	session := "cookie"
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if (session.GetUserByCookie(w, r)).Role > 1 {
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cahe.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to geta slice of all filepaths with
	// the extension '.page.tmpl'. This essetially gives us a slice of allthe
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}
	// Loop through the pages on-by-one.
	for _, page := range pages {
		// Extract the file name (lke 'home.pge.tmpl') from the full file path
		// and assign it to the name variable.

		name := filepath.Base(page)
		// Parse the page template ile in to a templae set.
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
		// templateset (in our case, it's just the 'footer' partial at th
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

func index(w http.ResponseWriter, r *http.Request) {
	cookiePost, err := cookie.GetCookieR(w, r)

	// nbPosts := 10
	// {page:'index', nb:0}
	// var tab []structs.Post

	if err != nil {
		cookie.SetCookie("Reference", "0", w, r)
		cookiePost, err = cookie.GetCookieR(w, r)
		if err == nil {
			fmt.Println("value : ", (cookiePost).Value)
		}
	} else {
		cookie.IncCookieVal(w, r)
		cookiePost, err = cookie.GetCookieR(w, r)
		if err == nil {
			fmt.Println("value : ", (cookiePost).Value)
		}
	}

	type testCookie struct {
		Page string
		Nb   int
	}
	var cookietest testCookie
	path := []byte("{Page:'index', Nb:0}")
	json.Unmarshal(path, &cookietest)
	fmt.Printf("page : %s\nnb : %d\n", cookietest.Page, cookietest.Nb)
	// if cookie.page == url {
	// 	cookie.nb++
	// }
	// for i := 0; i < cookie.nb; i++ {
	// 	temp := db.GetNbPost(nbPosts, i)
	// 	for j := 0; j < nbPosts; j++ {
	// 		tab = append(tab, (*temp)[j])
	// 	}
	// }

	var temp structs.Posts
	temp.Posts = (*db.GetNbPost(10, 0))
	temp.Error = false
	err = errorGestion(w, r, "index")
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		err = tmplCache["index.page.html"].Execute(w, temp)
		if err != nil {
			panic(err)
		}
	}
}

func redirectTo404(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/404", 302)
}

func error404(w http.ResponseWriter, r *http.Request) {
	temp := structs.Err0r{}
	temp.Error = true
	errorGestion(w, r, "error")
	err := tmplCache["error.page.html"].Execute(w, temp)
	if err != nil {
		panic(err)
	}
}

func dashBoard(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r, "dashBoard")
	err := tmplCache["categorie_dashboard.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func banList(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r, "banList")
	err := tmplCache["banList.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	mail := r.FormValue("mail")
	pass := r.FormValue("password")
	if mail != "" {
		session.LogIn(mail, pass, db, w, r)
		fmt.Println("sdklfjghj")
	}
	errorGestion(w, r, "login")
	err := tmplCache["login.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func comment(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r, "comment")
	err := tmplCache["comment.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func recupMdp(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r, "mdp")
	err := tmplCache["mdp.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func post(w http.ResponseWriter, r *http.Request) {

	errorGestion(w, r, "posts")
	err := tmplCache["posts.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func profil(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r, "profil")
	err := tmplCache["profil.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func search(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r, "search")
	err := tmplCache["research.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func signup(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r, "signup")
	err := tmplCache["signup.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
func tickets(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r, "tickets")
	err := tmplCache["tickets.page.html"].Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func errorGestion(w http.ResponseWriter, r *http.Request, location string) error {

	// if r.URL.Path != "/"+location {
	// 	http.Error(w, "404 - page not found", http.StatusNotFound)
	// }

	_, err := template.ParseFiles("./static/html/layout.html")

	if err != nil {
		http.Error(w, "Error 400 - Bad Request!", http.StatusBadRequest)
		return errors.New("error")
	}

	switch r.Method {
	case "POST":
		if err != nil {
			http.Error(w, "Error 500 - Internal Server Error", http.StatusInternalServerError)
			return errors.New("error")
		}
	}
	cookie.SetCookie("Session", "", w, r)
	return nil
}

// func loadPage(w http.ResponseWriter, r *http.Request) {

// 	url := r.URL.String()
// 	var temp string
// 	for i := len(url) - 1; i >= 0; i-- {
// 		temp += string(url[i])
// 	}
// 	position_slash := strings.Index(temp, "/")
// 	location = url[len(url)-position_slash:]

// 	fmt.Println("page : " + location)

// 	if r.URL.Path != "/"+location {
// 		http.Error(w, "404 - page not ound", http.StatusNotFound)
// 	}
// 	if r.URL.Path == "/login" {
// 		fmt.Println("login")
// 		// session.LogIn(mail, password, w, r)
// 	}

// 	_, err := template.ParseFiles("./static/html/layout.html")

// 	if err != nil {
// 		http.Error(w, "Error 400 - Bad Request!", http.StatusBadRequest)
// 		return
// 	}

// 	switch r.Method {
// 	case "POST":
// 		if err != nil {
// 			http.Error(w, "Error 500 - Internal Server Error", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	tests()

// 	cookie.SetCookie(w, r)
// 	// session.GetUserByCokie(w, r)

// 	page := location + ".page.html"
// 	fmt.Println("----------------")
// 	fmt.Println(page)
// 	fmt.Println("++++++++++++++++")

// 	err = tmplCache[page].Execute(w, nil)
// 	fmt.Println("****************")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// tmpl.Execute(w, nil)
// }
func tests() {
	// username : "cyp"
	// mail := "csiaud83gmail.com"
	// mdp := "Ynov"
	// avatar := "lehat.png"
	// err := db.CreateUser(uername, mail, mdp, avatar)

	// if err != nil {
	// 	fmt.Println("Uername / Mail déjà utilisé")
	// } else {
	// 	fmt.Pritln("Bienvenue dans la secte")
	// }
	fmt.Println(db.GetUser(8))
}
