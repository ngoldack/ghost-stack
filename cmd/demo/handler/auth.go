package handler

import (
	"net/http"

	"log/slog"

	"github.com/ngoldack/ghost-stack/cmd/demo/db"
)

type AuthHandler struct {
	db *db.DB
}

func NewAuthHandler(db *db.DB) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

func (ah *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get form values
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Create user
		user, err := ah.db.CreateUser(username, password)
		if err != nil {
			panic(err)
		}

		// Render template
		w.Write([]byte(`
		<h1>Successfully created user</h1>
		<p>ID: ` + user.ID + `</p>
		<p>Username: ` + user.Username + `</p>
	`))
	}
}

func (ah *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get form values
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Search user
		user, err := ah.db.GetUser(username, password)
		if err != nil {
			slog.Error("Failed to get user", "err", err)
			w.Write([]byte(`
				<h1>Failed to login</h1>
				<p>` + err.Error() + `</p>
			`))
			return
		}

		slog.Debug("Found user, creating session...", "user", user)

		// Create session
		session, err := ah.db.GetOrCreateSession(user.ID)
		if err != nil {
			slog.Error("Failed to create session", "err", err)
			w.Write([]byte(`
				<h1>Failed to login</h1>
				<p>` + err.Error() + `</p>
			`))
			return
		}

		// redirect to dashboard
		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   session.ID,
			Expires: session.Expiration,
		})
		http.Redirect(w, r, "/dashboard/", http.StatusFound)
	}
}
