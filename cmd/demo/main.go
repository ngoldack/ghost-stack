package main

import (
	"html/template"
	"net/http"
	"os"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/ngoldack/ghost-stack/cmd/demo/db"
	"github.com/ngoldack/ghost-stack/cmd/demo/handler"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	r := chi.NewRouter()
	r.Use(handler.Cors)

	db, err := db.NewDB()
	if err != nil {
		panic(err)
	}

	authHandler := handler.NewAuthHandler(db)
	dashboardHandler := handler.NewDashboardHandler()

	templates, err := template.ParseGlob("demo/template/*.tmpl")
	if err != nil {
		panic(err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := templates.Lookup("index.tmpl")
		if err := tmpl.Execute(w, nil); err != nil {
			panic(err)
		}
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl := templates.Lookup("login.tmpl")
		if err := tmpl.Execute(w, nil); err != nil {
			panic(err)
		}
	})

	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		tmpl := templates.Lookup("register.tmpl")
		if err := tmpl.Execute(w, nil); err != nil {
			panic(err)
		}
	})

	r.Post("/auth/login", authHandler.Login())
	r.Post("/auth/register", authHandler.Register())

	r.Route("/dashboard", func(r chi.Router) {
		r.Use(authHandler.Authenticated)

		r.Get("/", dashboardHandler.Dashboard())
	})

	for _, route := range r.Routes() {
		slog.Debug("Route", "route", route.Pattern)
	}

	slog.Info("Starting server on port http://localhost:8080")
	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}
