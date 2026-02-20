package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Practice3/internal/usecase"
	"Practice3/pkg/modules"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/users")
	if path == "" || path == "/" {
		h.handleUsers(w, r)
		return
	}

	idStr := strings.Trim(path, "/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	h.handleUserByID(w, r, id)
}

func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		users, err := h.usecase.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to fetch users"}`))
			return
		}

		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var user modules.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid request body"}`))
			return
		}

		id, err := h.usecase.CreateUser(&user)
if err != nil {
    log.Println("CREATE ERROR:", err)
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]string{
        "error": err.Error(),
    })
    return
}


		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int64{
			"id": id,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) handleUserByID(w http.ResponseWriter, r *http.Request, id int64) {

	switch r.Method {

	case http.MethodGet:
		user, err := h.usecase.GetUserByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"user not found"}`))
			return
		}

		json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		var user modules.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid request body"}`))
			return
		}

		user.ID = id

		err := h.usecase.UpdateUser(&user)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"user not found"}`))
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{
			"updated": true,
		})

	case http.MethodDelete:
		_, err := h.usecase.DeleteUser(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"user not found"}`))
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{
			"deleted": true,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
