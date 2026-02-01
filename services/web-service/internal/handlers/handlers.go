package handlers

import (
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Render error: " + err.Error(), http.StatusInternalServerError)
	}
}
