package cors

import (
	"net/http"
)

// func Middleware(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Response) {
// 		w.Header().Add("Access-Control-Allow-Origin", "*")
// 		w.Header().Add("Content-Type", "application/json")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET,PUT,DELETE,OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type,Content-lenght , Accept-Encoding, X-CSRF-Token, Authorization")
// 		handler.ServeHTTP(w, r)
// 	})
// }

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET,PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler.ServeHTTP(w, r)
	})
}
