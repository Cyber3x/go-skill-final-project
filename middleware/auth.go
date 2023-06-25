package middleware

import (
	"net/http"
)

func BasicAuth(next http.HandlerFunc, users []struct {
	Name     string `ymal:"name"`
	JMBAG    string `ymal:"jmbag"`
	Password string `ymal:"password"`
}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Extract the username and password from the request
		// Authorization header. If no Authentication header is present
		// or the header value is invalid, then the 'ok' return value
		// will be false.
		username, password, ok := req.BasicAuth()
		if ok {
			// Go through all users and mathc username and passowrd
			for _, user := range users {
				if user.Name == username &&
					user.Password == password {

					// when we find a match forward the request
					next.ServeHTTP(w, req)
					return
				}
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
	})
}
