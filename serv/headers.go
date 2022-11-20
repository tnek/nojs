package serv

import "net/http"

func commonHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "sameorigin")
	w.Header().Set("X-XSS-Protection", "0")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Vary", "Accept-Encoding")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'none';")
}
