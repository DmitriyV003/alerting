package middleware

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/dmitriy/alerting/internal/helpers"
)

func Decrypt(privateKey *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if privateKey == nil {
				next.ServeHTTP(w, r)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			decrypted, err := helpers.Decrypt(privateKey, body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(decrypted))
			next.ServeHTTP(w, r)
		})
	}
}
