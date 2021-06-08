// Appeler func LogInCookie avant de lire cookie (r.cookie)
// Voir emplacement des appels de fonctions (LogInCookie / GetUserByCookie)

package session

import (
	// "html/template"
	bdd "Forum/static/go/bdd"
	cookies "Forum/static/go/cookies"
	structs "Forum/static/go/structs"
	"errors"
	"fmt"
	"net/http"
	// "strconv"
)

// var cookie *http.Cookie
func GetUserByCookie(db bdd.MyDB, w http.ResponseWriter, r *http.Request) (structs.User, error) {

	user := structs.User{}
	cookie, err := r.Cookie("Session")
	fmt.Println(cookie)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return user, errors.New("error")
		}
		w.WriteHeader(http.StatusBadRequest)
		return user, errors.New("error")
	}

	//on converti la valeur du uuid du cookie en int
	cookieV := cookie.Value
	user = (*db.GetUserBySession(cookieV))

	return user, nil
}

func LogIn(mail string, password string, db bdd.MyDB, w http.ResponseWriter, r *http.Request) bool {
	if db.UserExist(mail) {
		id, err := db.CompareMdp(password, mail)
		if err != nil {
			fmt.Println(err)
			fmt.Println("erreur1")
			return false
		} else {
			session, err := r.Cookie("Session")
			if err != nil {
				fmt.Print("cookieError :")
				fmt.Println(err)
			} else {
				fmt.Println(session)
				fmt.Println(id)

				db.SetSession(session.Value, id)
				fmt.Println("session ok")
			}
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
	return true
}
