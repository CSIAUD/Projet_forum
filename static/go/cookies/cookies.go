package SetCookie

import (
	// "Forum/static/go/structs"
	// "fmt"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	// "io"
	// structs "Forum/static/go/structs"
)

var html string

//set cookie en plus de donner un uuid
func SetCookie(name string, value string, w http.ResponseWriter, r *http.Request) *http.Cookie {

	cookie, err := r.Cookie(name)

	if err != nil {
		if name == "Session" {
			value = (uuid.New()).String()
		}
		expiration := time.Now().Add(5 * time.Second)
		cookie = &http.Cookie{
			Name:     name,
			Expires:  expiration,
			Value:    value,
			HttpOnly: true,
			Path:     "/",
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
	cookie, err := GetCookieR(w, r)

	if err != nil {
		SetCookie("Reference", "0", w, r)
	} else {
		val, _ := strconv.ParseInt(cookie.Value, 10, 0)
		cookie.Value = strconv.Itoa(int(val) + 1)

	}
	return cookie
}

func GetCookieR(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("Reference")
	return cookie, err

}

func DestroyCookie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("destroy Function")
	// cookie, err := r.Cookie("SessionToken")

	if r.URL.Path == "/logout" {
		cookie := &http.Cookie{
			Name:   "Session",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)

	} else {
		cookie2 := &http.Cookie{
			Name:   "Reference",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie2)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
