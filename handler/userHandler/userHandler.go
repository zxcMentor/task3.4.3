package userHandler

import (
	"encoding/json"
	"net/http"
	"pet/mw"
	"pet/service/userService"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service *userService.UserService
}

func NewUserHandler(service *userService.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func RegisterUserHandlers(r *mux.Router, service *userService.UserService) {
	handler := NewUserHandler(service)

	r.HandleFunc("/users", mw.TokenAuthMiddleware(handler.GetHandler)).Methods("GET")
	r.HandleFunc("/users", mw.TokenAuthMiddleware(handler.PutHandler)).Methods("PUT")
	r.HandleFunc("/users", mw.TokenAuthMiddleware(handler.PostNewUserHandler)).Methods("POST")
	r.HandleFunc("/login", mw.TokenAuthMiddleware(handler.GetLoginHandler)).Methods("GET")
	r.HandleFunc("/logout", mw.TokenAuthMiddleware(handler.GetLogoutHandler)).Methods("GET")
	r.HandleFunc("/users", mw.TokenAuthMiddleware(handler.DeleteUserHandler)).Methods("DELETE")
	r.HandleFunc("/users/array", mw.TokenAuthMiddleware(handler.PostNewArrayOfUsersHandler)).Methods("POST")
	r.HandleFunc("/users/list", mw.TokenAuthMiddleware(handler.PostNewListOfUserHandler)).Methods("PUT")
}

func (h *UserHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	user, err := h.service.Get(name)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	respondJSON(w, user, http.StatusOK)
}

func (h *UserHandler) PutHandler(w http.ResponseWriter, r *http.Request) {
	var user userService.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Put(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Остальные обработчики аналогично...

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (h *UserHandler) GetLogoutHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = h.service.GetLogout(credentials.Username)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	err := h.service.DeleteUser(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) PostNewArrayOfUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []userRepo.User
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = h.service.PostNewArrayOfUsers(users)
	if err != nil {
		http.Error(w, "Failed to create new users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) PostNewListOfUserHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON-запрос в список User
	var users []userRepo.User
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = h.service.PostNewListOfUser(users...)
	if err != nil {
		http.Error(w, "Failed to create new users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
