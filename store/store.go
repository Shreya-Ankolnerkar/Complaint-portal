package store

import (
	"complaint-portal/models"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrComplaintNotFound = errors.New("complaint not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrEmailTaken        = errors.New("email already registered")
	ErrBadRequest        = errors.New("bad request")
)

type Store struct {
	mu         sync.RWMutex
	users      map[string]*models.User
	bySecret   map[string]string
	byEmail    map[string]string
	complaints map[string]*models.Complaint
	adminCode  string
}

func New(adminSecret string) *Store {
	return &Store{
		users:      make(map[string]*models.User),
		bySecret:   make(map[string]string),
		byEmail:    make(map[string]string),
		complaints: make(map[string]*models.Complaint),
		adminCode:  adminSecret,
	}
}

func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Store) isAdmin(secretCode string) bool {
	return secretCode == s.adminCode
}

func (s *Store) Register(name, email string) (*models.User, error) {
	if name == "" || email == "" {
		return nil, ErrBadRequest
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byEmail[email]; exists {
		return nil, ErrEmailTaken
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}
	secret, err := generateID()
	if err != nil {
		return nil, err
	}

	u := &models.User{
		ID:         id,
		SecretCode: secret,
		Name:       name,
		Email:      email,
		Complaints: []models.Complaint{},
	}

	s.users[id] = u
	s.bySecret[secret] = id
	s.byEmail[email] = id

	return u, nil
}

func (s *Store) Login(secretCode string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uid, ok := s.bySecret[secretCode]
	if !ok {
		return nil, ErrUserNotFound
	}
	return s.users[uid], nil
}

func (s *Store) SubmitComplaint(secretCode, title, summary string, severity int) (*models.Complaint, error) {
	if title == "" || summary == "" {
		return nil, ErrBadRequest
	}
	if severity < 1 || severity > 5 {
		return nil, errors.New("severity must be between 1 and 5")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	uid, ok := s.bySecret[secretCode]
	if !ok {
		return nil, ErrUnauthorized
	}
	user := s.users[uid]

	cid, err := generateID()
	if err != nil {
		return nil, err
	}

	c := &models.Complaint{
		ID:       cid,
		Title:    title,
		Summary:  summary,
		Severity: severity,
		Status:   "open",
		UserID:   uid,
	}

	s.complaints[cid] = c
	user.Complaints = append(user.Complaints, *c)

	return c, nil
}

func (s *Store) GetAllComplaintsForUser(secretCode string) ([]models.Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uid, ok := s.bySecret[secretCode]
	if !ok {
		return nil, ErrUnauthorized
	}

	var result []models.Complaint
	for _, c := range s.complaints {
		if c.UserID == uid {
			result = append(result, *c)
		}
	}
	if result == nil {
		result = []models.Complaint{}
	}
	return result, nil
}

func (s *Store) GetAllComplaintsForAdmin(secretCode string) ([]models.Complaint, error) {
	if !s.isAdmin(secretCode) {
		return nil, ErrUnauthorized
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Complaint
	for _, c := range s.complaints {
		cp := *c
		if u, ok := s.users[c.UserID]; ok {
			cp.UserName = u.Name
		}
		result = append(result, cp)
	}
	if result == nil {
		result = []models.Complaint{}
	}
	return result, nil
}

func (s *Store) ViewComplaint(secretCode, complaintID string) (*models.Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.complaints[complaintID]
	if !ok {
		return nil, ErrComplaintNotFound
	}

	if s.isAdmin(secretCode) {
		cp := *c
		if u, ok2 := s.users[c.UserID]; ok2 {
			cp.UserName = u.Name
		}
		return &cp, nil
	}

	uid, ok := s.bySecret[secretCode]
	if !ok {
		return nil, ErrUnauthorized
	}
	if c.UserID != uid {
		return nil, ErrUnauthorized
	}

	cp := *c
	return &cp, nil
}

func (s *Store) ResolveComplaint(secretCode, complaintID string) (*models.Complaint, error) {
	if !s.isAdmin(secretCode) {
		return nil, ErrUnauthorized
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.complaints[complaintID]
	if !ok {
		return nil, ErrComplaintNotFound
	}

	c.Status = "resolved"

	if u, ok2 := s.users[c.UserID]; ok2 {
		for i := range u.Complaints {
			if u.Complaints[i].ID == complaintID {
				u.Complaints[i].Status = "resolved"
				break
			}
		}
	}

	cp := *c
	return &cp, nil
}
