package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Practice5/internal/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func (h *Handler) GetCommonFriendsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	user1, _ := strconv.Atoi(query.Get("user1"))
	user2, _ := strconv.Atoi(query.Get("user2"))

	friends, err := h.Repo.GetCommonFriends(user1, user2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(friends)
}

func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize == 0 {
		pageSize = 10
	}

	filters := map[string]string{
		"id": query.Get("id"),
		"name": query.Get("name"),
		"email": query.Get("email"),
		"gender": query.Get("gender"),
		"birth_date": query.Get("birth_date"),
	}

	orderBy := query.Get("order_by")

	users, err := h.Repo.GetPaginatedUsers(page, pageSize, filters, orderBy)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(users)
}