// Appeler func LogInCookie avant de lire cookie (r.cookie)
// Voir emplacement des appels de fonctions (LogInCookie / GetUserByCookie)

package session

import (
	// "html/template"
	structs "Forum/static/go/structs"
	"net/http"

	// cookies "Forum/static/go/cookies"
	bdd "Forum/static/go/bdd"
	"fmt"
	// "strconv"
)

var db bdd.MyDB

// var cookie *http.Cookie
func GetUserByCookie(w http.ResponseWriter, r *http.Request) {

	var session structs.Session
	user := structs.User{}

	cookie, err := r.Cookie("SessionToken")
	fmt.Println(cookie)
	// cookies.LogInCookie(w, cookie, r, err)
	if cookie.Name == "SessionToken" {

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
			user = structs.User{SessionToken: cookieV}
			// fmt.Println("Session Token user :", user.SessionToken)
			// fmt.Println("Cookie Value Session Token:", cookieV)
		}

		session = structs.Session{cookieV, user.Id}
		fmt.Println("je suis un session token cookie")

	}

	fmt.Println(session)

}

func LogIn(Name string, password string) {

	user := structs.User{}

	if Name == user.Username {
		// bdd.db.compareMdp(password, user.Id)
	} else {
		fmt.Println("Pseudo ou email incorrect")
	}
}

func LogOut() {

}
