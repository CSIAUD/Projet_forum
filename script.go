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
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	// "io/iutil"
	// "log"
	// "time

	// "net/ul"
	"database/sql"
)

var tmplCache map[string]*template.Template
var db bdd.MyDB

func main() {
	var err error
	tmplCache, err = newTemplateCache("./static/html/")
	if err != nil {
		fmt.Printf("Cache Error : %s\n", err)
	}

	// Charger les fichiers du dossier 'static' ur le serveur :
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", redirectTo404)
	http.HandleFunc("/index", index)
	http.Handle("/dashBoard", http.HandlerFunc(dashBoard))
	http.HandleFunc("/banList", banList)
	http.HandleFunc("/comment", comment)
	http.HandleFunc("/recupMdp", recupMdp)
	http.HandleFunc("/post", post)
	http.Handle("/profil", http.HandlerFunc(profil))
	http.HandleFunc("/search", search)
	http.HandleFunc("/tickets", tickets)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/404", error404)

	db.DB, err = sql.Open("sqlite3", "./SQLite/mlcData.db")
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	if err != nil {
		panic(err)
	}
	fmt.Println("===========================")
	fmt.Println(tmplCache)
	fmt.Println("===========================")
	// tests()

	fmt.Println("Listening server at port 8000")
	http.ListenAndServe("localhost:8000", nil)
}

func mwIsLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GetUserByCookie(db, w, r)
		if err != nil {
			return
		}
		InitialCookie(w, r)
		next.ServeHTTP(w, r)
	})
}
func mwIsNotLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middlewareLogin")
		next.ServeHTTP(w, r)
		fmt.Println("Executing middlewareLogin again")
	})
}

func mwIsModo(next http.Handler) http.Handler {
	modoId := 2

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := session.GetUserByCookie(db, w, r)
		fmt.Println(user)
		fmt.Println(user.Role)
		if err != nil {
			fmt.Println(err)
			return
		} else if user.Role > modoId {
			fmt.Println(err)
			return
		}
		fmt.Println("ok")

		next.ServeHTTP(w, r)
	})
}

func mwIsAdmin(next http.Handler) http.Handler {
	adminId := 1
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := session.GetUserByCookie(db, w, r)
		if err != nil {
			return
		} else if user.Role > adminId {
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	// InitialCookie(w, r)
	// cookiePost, err := cookie.GetCookie("Reference", w, r)

	// if err != nil {
	// 	cookie.SetCookie("Reference", "0", w, r)
	// 	cookiePost, err = cookie.GetCookie("Reference", w, r)
	// 	if err == nil {
	// 		fmt.Println("value : ", (cookiePost).Value)
	// 	}
	// } else {
	// 	cookie.IncCookieVal(w, r)
	// 	cookiePost, err = cookie.GetCookie("Reference", w, r)
	// 	if err == nil {
	// 		fmt.Println("value : ", (cookiePost).Value)
	// 	}
	// }

	type testCookie struct {
		Page string
		Nb   int
	}
	var cookietest testCookie
	path := []byte("{Page:'index', Nb:0}")
	json.Unmarshal(path, &cookietest)
	// fmt.Printf("page : %s\nnb : %d\n", cookietest.Page, cookietest.Nb)
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
	user, err := session.GetUserByCookie(db, w, r)
	if err != nil {
		temp.User = structs.User{}
	} else {
		temp.User = user
	}
	err = errorGestion(w, r)
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
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/index", 302)
	} else {
		http.Redirect(w, r, "/404", 302)
	}
}

func error404(w http.ResponseWriter, r *http.Request) {
	temp := structs.Err0r{}
	temp.Error = true
	errorGestion(w, r)
	err := tmplCache["error.page.html"].Execute(w, temp)
	if err != nil {
		panic(err)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session.LogOut(w, r)
}
func dashBoard(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r)
	err := tmplCache["categorie_dashboard.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}
func banList(w http.ResponseWriter, r *http.Request) {
	ban := structs.BanLists{}

	ban.BanLists = (*db.GetAllBans())
	ban.Error = false
	err := errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		err := tmplCache["banList.page.html"].Execute(w, ban)
		if err != nil {
			panic(err)
		}
	}

}
func login(w http.ResponseWriter, r *http.Request) {
	temp := structs.Err0r{}
	user, err := session.GetUserByCookie(db, w, r)
	if err != nil {
		temp.User = structs.User{}
	} else {
		temp.User = user
	}
	mail := r.FormValue("mail")
	pass := r.FormValue("password")
	if mail != "" {
		session.LogIn(mail, pass, db, w, r)
		fmt.Println("sdklfjghj")
	}
	errorGestion(w, r)
	err = tmplCache["login.page.html"].Execute(w, temp)
	if err != nil {
		panic(err)
	}
}
func comment(w http.ResponseWriter, r *http.Request) {
	var com structs.Commentaires
	// post := structs.Post{}
	// like := structs.PostLike{}

	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
	}

	clef := keys[0]
	fmt.Println(clef)
	idPost, err := strconv.Atoi(clef)
	if err != nil {
		fmt.Printf("Convert key error : %s", err)
	}

	//mettre l'id du post dans l'url pour pouvoir ensuite le récupérer et le mettre en argument dans les fonctions pour get
	// -> quand on clique sur le post, on transmet l'id du post liés à ses commentaires pour les afficher
	com.Commentaires = (*db.GetComment(idPost))
	com.Error = false
	// like = (*db.GetPostLike(idPost))
	// post = (*db.GetPost(idPost))
	err = errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		fmt.Println("T3mPl4t3")
		err = tmplCache["comment.page.html"].Execute(w, structs.Err0r{})
		if err != nil {
			panic(err)
		}
	}

}
func recupMdp(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r)
	err := tmplCache["mdp.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}
func post(w http.ResponseWriter, r *http.Request) {

	err := errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["posts.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}
func profil(w http.ResponseWriter, r *http.Request) {
	var badge structs.BadgeUser
	// var badges []structs.Badge
	err := errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	user, _ := session.GetUserByCookie(db, w, r)
	temp := db.GetBadgeUser(user)
	badge.User = temp.User
	badge.Badges = temp.Badges
	badge.Error = false

	err = tmplCache["profil.page.html"].Execute(w, temp)
	if err != nil {
		panic(err)
	}
}
func search(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r)
	err := tmplCache["research.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}
func signup(w http.ResponseWriter, r *http.Request) {
	errorGestion(w, r)
	err := tmplCache["signup.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}
func tickets(w http.ResponseWriter, r *http.Request) {
	var ticket structs.Tickets

	ticket.Tickets = (*db.GetAllTickt())
	ticket.Error = false
	err := errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["tickets.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}

func InitialCookie(w http.ResponseWriter, r *http.Request) {
	cookiePost, err := cookie.GetCookie("Session", w, r)

	if err != nil {
		cookie.SetCookie("Session", "0", w, r)
		cookiePost, err = cookie.GetCookie("Session", w, r)
		if err == nil {
			fmt.Println("value : ", (cookiePost).Value)
		}
	} else {
		if cookiePost.Value == "0" {
			cookie.SetCookie("Session", "0", w, r)
			cookiePost, err = cookie.GetCookie("Session", w, r)
			if err == nil {
				fmt.Println("value : ", (cookiePost).Value)
			}
		} else {
			cookie.IncCookieVal(w, r)
			cookiePost, err = cookie.GetCookie("Session", w, r)
			if err == nil {
				fmt.Println("value : ", (cookiePost).Value)
			}
		}
	}
}

func errorGestion(w http.ResponseWriter, r *http.Request) error {
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
	_, err = r.Cookie("Session")
	if err != nil {
		cookie.SetCookie("Session", "", w, r)
	}
	fmt.Println("404 => ", err)
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
	username := "cyp83"
	mail := "csiaud83@gmail.com"
	mdp := "YnovAix13100#"
	avatar := "lechat.png"
	err := db.CreateUser(username, mail, mdp, avatar)

	if err != nil {
		fmt.Println("Uername / Mail déjà utilisé")
	} else {
		fmt.Println("Bienvenue dans la secte")
	}
	fmt.Println(db.GetUser(8))
}
