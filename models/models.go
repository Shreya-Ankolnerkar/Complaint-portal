package models

type User struct {
	ID         string      `json:"id"`
	SecretCode string      `json:"secret_code"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Complaints []Complaint `json:"complaints"`
}

type Complaint struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Severity int    `json:"severity"`
	Status   string `json:"status"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name,omitempty"`
}

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	SecretCode string `json:"secret_code"`
}

type SubmitComplaintRequest struct {
	SecretCode string `json:"secret_code"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Severity   int    `json:"severity"`
}

type GetComplaintsForUserRequest struct {
	SecretCode string `json:"secret_code"`
}

type GetAllComplaintsAdminRequest struct {
	SecretCode string `json:"secret_code"`
}

type ViewComplaintRequest struct {
	SecretCode  string `json:"secret_code"`
	ComplaintID string `json:"complaint_id"`
}

type ResolveComplaintRequest struct {
	SecretCode  string `json:"secret_code"`
	ComplaintID string `json:"complaint_id"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
