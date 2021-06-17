package SetCookie

import (
	// "Forum/static/go/structs"
	// "fmt"
	"time"

	"net/http"
	"strconv"

	"github.com/google/uuid"
	// "io"
)

var html string

//Création des cookies,pour le Session, la valeur est un uuid, autrement Reference a une valeur incrémenté
func SetCookie(name string, value string, w http.ResponseWriter, r *http.Request) *http.Cookie {

	cookie, err := r.Cookie(name)

	if err != nil {
		if name == "Session" {
			value = (uuid.New()).String()
		}
		expiration := time.Now().Add(24 * time.Hour)
		cookie = &http.Cookie{
			Name:     name,
			Value:    value,
			Expires:  expiration,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)

	}
	return cookie
}

//on actualise la valeur du cookie référence en fonction de son avancé dans le forum
func CookieActu(w http.ResponseWriter, r *http.Request, url string) *http.Cookie {

	cookie, err := r.Cookie("Reference")

	if err != nil {
		// CookiePosts(w, r)

	} else if url == cookie.Path {

	} else if url != cookie.Path {
		cookie.Path = url
		cookie.Value = "0"
	}
	return cookie
}

//on incrémente la valeur du cookie Référence
func IncCookieVal(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookie, err := GetCookie("Reference", w, r)

	if err != nil {
		SetCookie("Reference", "0", w, r)
	} else {
		val, _ := strconv.ParseInt(cookie.Value, 10, 0)
		cookie.Value = strconv.Itoa(int(val) + 1)

	}
	return cookie
}

//fonction qui get le cookie d'après le nom
func GetCookie(name string, w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	return cookie, err
}

//fonction qui détruit les cookies
func DestroyCookie(w http.ResponseWriter, r *http.Request) {
	// cookie, err := r.Cookie("SessionToken")

	if r.URL.Path == "/logout" {
		cookie := &http.Cookie{
			Name:     "Session",
			Value:    "0",
			MaxAge:   -1,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)

	} else {
		cookie2 := &http.Cookie{
			Name:   "Reference",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie2)
	}

	http.Redirect(w, r, "index", 302)
}
