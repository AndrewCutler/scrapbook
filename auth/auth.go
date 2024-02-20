package auth

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	utils "scrapbook/utils"

	"github.com/google/uuid"
)

func WithSession(next http.HandlerFunc, session *Session) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}

		var creds utils.Config
		fmt.Println(string(body))
		deserializeErr := json.Unmarshal(body, &creds)
		if deserializeErr != nil {
			fmt.Println(deserializeErr)
			http.Error(w, "Deserialization error", http.StatusBadRequest)
			return
		}

		hasher := sha1.New()
		hasher.Write([]byte(creds.Username))
		usernameSha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		hasher.Reset()
		hasher.Write([]byte(creds.Password))
		pwSha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		config := utils.ReadConfig()

		if usernameSha == config.Username && pwSha == config.Password {
			fmt.Println("Authenticated.")
			token := uuid.NewString()
			expiration := time.Now().Add(3000 * time.Second)

			session.expiration = expiration
			session.username = creds.Username
			session.token = token
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   token,
				Expires: expiration,
			})

			w.WriteHeader(http.StatusOK)
		} else {
			fmt.Println("Not authenticated.")
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func unauthorized(w http.ResponseWriter, s *Session, msg string) {
	fmt.Println(msg)
	s = new(Session)
	http.Error(w, "Unauthorized :(", http.StatusUnauthorized)
}

func Authenticate(next http.HandlerFunc, s *Session) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				unauthorized(w, s, "No cookie")
				return
			}
			unauthorized(w, s, "Error")
			return
		}

		if s.token != c.Value {
			unauthorized(w, s, "Token mismatch")
			return
		}

		if s.isExpired() {
			unauthorized(w, s, "Expired")
			return
		}

		next.ServeHTTP(w, r)
	})
}

type Session struct {
	username   string
	expiration time.Time
	token      string
}

func (s *Session) isExpired() bool {
	return s.expiration.Before(time.Now())
}
