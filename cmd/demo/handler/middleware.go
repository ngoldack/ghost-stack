package handler

import (
	"log/slog"
	"net/http"
)

func Cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("CORS", "origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		h.ServeHTTP(w, r)
	})
}

func (ah *AuthHandler) Authenticated(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Session From Cookie

		for _, c := range r.Cookies() {
			slog.Debug("Cookie", "name", c.Name, "value", c.Value)
		}

		sc, err := r.Cookie("session")
		if err != nil {
			slog.Error("Failed to get session cookie", "err", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ok, err := ah.db.ValidateSession(sc.Value)
		if err != nil {
			slog.Error("Failed to validate session", "err", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if !ok {
			slog.Error("Session is invalid")
			http.SetCookie(w, &http.Cookie{
				Name:   "session",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}
