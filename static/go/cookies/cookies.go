package SetCookie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	// "io"
	// structs "Forum/static/go/structs"
)

var html string

//set cookie en plus de donner un uuid
func SetCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {

	cookie, err := r.Cookie("Session")

	if err != nil {
		id := uuid.New()
		expiration := time.Now().Add(7 * 24 * time.Hour)
		cookie = &http.Cookie{
			Name:     "Session",
			Expires:  expiration,
			Value:    id.String(),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)

	}
	fmt.Println(cookie)
	return cookie

}

//dans la requête POST du serveur,
//si le mot de passe est bon et qu'il peut login
// func LogInCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
// 	//on compare l'input du mot de passe avec celle de la session de l'utilisateur
// 	cookie, err := r.Cookie("SessionToken")

// 	if err != nil {
// 		expiration := time.Now().Add(3 * 24 * time.Hour)
// 		id := uuid.New()
// 		cookie = &http.Cookie{
// 			Name:     "SessionToken",
// 			Value:    id.String(),
// 			Expires:  expiration,
// 			HttpOnly: true,
// 			// Path: "/",
// 		}
// 		http.SetCookie(w, cookie)
// 	}

// 	return cookie
// }

// 	// if cookie.Value == string(user.Id) {
// 	// 	html= `
// 	// 	<!DOCTYPE html>
// 	// 	<html lang="fr">
// 	// 	<head>
// 	// 		<meta charset="UTF-8">
// 	// 		<title></title>
// 	// 	</head>
// 	// 	<body>
// 	// 	<h1>CONNEXION</h1>
// 	// 	<form method="post" action="/">
// 	// 		<h3>Nom d'utilisateur</h3>
// 	// 		<input type="text" name="pseudo">
// 	// 		<h3>Mot de Passe</h3>
// 	// 		<input type="text" name="mdp">
// 	// 	</body>
// 	// 	</html>`
// 	// }
// 	// io.WriteString(r, html)
// 	return cookie

//log out et détruit le cookie
func DestroyCookie(w http.ResponseWriter, r *http.Request) {

	// cookie, err := r.Cookie("SessionToken")

	if r.URL.Path == "/logout" {
		// expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := &http.Cookie{
			Name: "Session",
			// Expires: expiration,
			// Value: "0",
			// HttpOnly: true,
			// Path: "/logout",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	// } else {
	// 	cookie := &http.Cookie{
	// 		Name:   "Session",
	// 		MaxAge: -1,
	// 	}
	// 	http.SetCookie(w, cookie)
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// }

	// if cookie.Value == "0" {
	// 	html= `
	// 	<!DOCTYPE html>
	// 	<html lang="fr">
	// 	<head>
	// 		<meta charset="UTF-8">
	// 		<title></title>
	// 	</head>
	// 	<body>
	// 	<h1><a href="/logout">DECONNEXION</a></h1>
	// 	</body>
	// 	</html>`
	// }
	// io.WriteString(r, html)
}
