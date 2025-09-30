package auth

import (
	"log"
	"net/http"
	"scavenger/internal/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store *sessions.CookieStore
}

func New(cfg models.AuthConfig) *AuthService {
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	return &AuthService{store: store}
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request, user *models.User, password string) bool {

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("Failed to compare hash and password: %v", err)
		return false
	}

	session, _ := s.store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = user.Username
	session.Values["name"] = user.Name
	session.Values["role"] = user.RoleName
	err := session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session: %v", err)
	}
	return true
}

func (s *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Values["authenticated"] = false
	sessions.Save(r, w)
}

func (s *AuthService) IsAuthenticated(r *http.Request) bool {
	session, _ := s.store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	return auth && ok
}

func (s *AuthService) GetUser(r *http.Request) *models.User {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	name, _ := session.Values["name"].(string)
	role, _ := session.Values["role"].(string)

	return &models.User{
		Username: username,
		Name:     name,
		RoleName: role,
	}
}

func (s *AuthService) GetUserRole(r *http.Request) string {
	session, _ := s.store.Get(r, "session")
	group, _ := session.Values["role"].(string)
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
