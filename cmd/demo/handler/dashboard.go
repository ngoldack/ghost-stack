package handler

import "net/http"

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (dh *DashboardHandler) Dashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		<h1>Dashboard</h1>
		<p>This is the dashboard</p>
	`))
	}
}
