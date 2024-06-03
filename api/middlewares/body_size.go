package middlewares

import "net/http"

func LimitBodySize(max_size int64) NextFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, max_size<<20)
			next.ServeHTTP(w, r)
		})
	}
}
