package auth

import (
	"net/http"
	"scavenger/internal/models"

	"github.com/gorilla/sessions"
)

type AuthService struct {
	store *sessions.CookieStore
}

func New(cfg models.AuthConfig) *AuthService {
	return &AuthService{store: sessions.NewCookieStore([]byte(cfg.SessionSecret))}
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request, users []models.User, username, password string) bool {
	for _, user := range users {
		if user.Username == username && user.Password == password {
			session, _ := s.store.Get(r, "session")
			session.Values["authenticated"] = true
			session.Values["username"] = user.Username
			session.Values["role"] = user.Role
			session.Values["Group"] = user.Group
			session.Save(r, w)
			return true
		}
	}
	return false
}

func (s *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	session.Values["authenticated"] = false
	sessions.Save(r, w)
}

func (s *AuthService) IsAuthenticated(r *http.Request) bool {
	session, _ := s.store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	return auth && ok
}

func (s *AuthService) GetUserRole(r *http.Request) string {
	session, _ := s.store.Get(r, "session")
	group, _ := session.Values["group"].(string)
	return group
}

func (s *AuthService) GetUsername(r *http.Request) string {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	return username
}

func (s *AuthService) GetGroup(r *http.Request) string {
	session, _ := s.store.Get(r, "session")
	group, _ := session.Values["group"].(string)
	return group

}
