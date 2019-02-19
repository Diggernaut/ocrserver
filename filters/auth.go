package filters

import (
	"net/http"
)

// AuthFilter ...
type AuthFilter struct {
	Apikey string
	Next   http.Handler
}

// ServeHTTP ...
func (f *AuthFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == f.Apikey {
		f.Next.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

// SetNext ...
func (f *AuthFilter) SetNext(next http.Handler) {
	f.Next = next
}
