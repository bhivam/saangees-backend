package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/util"
)

type ItemHandler struct {
	logger    *log.Logger
	itemStore data.ItemStore
}

func NewItemHandler(
	logger *log.Logger,
	itemStore data.ItemStore,
) *ItemHandler {
	return &ItemHandler{logger, itemStore}
}

func (itemHandler *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	if !ok || !user.IsAdmin {
		itemHandler.logger.Println("Unauthorized: only admins can create items")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		itemHandler.logger.Println("Error decoding item request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	item := &data.Item{
		Name: req.Name,
		Date: req.Date,
	}

	createdItem, err := itemHandler.itemStore.CreateItem(item)
	if err != nil {
		itemHandler.logger.Println("Error creating item:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := toItemResponse(createdItem)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (itemHandler *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	if !ok || !user.IsAdmin {
		itemHandler.logger.Println("Unauthorized: only admins can update items")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		itemHandler.logger.Println("Error decoding item request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	item := &data.Item{
		ID:   req.ID,
		Name: req.Name,
		Date: req.Date,
	}

	err := itemHandler.itemStore.UpdateItem(item)
	if err != nil {
		itemHandler.logger.Println("Error updating item:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (itemHandler *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	if !ok || user == data.AnonymousUser {
		itemHandler.logger.Println("Unauthorized: user not authenticated")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(itemID, 10, 64)
	if err != nil {
		itemHandler.logger.Println("Error parsing item ID:", err)
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := itemHandler.itemStore.GetItem(id)
	if err != nil {
		itemHandler.logger.Println("Error retrieving item:", err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	res := toItemResponse(item)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (itemHandler *ItemHandler) GetItemsByWeek(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	if !ok || user == data.AnonymousUser {
		itemHandler.logger.Println("Unauthorized: user not authenticated")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	items, err := itemHandler.itemStore.ComingWeekItems()
	if err != nil {
		itemHandler.logger.Println("Error retrieving items for the coming week:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var res ListItemsResponse
	for _, item := range items {
		res = append(res, *toItemResponse(item))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
