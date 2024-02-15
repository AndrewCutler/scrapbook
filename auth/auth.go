package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func UseBasicAuth(next http.HandlerFunc) http.HandlerFunc {
	var s = session{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, pw, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(pw))
			expectedUsernameHash := sha256.Sum256([]byte("username"))
			expectedPasswordHash := sha256.Sum256([]byte("password"))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				token := uuid.NewString()
				expiration := time.Now().Add(3000 * time.Second)

				s = session{
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

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

type session struct {
	username   string
	expiration time.Time
}

func (s *session) isExpired() bool {
	return s.expiration.Before(time.Now())
}

func createSession(username string) {
	// sessionToken := uuid.NewString()

}
