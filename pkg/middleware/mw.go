package middleware

import (
	appErr "Effective-Mobile-Music-Library/pkg/middleware/app-err"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				appError, ok := rec.(*appErr.Error)
				if !ok {
					log.Error(rec)
					appError = appErr.New(http.StatusInternalServerError, "Internal Server Error")
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(appError.StatusCode)
				json.NewEncoder(w).Encode(appError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
