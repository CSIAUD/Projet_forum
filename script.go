package main

import (
	cookie "Forum/static/go/cookies"
	"log"

	bdd "Forum/static/go/bdd"
	session "Forum/static/go/session"
	structs "Forum/static/go/structs"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	template "html/template"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	// "io/iutil"
	// "log"
	// "time
	// "net/ul"
)

var tmplCache map[string]*template.Template
var db bdd.MyDB

func main() {
	var err error
	tmplCache, err = newTemplateCache("./static/html/")
	if err != nil {
		fmt.Printf("Cache Error : %s\n", err)
	}

	// Charger les fichiers du dossier 'static' sur le serveur :
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", http.HandlerFunc(redirectTo404))
	http.Handle("/index", http.HandlerFunc(index)) // Modo a optionnalité en + pour supprimer le post
	http.Handle("/dashBoard", mwIsLogged(http.HandlerFunc(dashBoard)))
	http.Handle("/banList", mwIsLogged(mwIsModo(http.HandlerFunc(banList))))
	http.Handle("/supprpostmodo", mwIsModo(http.HandlerFunc(deletePostModo)))
	http.Handle("/supprpost", mwIsLogged(http.HandlerFunc(deletePost))) //delete posts avec justification pour Modo
	http.Handle("/supprcommentmodo", mwIsModo(http.HandlerFunc(deleteCommentModo)))
	http.Handle("/supprcomment", mwIsLogged(http.HandlerFunc(deleteComment)))
	http.Handle("/deleteUser", mwIsLogged(mwIsAdmin(http.HandlerFunc(deleteUser))))
	http.Handle("/deleteMyUser", mwIsLogged(http.HandlerFunc(deleteMyUser)))
	http.Handle("/comment", http.HandlerFunc(comment))
	http.Handle("/recupmdp", mwIsNotLogged(http.HandlerFunc(recupmdp)))
	http.Handle("/mdp", mwIsNotLogged(http.HandlerFunc(changeMdp)))
	http.Handle("/profil", mwIsLogged(http.HandlerFunc(profil))) // Modo a optionnalité en + pour supprimer profil user
	http.Handle("/profilPost", mwIsLogged(http.HandlerFunc(profilPost)))
	http.Handle("/profilComment", mwIsLogged(http.HandlerFunc(profilComment)))
	http.Handle("/profilInfos", mwIsLogged(http.HandlerFunc(profilInfos)))
	http.Handle("/profilPostUser", mwIsLogged(http.HandlerFunc(profilPostUser))) //affichage des posts de l'utilisateur, voit ceux cacher
	http.Handle("/search", http.HandlerFunc(search))
	http.Handle("/tickets", mwIsLogged(mwIsModo(http.HandlerFunc(tickets))))
	http.Handle("/role", mwIsLogged(mwIsAdmin(http.HandlerFunc(roles))))
	http.Handle("/catModo", mwIsLogged(mwIsAdmin(http.HandlerFunc(catModo))))
	http.Handle("/categories", mwIsLogged(mwIsAdmin(http.HandlerFunc(categories))))
	http.Handle("/signup", mwIsNotLogged(http.HandlerFunc(signup)))
	http.Handle("/login", mwIsNotLogged(http.HandlerFunc(login)))
	http.Handle("/logout", mwIsLogged(http.HandlerFunc(logout)))
	http.Handle("/404", http.HandlerFunc(error404))
	http.Handle("/getStats", mwIsLogged(mwIsModo(http.HandlerFunc(getStat))))
	http.Handle("/share", http.HandlerFunc(share))
	http.Handle("/verif", http.HandlerFunc(verifEmail))
	http.Handle("/vote", mwIsLogged(http.HandlerFunc(vote)))
	http.Handle("/createPost", mwIsLogged(http.HandlerFunc(createPost)))

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

	fmt.Println("Listening server at port 8000.")
	err = http.ListenAndServeTLS(":8000", "https-server.crt", "https-server.key", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//middleware pour vérifier que l'utilisateur est bien connecté avant d'accéder à la page
func mwIsLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GetUserByCookie(db, w, r)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
		} else {
			InitialCookie(w, r)
			next.ServeHTTP(w, r)
			fmt.Println("C0nn3cté")
		}
	})
}

//middleware pour vérifier que l'utilisateur n'est pas connecté avant d'accéder à la page
func mwIsNotLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middlewareLogin")
		next.ServeHTTP(w, r)
		fmt.Println("Executing middlewareLogin again")
	})
}

//middleware pour vérifier que l'utilisateur est modérateur avant d'accéder à la page
func mwIsModo(next http.Handler) http.Handler {
	modoId := 2

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := session.GetUserByCookie(db, w, r)
		fmt.Println(user)
		fmt.Println(user.Role)
		if err != nil {
			fmt.Println(err)
			redirectTo404(w, r)
		} else if user.Role > modoId {
			fmt.Println(err)
			redirectTo404(w, r)
		} else {
			fmt.Println("ok")
			next.ServeHTTP(w, r)
		}
	})
}

//middleware pour vérifier que l'utilisateur est admin avant d'accéder à la page
func mwIsAdmin(next http.Handler) http.Handler {
	adminId := 1
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := session.GetUserByCookie(db, w, r)
		if err != nil {
			redirectTo404(w, r)
		} else if user.Role > adminId {
			redirectTo404(w, r)
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

//fonction pour remplir les pages à l'aide des structs
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

//fonction pour remplis la page index avec les posts de la bdd
func index(w http.ResponseWriter, r *http.Request) {
	// InitialCookie(w, r)
	// cookiePost, err := cookie.GetCookie("Reference", w, r)
	TheCookie := cookie.SetCookie("Session", "", w, r)
	if TheCookie.Value == "0" {
		TheCookie.Value = (uuid.New()).String()
	}
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
	user, _ := session.GetUserByCookie(db, w, r)
	temp.User = user
	err := errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		err = tmplCache["index.page.html"].Execute(w, temp)
		if err != nil {
			panic(err)
		}
	}
}

//fonction de la page create_post qui permet la création d'un post et de l'enregistrer dans la bdd
func createPost(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	fmt.Println(keys)

	if len(keys) == 3 {
		content := keys["content"][0]
		user, _ := strconv.Atoi(keys["userId"][0])
		cat := db.GetCategoryId(keys["categorie"][0])
		db.CreatePost(content, user, cat)
		http.Redirect(w, r, "/index", 302)
	} else {
		data := structs.Categories{}
		data.Categories = (*db.GetAllCategory())
		data.Error = false
		user, _ := session.GetUserByCookie(db, w, r)
		data.User = user
		err := errorGestion(w, r)
		if err != nil {
			http.Redirect(w, r, "/404", 302)
		} else {
			err = tmplCache["create_posts.page.html"].Execute(w, data)
			if err != nil {
				panic(err)
			}
		}
	}
}

//fonction qui redirige l'utilisateur sur une page précise si l'URL n'est pas bonne
func redirectTo404(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/index", 302)
	} else {
		http.Redirect(w, r, "/404", 302)
	}
}

//fonction page error 404
func error404(w http.ResponseWriter, r *http.Request) {
	temp := structs.Err0r{}
	temp.Error = true
	errorGestion(w, r)
	err := tmplCache["error.page.html"].Execute(w, temp)
	if err != nil {
		panic(err)
	}
}

//fonction qui donne le chemin pour détruire le cookie afin de déconnecter l'utilisateur
func logout(w http.ResponseWriter, r *http.Request) {
	session.LogOut(w, r)
}

//fonction de la page des statistiques
func dashBoard(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Error bool
		User  structs.User
		Page  string
		Stats structs.Stats
	}
	data.Stats.Month, _ = db.GetStats(30)
	data.Stats.Seven, _ = db.GetStats(7)
	data.Stats.All, _ = db.GetStats(0)

	// fmt.Println(stats)
	errorGestion(w, r)
	err := tmplCache["categorie_dashboard.page.html"].Execute(w, data)
	if err != nil {
		panic(err)
	}
}

//fonction de la page statisque du forum
func getStat(w http.ResponseWriter, r *http.Request) {
	stats, _ := db.GetStats(0)
	errorGestion(w, r)
	err := tmplCache["stats.page.html"].Execute(w, stats)
	if err != nil {
		panic(err)
	}
}

//fonction de la page qui liste tous les utilisateurs Ban
func banList(w http.ResponseWriter, r *http.Request) {
	ban := structs.BanLists{}

	ban.BanLists = (*db.GetAllBans())
	ban.Error = false
	user, _ := session.GetUserByCookie(db, w, r)
	ban.User = user
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

//fonction qui permet la connection de l'utilisateurvia la page login
func login(w http.ResponseWriter, r *http.Request) {
	temp := structs.Err0r{}
	var user structs.User
	var err error
	user, err = session.GetUserByCookie(db, w, r)
	fmt.Println("User:", user)
	if err != nil {
		temp.User = structs.User{}
	} else {
		temp.User = user
	}

	mail := r.FormValue("mail")
	pass := r.FormValue("password")
	if mail != "" {
		ok := session.LogIn(mail, pass, db, w, r)
		if ok {
			http.Redirect(w, r, "/index", 302)
		}
	}
	errorGestion(w, r)
	err = tmplCache["login.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}

//fonction de la page qui affiche le post avec ses commentaires liés
func comment(w http.ResponseWriter, r *http.Request) {
	var com structs.Commentaires

	var idPost int
	var idUser int
	var content string

	keys := r.URL.Query()

	if len(keys) == 1 {
		idPostS := keys["id"][0]
		idPost, _ = strconv.Atoi(idPostS)
		com.Commentaires = (*db.GetComment(idPost))
		com.Error = false
		com.Post = (*db.GetPost(idPost))
		user, _ := session.GetUserByCookie(db, w, r)
		com.User = user
		com.Commentaires = (*db.GetComment(idPost))
		com.Error = false
		com.Post = (*db.GetPost(idPost))
		user, _ = session.GetUserByCookie(db, w, r)
		com.User = user

		// keys, ok := r.URL.Query()["id"]
		// if !ok || len(keys[0]) < 1 {
		// 	fmt.Println("Url Param 'key' is missing")
		// }

		// clef := keys[0]
		// fmt.Println(clef)
		// idPost, err := strconv.Atoi(clef)
		// if err != nil {
		// 	fmt.Printf("Convert key error : %s", err)
		// }

		err := errorGestion(w, r)
		if err != nil {
			http.Redirect(w, r, "/404", 302)
		} else {
			fmt.Println("T3mPl4t3")
			err = tmplCache["comment.page.html"].Execute(w, com)
			if err != nil {
				panic(err)
			}
		}

	} else if len(keys) > 1 {
		idPostS := keys["id"][0]
		idPost, _ = strconv.Atoi(idPostS)
		users := keys["userId"][0]
		idUser, _ = strconv.Atoi(users)
		content = keys["content"][0]
		db.CreateComment(content, idUser, idPost)
		http.Redirect(w, r, "/comment?id="+strconv.Itoa(idPost), 302)

	}

}

//fonction de la page qui va permettre de partager un post précisavec un url personnalisé
func share(w http.ResponseWriter, r *http.Request) {
	var posts structs.Posts
	var post structs.Post

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
	post.Id = idPost

	//mettre l'id du post dans l'url pour pouvoir ensuite le récupérer et le mettre en argument dans les fonctions pour get
	// -> quand on clique sur le post, on transmet l'id du post liés à ses commentaires pour les afficher
	// com.Commentaires = (*db.GetComment(idPost))
	// com.Error = false
	posts.Posts = append(posts.Posts, post)
	posts.Error = false
	temp, _ := session.GetUserByCookie(db, w, r)
	posts.User = temp

	err = errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		err = tmplCache["share.page.html"].Execute(w, posts)
		if err != nil {
			panic(err)
		}
	}

}

// func post(w http.ResponseWriter, r *http.Request) {

// 	err := errorGestion(w, r)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	err = tmplCache["posts.page.html"].Execute(w, structs.Err0r{})
// 	if err != nil {
// 		panic(err)
// 	}
// }
func deleteComment(w http.ResponseWriter, r *http.Request) {
	var comment structs.Commentaire

	fmt.Println("\n\n\nDELETE COMMENT\n\n\n++++++++++++++++++++++++++++++++++++++++++++")
	keys := r.URL.Query()
	if len(keys) > 0 {
		idComment, err := strconv.Atoi(keys["id"][0])
		comment = (*db.GetCommentById(idComment))
		fmt.Println(comment)
		fmt.Printf("ID commentaire => %d\n\n", idComment)
		if err != nil {
			fmt.Printf("Convert key error : %s", err)
		}
		db.DeleteComment(idComment)
		err = errorGestion(w, r)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/comment?id="+strconv.Itoa(comment.Post.Id), 302)
	} else {
		fmt.Println("nope")
		// return
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++")

}

func deleteUser(w http.ResponseWriter, r *http.Request) { //supprimer définitivement un utilisateur, ses posts et ses commentaires

	keys, ok := r.URL.Query()["userId"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
	}

	clef := keys[0]
	fmt.Println(clef)
	idUser, err := strconv.Atoi(clef)
	if err != nil {
		fmt.Printf("Convert key error : %s", err)
	}
	db.DeleteUser(idUser)

	err = errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		http.Redirect(w, r, "/index", 302)
	}
}

func deleteMyUser(w http.ResponseWriter, r *http.Request) { //supprimer définitivement mon compte, sans mes posts et mes commentaires

	keys, ok := r.URL.Query()["userId"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
	}

	clef := keys[0]
	fmt.Println(clef)
	idUser, err := strconv.Atoi(clef)
	if err != nil {
		fmt.Printf("Convert key error : %s", err)
	}
	db.DeleteMyUser(idUser)

	err = errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		http.Redirect(w, r, "/logout", 302)
	}
}

//fonction de la page où le modérateur et l'administrateur peuvent supprimer le post d'un utilisateur en justifiant le motif
func deletePostModo(w http.ResponseWriter, r *http.Request) {
	var post structs.Post
	var user structs.User
	var data struct {
		Error bool
		User  structs.User
		Page  string
		Id    int
	}

	keys := r.URL.Query()
	if len(keys) < 1 {
		fmt.Println("Url Param 'key' is missing")
	} else {
		fmt.Println(len(keys["id"]))
		fmt.Println(len(keys["comment"]))
		clef := keys["id"][0]
		fmt.Println(clef)
		idPost, err := strconv.Atoi(clef)
		if err != nil {
			fmt.Printf("Convert key error : %s", err)
		}
		post = (*db.GetPost(idPost))
		user = (*db.GetUser(post.User.Id))
		if len(keys) > 1 {
			explanation := keys["content"][0]
			db.DeletePost(idPost)
			SendEmailSuppPost(explanation, user)
			http.Redirect(w, r, "/index", 302)
		}
	}

	data.User, _ = session.GetUserByCookie(db, w, r)
	data.Id, _ = strconv.Atoi(keys["id"][0])
	err := errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["supprpostmodo.page.html"].Execute(w, data)
	if err != nil {
		panic(err)
	}
}

//suppression de son propre post
func deletePost(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	if len(keys) < 1 {
		fmt.Println("Url Param 'key' is missing")
	} else {
		fmt.Println(len(keys["id"]))
		clef := keys["id"][0]
		fmt.Println(clef)
		idPost, err := strconv.Atoi(clef)
		if err != nil {
			fmt.Printf("Convert key error : %s", err)
		}
		db.DeletePost(idPost)
	}

	err := errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/index", 302)

}

//fonction de la page où le modérateur et l'administrateur peuvent supprimer le commentaire d'un utilisateur en justifiant le motif
func deleteCommentModo(w http.ResponseWriter, r *http.Request) {
	var comment structs.Commentaire
	var user structs.User
	var retour struct {
		Error bool
		User  structs.User
		Page  string
		Id    int
	}
	retour.User, _ = session.GetUserByCookie(db, w, r)

	keys := r.URL.Query()
	if len(keys) < 1 {
		fmt.Println("Url Param 'key' is missing")
	} else if len(keys) == 1 {
		clef := keys["comment"][0]
		fmt.Println(clef)
		idComment, err := strconv.Atoi(clef)
		if err != nil {
			fmt.Printf("Convert key error : %s", err)
			fmt.Println(idComment)
		}
		retour.Id = idComment
		err = errorGestion(w, r)
		if err != nil {
			fmt.Println(err)
		}
		err = tmplCache["supprcommentmodo.page.html"].Execute(w, retour)
		if err != nil {
			panic(err)
		}
	} else {
		clef := keys["id"][0]
		fmt.Println(clef)
		idComment, err := strconv.Atoi(clef)
		if err != nil {
			fmt.Printf("Convert key error : %s", err)
		}
		comment = (*db.GetCommentById(idComment))
		user = (*db.GetUser(comment.User.Id))

		explanation := keys["content"][0]
		fmt.Println("motif supp comment :", explanation)
		fmt.Println(user)
		fmt.Println(comment)
		db.DeleteComment(idComment)
		SendEmailSuppComment(explanation, user)
		http.Redirect(w, r, "/comment?id="+strconv.Itoa(comment.Post.Id), 302)
	}

}

//envoit dumail de l'utilisateur sur le motif de la suppression de son post
func SendEmailSuppComment(explanation string, user structs.User) {

	// Sender data.
	fmt.Println("MAIL DELETE POST!!!!!")

	from := "themlcforum@gmail.com"
	password := "YnovAix13100#"

	// Receiver email address.
	to := []string{
		user.Mail,
	}

	fmt.Println("EXPLANATION COMMENTAIRE +++++", explanation)

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	body := fmt.Sprintf("Subject: Suppression de l'un de vos commentaires\nBonjour %s,\n\nVotre commentaire a été supprimé car il ne respectait pas les conditions d'utilisation de MLC Forum.\n\nJustification : %s\n\nAttention, plusieurs infractions peuvent entraîner un bannissement temporaire ou définif du site. \n\n-L'équipe MLC Forum", user.Username, explanation)
	// Message.
	message := []byte(body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("erreur dans l'envoi du mail", err)
		return
	}
	fmt.Println("Email Delete Post Sent Successfully!")
}

//envoit dumail de l'utilisateur sur le motif de la suppression de son post
func SendEmailSuppPost(explanation string, user structs.User) {

	// Sender data.
	fmt.Println("MAIL DELETE POST!!!!!")

	from := "themlcforum@gmail.com"
	password := "YnovAix13100#"

	// Receiver email address.
	to := []string{
		user.Mail,
	}

	fmt.Println("EXPLANATION +++++", explanation)

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	body := fmt.Sprintf("Subject: Suppression de l'un de vos posts\nBonjour %s,\n\nVotre article a été supprimé car il ne respectait pas les conditions d'utilisation de MLC Forum.\n\nJustification : %s\n\nAttention, plusieurs infractions peuvent entraîner un bannissement temporaire ou définif du site. \n\n-L'équipe MLC Forum", user.Username, explanation)
	// Message.
	message := []byte(body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("erreur dans l'envoi du mail", err)
		return
	}
	fmt.Println("Email Delete Post Sent Successfully!")
}

//fonction page profil des utilisateurs
func profil(w http.ResponseWriter, r *http.Request) {
	// var badge structs.BadgeUser
	// var user structs.User

	var profil struct {
		User        structs.User
		UserSession structs.User
		Badges      structs.BadgeUser
		Posts       []structs.Post
		Comments    structs.Commentaires
		Error       error
		Page        string
	}
	// var badges []structs.Badge
	keys := r.URL.Query()
	if len(keys) == 1 {
		clef := keys["id"][0]
		idUser, err := strconv.Atoi(clef)
		if err != nil {
			fmt.Printf("Convert key error : %s", err)
		}

		profil.UserSession = (*db.GetUser(idUser))
		profil.Posts = (*db.GetUserPost(idUser))
		fmt.Println("POSTS:", profil.Posts)
		profil.Comments.Commentaires = (*db.GetUserComment(idUser))

		user, err := session.GetUserByCookie(db, w, r)
		if err != nil {
			profil.User = structs.User{}
		} else {
			profil.User = user
		}
		profil.Badges = (*db.GetBadgeUser(profil.User))
		// badge.User = temp.User
		// badge.Badges = temp.Badges
		// badge.Error = false
	}
	err := errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["profil.page.html"].Execute(w, profil)
	if err != nil {
		panic(err)
	}
}

//fonction de la page qui permet l'enregistrement d'un utilisateur dans ma bdd du forum
func signup(w http.ResponseWriter, r *http.Request) {

	username := ""
	email := ""
	mdp := ""
	confMdp := ""

	keys := r.URL.Query()

	if len(keys) > 0 {
		username = keys["username"][0]
		email = keys["email"][0]
		mdp = keys["password"][0]
		confMdp = keys["conf_password"][0]
		if mdp == confMdp {
			db.CreateUser(username, email, mdp)
			SendEmail(email, username, w, r)
		}
	}
	fmt.Println(username)
	fmt.Println(email)
	fmt.Println(mdp)
	fmt.Println(confMdp)

	errorGestion(w, r)
	err := tmplCache["signup.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}

//fonction de la page des tickets qui répertorient tous les tickets(ouver, fermé, en attente)
func tickets(w http.ResponseWriter, r *http.Request) {
	var ticket structs.Tickets
	attente := 0
	ouvert := 1
	ferme := 2

	ticket.Wait = (*db.GetAllTickt(attente))
	ticket.Open = (*db.GetAllTickt(ouvert))
	ticket.Close = (*db.GetAllTickt(ferme))
	ticket.Error = false
	user, err := session.GetUserByCookie(db, w, r)
	if err != nil {
		ticket.User = structs.User{}
	} else {
		ticket.User = user
	}
	fmt.Println("user: ", ticket.User)
	err = errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["tickets.page.html"].Execute(w, ticket)
	if err != nil {
		panic(err)
	}
}

//fonction page role qui repertorie tous les utilisateurs et applique la modification des rôles de ces derniers
func roles(w http.ResponseWriter, r *http.Request) {
	var role structs.Roles
	var user structs.User
	keys := r.URL.Query()

	if len(keys) == 0 {
		utilisateur := 3
		moderateur := 2
		administrateur := 1

		role.Users = (*db.GetUserByRole(utilisateur))
		role.Modo = (*db.GetUserByRole(moderateur))
		role.Admin = (*db.GetUserByRole(administrateur))
		role.Error = false

	} else {
		input := ""
		roleString := ""
		roleId := 3
		input = keys["search"][0]
		roleString = keys["select"][0]
		roleInt, _ := strconv.Atoi(roleString)
		user = (*db.GetUserByName(input))
		db.UpdateRole(roleInt, user.Id)
		fmt.Println(user)
		fmt.Println("Users:", role.Users)

		if len(keys) > 0 {
		}

		if roleId > 3 || roleId < 0 {
			_ = errorGestion(w, r)
		}
		http.Redirect(w, r, "/role", 302)
	}

	user, err := session.GetUserByCookie(db, w, r)
	if err != nil {
		role.User = structs.User{}
	} else {
		role.User = user
	}
	err = errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["role.page.html"].Execute(w, role)
	if err != nil {
		panic(err)
	}
}

//fonction de la page qui permet la gestion des catégories des modos
func catModo(w http.ResponseWriter, r *http.Request) {
	var role struct {
		Error      bool
		Modo       []structs.UserCat
		User       structs.User
		Page       string
		Categories []structs.Categorie
	}
	var userIdString string
	var userId int
	var catIdString string
	var catId int
	user, err := session.GetUserByCookie(db, w, r)
	if err != nil {
		role.User = structs.User{}
	} else {
		role.User = user
	}

	keys := r.URL.Query()

	if len(keys) == 0 {
		fmt.Println("<===========================================================>")
		fmt.Println("No Redirect to CatModo")
		role.Modo = (*db.GetUserModo())
		role.Error = false
		role.Categories = (*db.GetAllCategory())

		err := errorGestion(w, r)
		if err != nil {
			http.Redirect(w, r, "/404", 302)
		} else {
			err = tmplCache["catModo.page.html"].Execute(w, role)
			if err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Println("je passe ici catMODO!!!")
		userIdString = keys["id"][0]
		catIdString = keys["categorie"][0]
		// userIdString = strings.split(catIdString, "_")[0]
		// catString = strings.split(catIdString, "_")[1]
		// if len(catIdString) > 1 {
		// 	fmt.Println("plusieurs cat!!!!")
		// 	for i := 0; i < len(catIdString); i++ {
		// 		catId, _ = strconv.Atoi(catString[i])
		// 		fmt.Println("les cats:", catId)
		// 		db.DeleteUserCat(userId, catId)
		// 	}

		// } else {
		// fmt.Println("pas plusieurs cat !!!!!!!!")
		userId, _ = strconv.Atoi(userIdString)
		catId, _ = strconv.Atoi(catIdString)
		db.DeleteUserCat(userId, catId)
		// }

		fmt.Println("Redirect to CatModo")
		// http.Redirect(w, r, "/catModo", 302)

	}
}

//fonction de la page pour gérer les catégories
func categories(w http.ResponseWriter, r *http.Request) {
	var choix string
	var idCat int
	fmt.Println(idCat)

	keys := r.URL.Query()
	if len(keys) > 0 {
		fmt.Println(keys)
		if len(keys["delete"]) > 0 {
			fmt.Println("delete")
			choix = keys["delete"][0]
			idCat, _ := strconv.Atoi(choix)
			db.DeleteCategory(idCat)
		} else if len(keys["add"]) > 0 {
			fmt.Println("add")
			choix = keys["add"][0]
			db.CreateCategory(choix)

		} else if len(keys["update"]) > 0 {
			fmt.Println("update")
			choix = keys["update"][0]
			idCat, _ := strconv.Atoi(choix)
			db.UpdateCategory(idCat, keys["name"][0])
		}
		http.Redirect(w, r, "/categories", 302)
	} else {
		data := structs.Categories{}
		data.Categories = (*db.GetAllCategory())
		data.Error = false
		user, _ := session.GetUserByCookie(db, w, r)
		data.User = user
		err := errorGestion(w, r)
		if err != nil {
			http.Redirect(w, r, "/404", 302)
		} else {
			err = tmplCache["categories.page.html"].Execute(w, data)
			if err != nil {
				panic(err)
			}
		}
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

//fonction de gestion des erreurs pour toutes les pages
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

//fonction envoie de mail pour vérification de mail à l'inscription
func SendEmail(mail string, username string, w http.ResponseWriter, r *http.Request) {
	// Sender data.

	from := "themlcforum@gmail.com"
	password := "YnovAix13100#"
	// subject:= "Confirmation de l'email"
	// Receiver email address.
	to := []string{
		mail,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	user := (*db.GetUserByName(username))
	// if err != nil {
	// 	fmt.Println("No user cookie")
	// 	return
	// }
	url := fmt.Sprintf("https://localhost:8000/verif?id=%d", user.Id)
	body := fmt.Sprintf("Subject: Confirmation de l'email\nBonjour %s,\n\nBienvenue sur MLC Forum! Afin de finaliser votre inscription, veuillez cliquer sur le lien ci-dessous pour confirmer votre adresse mail:\n%s\n\nEn vous souhaitant une agréable navigation au sein de notre petit navire!\n\n-L'équipe MLC Forum", username, url)
	// Message.
	message := []byte(body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

//fonction de la page qui confirme la vérification du mail de l'utilisateur
func verifEmail(w http.ResponseWriter, r *http.Request) {
	var user structs.User

	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
	}

	clef := keys[0]
	idUser, err := strconv.Atoi(clef)
	if err != nil {
		fmt.Printf("Convert key error : %s", err)
	}

	user = (*db.GetUser(idUser))

	err = errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		err = tmplCache["verif.page.html"].Execute(w, structs.Err0r{})
		if err != nil {
			panic(err)
		} else {
			db.UserVerified(user.Mail)
			fmt.Println(user)
		}
	}
}

//fonction page qui permet à l'utilisateur d'indiquer son mail pourobtenir un mail de récupération de mdp
func recupmdp(w http.ResponseWriter, r *http.Request) {
	var user structs.User
	email := r.FormValue("mail")

	// keys := r.URL.Query()
	fmt.Println("fonction recup MDP")
	fmt.Println("email user: ", email)
	user = *(db.GetUserByMail(email))

	fmt.Println("USER!!!:", user)

	token := (uuid.New()).String()
	fmt.Println("TOKEN:", token)
	db.CreateUserToken(user.Id, token)

	SendEmailMdp(user, token)

	errorGestion(w, r)
	err := tmplCache["recupmdp.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}

//mail de récupération de mdp
func SendEmailMdp(user structs.User, token string) {
	// Sender data.
	fmt.Println("MAIL RECUP MDP!!!!!")

	from := "themlcforum@gmail.com"
	password := "YnovAix13100#"

	// Receiver email address.
	to := []string{
		user.Mail,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	url := fmt.Sprintf("https://localhost:8000/mdp?token=%s", token)
	body := fmt.Sprintf("Subject: Récupération de compte\nBonjour %s,\n\nUne demande de changement de mot de passe a été faite sur votre compte. Afin de finaliser cette récupération, veuillez cliquer sur le lien ci-dessous :\n%s\n\nSi vous n'êtes pas à l'origine de cette demande, merci d'ignorer cet e-mail.\n\n-L'équipe MLC Forum", user.Username, url)
	// Message.
	message := []byte(body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("erreur dans l'envoi du mail", err)
		return
	}
	fmt.Println("Email Reset Password Sent Successfully!")
}

//fonction de la page qui va permettre à l'utilisateur de rentrer un nouveau mdp si ce dernier l'a oublié
func changeMdp(w http.ResponseWriter, r *http.Request) {
	var user structs.User
	var userToken structs.UserToken
	r.ParseForm()
	pass := r.FormValue("password")
	confPass := r.FormValue("conf_password")

	keys, ok := r.URL.Query()["token"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
	}
	fmt.Println("nouveau mdp:", pass)
	fmt.Println("nouveau mdp conf:", confPass)
	clef := keys[0]
	fmt.Println("KEYS:", clef)

	userToken = (*db.GetUserToken(clef))
	fmt.Println("USER TOKEN:", userToken)

	if (confPass == pass) && (len(pass) > 0) {
		fmt.Println("MDP EQUIVALENT!!!")
		db.UpdateMdpForgotten(userToken.Userid, pass)
		user = (*db.GetUser(userToken.Userid))
		fmt.Println("User: ", user)
	} else {
		fmt.Println("can't do the update password")
	}

	errorGestion(w, r)
	err := tmplCache["mdp.page.html"].Execute(w, structs.Err0r{})
	if err != nil {
		panic(err)
	}
}

//fonction qui rempli la page de profil pour l'onglet post avec les posts des utilisateurs (les non masqués uniquement)
func profilPost(w http.ResponseWriter, r *http.Request) {
	var post structs.Posts

	post.Error = false

	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
	}

	clef := keys[0]
	fmt.Println(clef)
	idUser, err := strconv.Atoi(clef)
	if err != nil {
		fmt.Printf("Convert key error : %s", err)
	}
	post.Posts = (*db.GetUserPost(idUser))
	post.User = (*db.GetUser(idUser))

	err = errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["profilPost.page.html"].Execute(w, post)
	if err != nil {
		panic(err)
	}
}

//fonction qui rempli la page de profil pour l'onglet post avec les posts de utilisateurs (les masqués compris)
func profilPostUser(w http.ResponseWriter, r *http.Request) {
	var post structs.Posts

	post.Error = false

	user, err := session.GetUserByCookie(db, w, r)
	if err != nil {
		post.User = structs.User{}
	} else {
		post.User = user
	}
	post = (*db.GetTheUserPost(user))

	err = errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["profilPostUser.page.html"].Execute(w, post)
	if err != nil {
		panic(err)
	}
}

//onglet profil qui répertorie les infos de l'utilisateur
func profilInfos(w http.ResponseWriter, r *http.Request) {
	var user structs.User

	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
	}

	clef := keys[0]
	fmt.Println(clef)
	idUser, err := strconv.Atoi(clef)
	if err != nil {
		fmt.Printf("Convert key error : %s", err)
	}

	user = (*db.GetUser(idUser))

	err = errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["profilInfos.page.html"].Execute(w, user)
	if err != nil {
		panic(err)
	}
}

//onglet qui répertorie tous les commentaires de l'utilisateur
func profilComment(w http.ResponseWriter, r *http.Request) {
	var comment structs.Commentaires

	comment.Error = false

	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
	}

	clef := keys[0]
	fmt.Println(clef)
	idUser, err := strconv.Atoi(clef)
	if err != nil {
		fmt.Printf("Convert key error : %s", err)
	}
	comment.Commentaires = (*db.GetUserComment(idUser))
	comment.User = (*db.GetUser(idUser))

	err = errorGestion(w, r)
	if err != nil {
		fmt.Println(err)
	}
	err = tmplCache["profilComment.page.html"].Execute(w, comment)
	if err != nil {
		panic(err)
	}
}

//fonction de la page search qui permet de rechercher des posts ainsi que filtrer par catégories
func search(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Categories []structs.Categorie
		Error      bool
		User       structs.User
		Page       string
		Posts      []structs.Post
	}
	data.Categories = (*db.GetAllCategory())
	data.Error = false
	user, _ := session.GetUserByCookie(db, w, r)
	data.User = user
	data.Page = "search"

	keys := r.URL.Query()
	if len(keys) < 1 {
		data.Posts = (*db.GetNbPost(10, 0))
	} else {
		txt := ""
		if len(keys["search"]) > 0 {
			txt = keys["search"][0]
		}
		data.Posts = (*db.Search(txt, keys["categorie"], 10, 0))
	}

	errorGestion(w, r)
	err := tmplCache["research.page.html"].Execute(w, data)
	if err != nil {
		panic(err)
	}
}

//fonction sur le like et dislike des posts
func vote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in VOTE")
	var data struct {
		Val string
	}
	user, _ := session.GetUserByCookie(db, w, r)
	keys := r.URL.Query()
	if len(keys) == 2 {
		post, _ := strconv.Atoi(keys["PostId"][0])
		vote, _ := strconv.Atoi(keys["Vote"][0])
		fmt.Printf("Vote : %d\n", vote)
		db.CreatePostLike(user.Id, post, vote)

		note, total := db.GetLike(post)
		data.Val = note + " like.s pour " + total + " vote.s"
	}

	err := errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		err = tmplCache["data.page.html"].Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func tests(w http.ResponseWriter, r *http.Request) {
	err := errorGestion(w, r)
	if err != nil {
		http.Redirect(w, r, "/404", 302)
	} else {
		err = tmplCache["test.page.html"].Execute(w, structs.Err0r{})
		if err != nil {
			panic(err)
		}
	}
}
