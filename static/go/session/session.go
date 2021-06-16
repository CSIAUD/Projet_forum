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

	"github.com/google/uuid"
	// "strconv"
)

// Fonction qui récupèrel'utilisateur connectégrâce au cookie session
func GetUserByCookie(db bdd.MyDB, w http.ResponseWriter, r *http.Request) (structs.User, error) {

	user := structs.User{}
	cookie, err := r.Cookie("Session")
	if err != nil {
		cookies.SetCookie("Session", "", w, r)
	}
	fmt.Println(cookie)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return user, errors.New("error")
		}
		w.WriteHeader(http.StatusBadRequest)
		return user, errors.New("error")
	}

	cookieV := cookie.Value
	user = (*db.GetUserBySession(cookieV))

	return user, nil
}

//fonction qui permet de vérifier si l'utilisateur existeet si le mdp est valid pour permettre la connexion
func LogIn(mail string, password string, db bdd.MyDB, w http.ResponseWriter, r *http.Request) bool {
	if db.UserExist(mail) {
		id, err := db.CompareMdp(password, mail)
		if err != nil {
			return false
		} else {
			session, err := r.Cookie("Session")
			if err != nil {
				fmt.Print("cookieError :")
				fmt.Println(err)
			} else {
				if session.Value == "0" {
					session.Value = (uuid.New()).String()

				}
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

//fonction qui redirige à la destruction du cookie lors de la déconnexion, retourneun true une fois fait
func LogOut(w http.ResponseWriter, r *http.Request) bool {
	cookies.DestroyCookie(w, r)
	return true
}
