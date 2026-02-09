package session

import (
	"encoding/gob"
	"log"
	"net/http"
)

func init() {
	gob.Register(Flash{})
}

type FlashType string

type Flash struct {
	Type    FlashType
	Message string
}

const (
	ErrorFlash   FlashType = "danger"
	SuccessFlash FlashType = "success"
	WarningFlash FlashType = "warning"
	InfoFlash    FlashType = "info"
)

func (s *SessionService) AddFlash(w http.ResponseWriter, r *http.Request, flash Flash) error {
	session, err := s.sm.Get(r, sessionKey)
	if err != nil {
		return err
	}

	session.AddFlash(flash)
	session.Save(r, w)

	return nil
}

func (s *SessionService) FlashSuccess(w http.ResponseWriter, r *http.Request, message string) {
	err := s.AddFlash(w, r, Flash{
		Type:    SuccessFlash,
		Message: message,
	})

	if err != nil {
		log.Printf("WARN Failed to flash success")
	}
}

func (s *SessionService) FlashError(w http.ResponseWriter, r *http.Request, message string) {
	err := s.AddFlash(w, r, Flash{
		Type:    ErrorFlash,
		Message: message,
	})

	if err != nil {
		log.Printf("WARN Failed to flash error: %v", err)
	}
}

func (s *SessionService) GetFlashes(w http.ResponseWriter, r *http.Request) []Flash {
	session, err := s.sm.Get(r, sessionKey)
	if err != nil {
		log.Printf("WARN Failed to get flashes session: %v", err)
		return nil
	}

	flashes := session.Flashes()
	session.Save(r, w)

	var fls []Flash
	for _, flash := range flashes {
		if f, ok := flash.(Flash); ok {
			fls = append(fls, f)
		}
	}

	return fls
}
