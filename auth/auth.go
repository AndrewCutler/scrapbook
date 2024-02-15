package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func WithSession(next http.HandlerFunc) http.HandlerFunc {
	session := new(Session)

	return Authenticate(next, session)
}

func Authenticate(next http.HandlerFunc, s *Session) http.HandlerFunc {
	fmt.Println("use auth")
	// var s = Session{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		un := r.Form.Get("username")
		passw := r.Form.Get("password")

		fmt.Println("username and password: ", un, " ", passw)

		username, pw, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(pw))
			expectedUsernameHash := sha256.Sum256([]byte("username"))
			expectedPasswordHash := sha256.Sum256([]byte("password"))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			cookie, err := r.Cookie("session_token")
			if err != nil {
				fmt.Println(err)
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fmt.Println(cookie)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				token := uuid.NewString()
				expiration := time.Now().Add(3000 * time.Second)

				s = &Session{
					expiration: expiration,
					username:   username,
				}
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   token,
					Expires: expiration,
				})

				fmt.Println(s)
				return
			}
		}

		http.Error(w, "Unauthorized :(", http.StatusUnauthorized)
	})
}

type Session struct {
	username   string
	expiration time.Time
}

func (s *Session) isExpired() bool {
	return s.expiration.Before(time.Now())
}

func createSession(username string) {
	// sessionToken := uuid.NewString()

}
