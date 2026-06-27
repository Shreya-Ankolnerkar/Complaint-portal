package handlers

import (
	"complaint-portal/models"
	"complaint-portal/store"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	Store *store.Store
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func ok(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, models.Response{Success: true, Data: data})
}

func fail(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.Response{Success: false, Message: msg})
}

func decode(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	return json.NewDecoder(r.Body).Decode(dst)
}

func mapStoreError(err error) (int, string) {
	switch {
	case errors.Is(err, store.ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	case errors.Is(err, store.ErrComplaintNotFound):
		return http.StatusNotFound, "complaint not found"
	case errors.Is(err, store.ErrUnauthorized):
		return http.StatusUnauthorized, "unauthorized"
	case errors.Is(err, store.ErrEmailTaken):
		return http.StatusConflict, "email already registered"
	case errors.Is(err, store.ErrBadRequest):
		return http.StatusBadRequest, "missing required fields"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

// POST /register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.RegisterRequest
	if err := decode(r, &req); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	user, err := h.Store.Register(req.Name, req.Email)
	if err != nil {
		code, msg := mapStoreError(err)
		fail(w, code, msg)
		return
	}

	ok(w, user)
}

// POST /login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.LoginRequest
	if err := decode(r, &req); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.SecretCode == "" {
		fail(w, http.StatusBadRequest, "secret_code is required")
		return
	}

	user, err := h.Store.Login(req.SecretCode)
	if err != nil {
		code, msg := mapStoreError(err)
		fail(w, code, msg)
		return
	}

	ok(w, user)
}

// POST /submitComplaint
func (h *Handler) SubmitComplaint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.SubmitComplaintRequest
	if err := decode(r, &req); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.SecretCode == "" || req.Title == "" || req.Summary == "" {
		fail(w, http.StatusBadRequest, "secret_code, title, and summary are required")
		return
	}

	complaint, err := h.Store.SubmitComplaint(req.SecretCode, req.Title, req.Summary, req.Severity)
	if err != nil {
		code, msg := mapStoreError(err)
		fail(w, code, msg)
		return
	}

	ok(w, complaint)
}

// POST /getAllComplaintsForUser
func (h *Handler) GetAllComplaintsForUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.GetComplaintsForUserRequest
	if err := decode(r, &req); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.SecretCode == "" {
		fail(w, http.StatusBadRequest, "secret_code is required")
		return
	}

	complaints, err := h.Store.GetAllComplaintsForUser(req.SecretCode)
	if err != nil {
		code, msg := mapStoreError(err)
		fail(w, code, msg)
		return
	}

	ok(w, complaints)
}

// POST /getAllComplaintsForAdmin
func (h *Handler) GetAllComplaintsForAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.GetAllComplaintsAdminRequest
	if err := decode(r, &req); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.SecretCode == "" {
		fail(w, http.StatusBadRequest, "secret_code is required")
		return
	}

	complaints, err := h.Store.GetAllComplaintsForAdmin(req.SecretCode)
	if err != nil {
		code, msg := mapStoreError(err)
		fail(w, code, msg)
		return
	}

	ok(w, complaints)
}

// POST /viewComplaint
func (h *Handler) ViewComplaint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.ViewComplaintRequest
	if err := decode(r, &req); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.SecretCode == "" || req.ComplaintID == "" {
		fail(w, http.StatusBadRequest, "secret_code and complaint_id are required")
		return
	}

	complaint, err := h.Store.ViewComplaint(req.SecretCode, req.ComplaintID)
	if err != nil {
		code, msg := mapStoreError(err)
		fail(w, code, msg)
		return
	}

	ok(w, complaint)
}

// POST /resolveComplaint
func (h *Handler) ResolveComplaint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.ResolveComplaintRequest
	if err := decode(r, &req); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.SecretCode == "" || req.ComplaintID == "" {
		fail(w, http.StatusBadRequest, "secret_code and complaint_id are required")
		return
	}

	complaint, err := h.Store.ResolveComplaint(req.SecretCode, req.ComplaintID)
	if err != nil {
		code, msg := mapStoreError(err)
		fail(w, code, msg)
		return
	}

	ok(w, complaint)
}
