package middleware

import (
	"net"
	"net/http"
)

func TrustedNet(snet string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if snet == "" {
				next.ServeHTTP(w, r)
				return
			}

			_, ipNet, err := net.ParseCIDR(snet)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			realIP := r.Header.Get("X-Real-IP")

			ip := net.ParseIP(realIP)
			if !ipNet.Contains(ip) {
				http.Error(w, "Untrusted subnet", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
