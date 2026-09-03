package sessions

import (
	"errors"
	"layouts-orders-bot/internal/entity"
	"sync"
)

var (
	ErrNoSession = errors.New("no session for current chat")
)

type session struct {
	currentQuestion int
	currentMessage  int
	answers         map[entity.State]string
}

type Repository struct {
	mu       sync.RWMutex
	sessions map[int]*session
}

func New() *Repository {
	return &Repository{sessions: make(map[int]*session)}
}

func (s *Repository) CurQuestion(chatID int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[chatID]
	if !ok {
		return 0, ErrNoSession
	}
	return session.currentQuestion, nil
}

func (s *Repository) CurMessageID(chatID int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[chatID]
	if !ok {
		return 0, ErrNoSession
	}
	return session.currentMessage, nil
}

func (s *Repository) SetMessageID(chatID int, messageID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[chatID]
	if !ok {
		return ErrNoSession
	}
	session.currentMessage = messageID
	return nil
}

func (s *Repository) NextQuestion(chatID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[chatID]
	if !ok {
		return ErrNoSession
	}
	session.currentQuestion++
	return nil
}

func (s *Repository) SaveAnswer(chatID int, key entity.State, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[chatID]
	if !ok {
		return ErrNoSession
	}

	session.answers[key] = value
	return nil
}

func (s *Repository) GetAnswers(chatID int) (map[entity.State]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[chatID]
	if !ok {
		return nil, ErrNoSession
	}
	return session.answers, nil
}

func (s *Repository) StartSession(chatID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[chatID] = &session{
		currentQuestion: 0,
		answers:         make(map[entity.State]string),
	}
	return nil
}

func (s *Repository) FinishSession(chatID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, chatID)
	return nil
}
