package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"scavenger/internal/repository"
	"scavenger/internal/reqlog"
	"scavenger/internal/storage"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) FileDownload(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	id := chi.URLParam(r, "id")

	f, rc, err := h.svc.File.Open(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound), errors.Is(err, storage.ErrNotFound):
			http.NotFound(w, r)
		default:
			log.Error("open file", "err", err)
			http.Error(w, "internal", http.StatusInternalServerError)
		}
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", f.ContentType)
	w.Header().Set("Content-Disposition", contentDisposition(f.OriginalName))
	// ServeContent сам проставит Content-Length, Last-Modified,
	// обработает If-Modified-Since и Range.
	http.ServeContent(w, r, f.OriginalName, f.CreatedAt, rc)
}

// contentDisposition отдаёт имя и в ASCII-форме (для старых браузеров),
// и по RFC 5987 (UTF-8) — иначе русские имена приходят «кракозябрами».
func contentDisposition(name string) string {
	ascii := strings.Map(func(r rune) rune {
		if r < 32 || r > 126 || r == '"' || r == '\\' {
			return '_'
		}
		return r
	}, name)
	return fmt.Sprintf(`inline; filename="%s"; filename*=UTF-8''%s`,
		ascii, url.PathEscape(name))
}
