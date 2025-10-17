package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TeseySTD/GoHospitalApi/auth"
	"github.com/TeseySTD/GoHospitalApi/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Message  string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Only POST method allowed")
		return
	}
	
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if loginReq.Username == "" || loginReq.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}
	
	user, err := auth.AuthenticateUser(loginReq.Username, loginReq.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	
	token, err := auth.GenerateToken(user.Username, user.Role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	
	// Відправка токена
	response := LoginResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
		Message:  "Login successful",
	}
	
	utils.RespondJSON(w, http.StatusOK, response)
}

func UsersListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Only GET method allowed")
		return
	}
	
	type UserInfo struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		Note     string `json:"note"`
	}
	
	users := []UserInfo{
		{
			Username: "admin",
			Role:     auth.RoleAdmin,
			Note:     "Full access - password: admin123",
		},
		{
			Username: "reader",
			Role:     auth.RoleReader,
			Note:     "Read-only access - password: reader123",
		},
		{
			Username: "doctor",
			Role:     auth.RoleAdmin,
			Note:     "Full access - password: doctor123",
		},
		{
			Username: "viewer",
			Role:     auth.RoleReader,
			Note:     "Read-only access - password: viewer123",
		},
	}
	
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Available test users",
		"users":   users,
	})
}