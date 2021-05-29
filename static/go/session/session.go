// Vérification des sessions
// Dans fonction autorisations : autorisations selon rôle => séparation
// Appeler func LogInCookie avant de lire cookie (r.cookie)
// Voir emplacement des appels de fonctions (LogInCookie / GetUserByCookie)
// Vérifier fonctionnalités du cookie.Name Session Token

package session

import (
	// "html/template"
	"net/http"
	structs "Forum/static/go/structs"
	// cookies "Forum/static/go/cookies"
	bdd "Forum/static/go/bdd"
	"strconv"
	"fmt"
)

var session structs.Session

var db bdd.MyDB
// var cookie *http.Cookie
func GetUserByCookie(w http.ResponseWriter, r *http.Request) {

	user := structs.User{}

	guest := 0

	cookie, err := r.Cookie("SessionToken")
	// cookies.LogInCookie(w, cookie, r, err)

	if (cookie.Name == "SessionToken") {
		

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//on converti la valeur du uuid du cookie en int
		cookieV, err := strconv.Atoi(cookie.Value)
		if err != nil {
			fmt.Println(err)
		}
	
		//on compare l'uuid de l'utilisateur avec celui du cookie
		if cookieV == user.SessionToken {
			session = structs.Session{cookieV, user.Id}
		}
		fmt.Println("je suis un session token cookie")
	} else {
		user.Role = guest
		fmt.Println("hello")
	}
	
	Authorisation(user.Role)

	// User := db.GetUser(session.User_Id)

	fmt.Println(session)

}

func Authorisation(role int) {
	
	guest := 0
	user := 1
	modo := 2
	admin := 3

	if role == guest {
		
	}else if role == user {
		
	}else if role == modo {
		
	}else if role == admin {
		
	}
}

























// if err != nil {
// 	if err == http.ErrNoCookie {
// 		fmt.Println("y a pas de cookie sessionToken")
// 		// If the cookie is not set, return an unauthorized status
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
// 	// For any other type of error, return a bad request status
// 	w.WriteHeader(http.StatusBadRequest)
// 	return
// }


// var tpl *template.Template
// var dbUsers = map[string]structs.User{} //ID utilisateur, struct User
// var dbSession = map[string]string{}// ID session, ID User

// func init() {
// 	//Must: aide qui entour un appel de fonction, il retourne un panic en cas d'erreur, utilisé pour des initialisation de variable
// 	//ParseGlob: créé nouveau template. Equivalent à ParseFiles mais il aura la liste de fichiers correspondant au motif en plus
// 	tpl = template.Must(template.ParseGlob("templates/*"))
// }

// func CheckLogIn(w http.ResponseWriter, r *http.Request) {
// 	user := structs.User{}
// 	coookie := structs.Cookie{}

// 	var id int
// 	var Un string
// 	var maill string
// 	var av string


// 	//get the cookie
// 	cookie, err := r.Cookie("logged-in")
// 	cookies.SetCookie(cookie, err, w, r)

// 	//si l'utilisateur existe déjà, on récup ses infos
// 	//la valeur du cookie est l'UUID, un prend donc la valeur de la session ID
// 	if un, ok := dbSession[coookie.Value]; ok {
// 		user = dbUsers[un]//ok il existe, user devien User ID
// 	}

// 	if r.Method == http.MethodPost {
// 		un := r.FormValue(Un)
// 		mail := r.FormValue(maill)
// 		avatar := r.FormValue(av)
// 		user = structs.User{id, un, mail, avatar}
// 		dbSession[coookie.Value] = un
// 		dbUsers[un] = user
// 	}

// 	//voir pour le fichier dans la ligne ci-dessous. Possible qu'on doit mettre autre chose
// 	tpl.ExecuteTemplate(w, "index.html", user)

// }

// func SignIn() {

// }

// func LogOut() {
	
// }
