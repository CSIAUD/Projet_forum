package session

import (
	"html/template"
	"net/http"
	structs "Forum/static/go/structs"
	cookies "Forum/static/go/cookies"
)

var tpl *template.Template
var dbUsers = map[string]structs.User{} //ID utilisateur, struct User
var dbSession = map[string]string{}// ID session, ID User

func init() {
	//Must: aide qui entour un appel de fonction, il retourne un panic en cas d'erreur, utilisé pour des initialisation de variable
	//ParseGlob: créé nouveau template. Equivalent à ParseFiles mais il aura la liste de fichiers correspondant au motif en plus
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func CheckLogIn(w http.ResponseWriter, r *http.Request) {
	user := structs.User{}
	coookie := structs.Cookie{}

	var id int
	var Un string
	var maill string
	var mdpp string
	var av string


	//get the cookie
	cookie, err := r.Cookie("logged-in")
	cookies.SetCookie(cookie, err, w, r)

	//si l'utilisateur existe déjà, on récup ses infos
	//la valeur du cookie est l'UUID, un prend donc la valeur de la session ID
	if un, ok := dbSession[coookie.Value]; ok {
		user = dbUsers[un]//ok il existe, user devien User ID
	}

	if r.Method == http.MethodPost {
		un := r.FormValue(Un)
		mail := r.FormValue(maill)
		mdp := r.FormValue(mdpp)
		avatar := r.FormValue(av)
		user = structs.User{id, un, mail, mdp, avatar}
		dbSession[coookie.Value] = un
		dbUsers[un] = user
	}

	//voir pour le fichier dans la ligne ci-dessous. Possible qu'on doit mettre autre chose
	tpl.ExecuteTemplate(w, "index.html", user)

}

func SignIn() {

}

func LogOut() {
	
}

// func setSession(userName string, response http.ResponseWriter) {
//      value := map[string]string{
//          "name": userName,
//     }
//      if encoded, err := cookieHandler.Encode("logged-in", value); 
// 			err == nil {
//          cookie := &http.Cookie{
//              Name:  "session",
//              Value: encoded,
//              Path:  "/",
//          }
//          http.SetCookie(response, cookie)
//      }
//  }

// func getUserName(request *http.Request) (userName string) {
//     if cookie, err := request.Cookie("logged-in"); err == nil {
//          cookieValue := make(map[string]string)
//     }   
//      return userName
// }
 
// func clearSession(response http.ResponseWriter) {
//     cookie := &http.Cookie { 
//         Name:   "session",
//         Value:  "",
//         Path:   "/",
//         MaxAge: -1,
//     }
//     http.SetCookie(response, cookie)
// }

// func LogoutHandler(w http.ResponseWriter, r *http.Request) {
// 	//read from session
// 	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

// 	// remove the username
// 	session.Values["username"] = ""
// 	err := session.Save(r, w)

// 	if err != nil {
// 			log.Println(err)
// 	}
// 	w.Write([]byte("Logged out!"))
// }
