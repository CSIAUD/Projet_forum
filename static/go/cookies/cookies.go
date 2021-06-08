package SetCookie

import (
	// "Forum/static/go/structs"
	// "fmt"

	"net/http"
	"strconv"

	"github.com/google/uuid"
	// "io"
)

var html string

//set cookie en plus de donner un uuid
func SetCookie(name string, value string, w http.ResponseWriter, r *http.Request) *http.Cookie {

	cookie, err := r.Cookie(name)

	if err != nil {
		if name == "Session" {
			value = (uuid.New()).String()
		}
		cookie = &http.Cookie{
			Name:     name,
			Value:    value,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)

	}
	return cookie
}

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

func GetCookie(name string, w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	return cookie, err
}

func DestroyCookie(w http.ResponseWriter, r *http.Request) {
	// cookie, err := r.Cookie("SessionToken")

	if r.URL.Path == "/logout" {
		cookie := &http.Cookie{
			Name:     "Session",
			Value:    "0",
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
