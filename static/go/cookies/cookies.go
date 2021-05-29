package SetCookie

import(
	"net/http"
	"github.com/google/uuid"
	"time"
	"fmt"
	// "io"
	// structs "Forum/static/go/structs"
)

var html string

//set cookie en plus de donner un uuid
func SetCookie(cookie *http.Cookie, err error, w http.ResponseWriter, r *http.Request) *http.Cookie {

	if cookie == nil {
		id := uuid.New()
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie = &http.Cookie{
			Name: "logged-in",
			Expires: expiration,
			Value: id.String(),
			HttpOnly: true,
			Path: "/",
		}
		http.SetCookie(w, cookie)
		
	}
	// coookie, _ := r.Cookie("logged-in")
	// fmt.Fprint(w, r, coookie)
	fmt.Println(cookie)
	return cookie

}

//dans la requête POST du serveur,
//si le mot de passe est bon et qu'il peut login
// func LogInCookie(password string, cookie *http.Cookie) *http.Cookie{
// 	//on compare l'input du mot de passe avec celle de la session de l'utilisateur

// 	user := structs.User{}
// 	if password == user.Mdp {
// 		expiration := time.Now().Add(365 * 24 * time.Hour)
// 		cookie = &http.Cookie{
// 			Name: "logged-in",
// 			Value: string(user.Id),
// 			Expires: expiration,
// 			HttpOnly: true,
// 			Path: "/",
// 		}
// 	}
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

// }

//log out et détruit le cookie
func DestroyCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {

	cookie, err := r.Cookie("logged-in") 

	if r.URL.Path == "/logout" {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie = &http.Cookie{
			Name: "logged-in",
			Expires: expiration,
			Value: "0",
			HttpOnly: true,
			Path: "/logout",
			MaxAge:-1,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		if err != nil {
			http.Redirect(w, r, "/set", http.StatusSeeOther) // cookie delete quand on va dans une autre page ?
		}
		cookie.MaxAge = -1 // delete cookie 
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

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
	return cookie
}
