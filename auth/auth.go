package auth

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	Username string
	Password string
}

func readConfig() Config {
	file, readErr := os.Open("./config.json")
	if readErr != nil {
		fmt.Println(readErr)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	decodeErr := decoder.Decode(&config)
	if decodeErr != nil {
		fmt.Println("decodeErr: ", decodeErr)
	}

	return config
}

func WithSession(next http.HandlerFunc, session *Session) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}

		var creds Config
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

		config := readConfig()

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

	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("session from authenticate ", s)
	// 	if s.expiration.IsZero() {
	// 		http.Error(w, "Unauthorized :(", http.StatusUnauthorized)
	// 		return
	// 	}

	// if expired, return unauthorized
	// if s.expiration < time.Now() {
	// 	http.Error(w, "Unauthorized: session expired", http.StatusUnauthorized)
	// 	s = new(Session)
	// 	return
	// }

	// next.ServeHTTP(w, r)

	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, "Error parsing form data", http.StatusBadRequest)
	// 	return
	// }

	// un := r.Form.Get("username")
	// passw := r.Form.Get("password")

	// fmt.Println("username and password: ", un, " ", passw)

	// username, pw, ok := r.BasicAuth()
	// if ok {
	// 	usernameHash := sha256.Sum256([]byte(username))
	// 	passwordHash := sha256.Sum256([]byte(pw))
	// 	expectedUsernameHash := sha256.Sum256([]byte("username"))
	// 	expectedPasswordHash := sha256.Sum256([]byte("password"))

	// 	usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
	// 	passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

	// 	cookie, err := r.Cookie("session_token")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		if err == http.ErrNoCookie {
	// 			w.WriteHeader(http.StatusUnauthorized)
	// 			return
	// 		}

	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}

	// 	fmt.Println(cookie)

	// 	if usernameMatch && passwordMatch {
	// 		next.ServeHTTP(w, r)
	// 		token := uuid.NewString()
	// 		expiration := time.Now().Add(3000 * time.Second)

	// 		// s = &Session{
	// 		// 	expiration: expiration,
	// 		// 	username:   username,
	// 		// }
	// 		http.SetCookie(w, &http.Cookie{
	// 			Name:    "session_token",
	// 			Value:   token,
	// 			Expires: expiration,
	// 		})

	// 		fmt.Println(s)
	// 		return
	// 	}
	// }

	// http.Error(w, "Unauthorized :(", http.StatusUnauthorized)
	// })
}

type Session struct {
	username   string
	expiration time.Time
	token      string
}

func (s *Session) isExpired() bool {
	return s.expiration.Before(time.Now())
}

func createSession(username string) {
	// sessionToken := uuid.NewString()

}
