package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"web-service/internal/session"
	"web-service/internal/views/components"
	"web-service/internal/views/layouts"

	"github.com/a-h/templ"
)

func Base(w http.ResponseWriter, r *http.Request, title string, component templ.Component) error {
	fls, ok := r.Context().Value("flashes").([]session.Flash)
	if !ok {
		log.Printf("WARN Failed to get flashes from context")
	}

	var flashes []session.Flash
	flashes = append(flashes, fls...)
	fls = []session.Flash{}

	base := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		components.Alerts(flashes...).Render(r.Context(), w)
		io.WriteString(w, "<div class=\"container mt-4\"")
		component.Render(r.Context(), w)
		io.WriteString(w, "</div>")
		return nil
	})

	return layouts.Base(title, base).Render(r.Context(), w)
}

func BaseWithNavbar(w http.ResponseWriter, r *http.Request, title string, component templ.Component) error {
	userSession, ok := r.Context().Value("session").(*session.UserSession)
	if !ok {
		return fmt.Errorf("failed to get session")
	}

	fls, ok := r.Context().Value("flashes").([]session.Flash)
	if !ok {
		log.Printf("WARN Failed to get flashes from context")
	}

	var flashes []session.Flash
	flashes = append(flashes, fls...)
	fls = []session.Flash{}

	base := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		components.Navbar(userSession.User).Render(r.Context(), w)
		components.Alerts(flashes...).Render(r.Context(), w)
		io.WriteString(w, "<div class=\"container mt-4\"")
		component.Render(r.Context(), w)
		io.WriteString(w, "</div>")
		return nil
	})

	return layouts.Base(title, base).Render(r.Context(), w)
}
