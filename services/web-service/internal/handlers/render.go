package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"web-service/internal/alerts"
	"web-service/internal/session"
	"web-service/internal/views/components"
	"web-service/internal/views/layouts"

	"github.com/a-h/templ"
)

func Base(w http.ResponseWriter, r *http.Request, title string, component templ.Component) error {
	als := alerts.GetAlerts(r.Context())
	base := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		components.Alerts(als...).Render(r.Context(), w)
		component.Render(r.Context(), w)
		return nil
	})

	return layouts.Base(title, base).Render(r.Context(), w)
}

func BaseWithNavbar(w http.ResponseWriter, r *http.Request, title string, component templ.Component) error {
	session, ok := r.Context().Value("session").(*session.UserSession)
	if !ok {
		return fmt.Errorf("failed to get session")
	}
	
	als := alerts.GetAlerts(r.Context())
	base := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		components.Navbar(session.User).Render(r.Context(), w)
		components.Alerts(als...).Render(r.Context(), w)
		component.Render(r.Context(), w)
		return nil
	})

	return layouts.Base(title, base).Render(r.Context(), w)
}
