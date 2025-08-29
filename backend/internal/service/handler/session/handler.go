package session

import "skillspark/internal/storage"

type Handler struct {
	sessionRepository storage.SessionRepository
}

func NewHandler(sessionRepository storage.SessionRepository) *Handler {
	return &Handler{
		sessionRepository,
	}
}
