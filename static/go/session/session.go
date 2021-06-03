// Appeler func LogInCookie avant de lire cookie (r.cookie)
// Voir emplacement des appels de fonctions (LogInCookie / GetUserByCookie)

package session

import (
	// "html/template"
	structs "Forum/static/go/structs"
	"net/http"

	cookies "Forum/static/go/cookies"
	"fmt"

	bdd "Forum/static/go/bdd"
	// "strconv"
)

// var cookie *http.Cookie
func GetUserByCookie(db bdd.MyDB, w http.ResponseWriter, r *http.Request) {

	var session structs.Session
	user := structs.User{}

	cookie, err := r.Cookie("Session")
	fmt.Println(cookie)
	// cookies.LogInCookie(w, cookie, r, err)
	if cookie.Name == "Session" {

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//on converti la valeur du uuid du cookie en int
		cookieV := cookie.Value

		fmt.Println("cookie : ", cookieV)
		fmt.Println("session : ", user.SessionToken)
		//on compare l'uuid de l'utilisateur avec celui du cookie
		if cookieV == user.SessionToken {
			fmt.Println("same")
		} else {
			fmt.Println("not same")
			db.SetSession("sessionToken", 1)
			user = structs.User{SessionToken: cookieV}
			// fmt.Println("Session Token user :", user.SessionToken)
			// fmt.Println("Cookie Value Session Token:", cookieV)
		}

		session = structs.Session{Uuid: cookieV}
		fmt.Println("je suis un session token cookie")

	}

	fmt.Println(session)
}

func LogIn(mail string, password string, db bdd.MyDB, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println(mail)
	if db.UserExist(mail) {
		err, id := db.CompareMdp(password, mail)
		if err != nil {
			session, _ := r.Cookie("Session")
			db.SetSession(session.Value, id)
			fmt.Println("session ok")
			// GetUserByCookie(w, r)
		} else {
			fmt.Println("erreur1")
			return false
		}
	} else {
		fmt.Println("erreur2")
		return false
	}
	fmt.Println("connecté")
	return true
}

func LogOut(w http.ResponseWriter, r *http.Request) bool {
	cookies.DestroyCookie(w, r)
	http.Redirect(w, r, "/", 302)
	return true
}
