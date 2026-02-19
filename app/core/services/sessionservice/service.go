package sessionservice

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"scavenger/core/models"
	"scavenger/core/repositories/sessionrepo"

	"github.com/gorilla/sessions"
)

type SessionService struct {
	cookieStore   *sessions.CookieStore
	sessionRepo   *sessionrepo.SessionRepository
	sessionPrefix string
	httpClent     *http.Client
}

const sessionKey = "session"

func NewSessionService(sessionRepo *sessionrepo.SessionRepository, cookieStore *sessions.CookieStore) *SessionService {
	return &SessionService{
		cookieStore: cookieStore,
		sessionRepo: sessionRepo,
		httpClent: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *SessionService) SetCodeCookie(w http.ResponseWriter,r *http.Request, code *models.RegistrationCode) {
	session, _ := s.cookieStore.Get(r, sessionKey)
	session.Values["reg_code"] = code.Code
	session.Values["reg_name"] = code.Name
	session.Values["reg_email"] = code.Email
	session.Values["reg_role"] = code.Role
	if code.GroupID != nil {
		session.Values["reg_group_id"] = *code.GroupID
	}

	session.Save(r, w)
}

func (s *SessionService) GetCodeCookie(w http.ResponseWriter, r *http.Request) (*models.RegistrationCode, error) {
	session, _ := s.cookieStore.Get(r, sessionKey)
	code := &models.RegistrationCode{}
	var ok bool
	code.Code, ok = session.Values["reg_code"].(string)
	if !ok {
		return nil, fmt.Errorf("failed get reg_code from cookie")
	}
	code.Name = session.Values["reg_name"].(string)
	code.Email = session.Values["reg_email"].(string)
	code.Role = session.Values["reg_role"].(string)
	regGroupID, hasGroup := session.Values["reg_group_id"].(int)
	if hasGroup {
		code.GroupID = &regGroupID
	}

	return code, nil
}

func (s *SessionService) DeleteCodeCookie(w http.ResponseWriter, r *http.Request) {
	session, _ := s.cookieStore.Get(r ,sessionKey)
	delete(session.Values, "reg_code")
	delete(session.Values, "reg_name")
	delete(session.Values, "reg_email")
	delete(session.Values, "reg_role")
	delete(session.Values, "reg_group_id")
	session.Save(r, w)
}

func (s *SessionService) SetSessionCockie(w http.ResponseWriter, r *http.Request, sessionID string, expires_at time.Time) {
	session, err := s.cookieStore.Get(r, sessionKey)
	if err != nil {
		log.Printf("WARN Failed to get session cookie: %v", err)
		return
	}

	session.Values["session_id"] = sessionID
	session.Save(r, w)
}

func (s *SessionService) GetSession(r *http.Request) (*models.UserSession, error) {
	session, err := s.cookieStore.Get(r, sessionKey)
	if err != nil {
		return nil, err
	}

	sessionID, ok := session.Values["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("Failed to get session_id from cookie session")
	}

	userSession, err := s.sessionRepo.GetSession(r.Context(), sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session from repo: %w", err)
	}

	return userSession, nil
}

func (s *SessionService) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session, err := s.cookieStore.Get(r, sessionKey)
	if err != nil {
		log.Println("WARN Failed to get session to delete")
	}

	session.Options.MaxAge = -1
	session.Save(r, w)
}
