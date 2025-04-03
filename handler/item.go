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

	item := toDataItem(req)

	createdItem, err := itemHandler.itemStore.CreateItem(&item)
	if err != nil {
		itemHandler.logger.Println("Error creating item:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdItem)
}

func (itemHandler *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	if !ok || !user.IsAdmin {
		itemHandler.logger.Println("Unauthorized: only admins can update items")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var item data.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		itemHandler.logger.Println("Error decoding item request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO Add Validation

	err := itemHandler.itemStore.UpdateItem(&item)
	if err != nil {
		itemHandler.logger.Println("Error updating item:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TODO depending on level of access, return all items or only published items
func (itemHandler *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(itemID, 10, 64)
	if err != nil {
		itemHandler.logger.Println("Error parsing item ID:", err)
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := itemHandler.itemStore.GetItem(uint(id))
	if err != nil {
		itemHandler.logger.Println("Error retrieving item:", err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if !item.Published {
		user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)

		if !ok {
			itemHandler.logger.Println("Failed to get user context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if user.IsAnonymous() {
			itemHandler.logger.Println("Unauthorized users cannot get unpublished item.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else if !user.IsAdmin {
			itemHandler.logger.Println("Must be admin to get unpublished item")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// TODO Sort response by date
func (itemHandler *ItemHandler) GetItemsByWeek(w http.ResponseWriter, r *http.Request) {
	items, err := itemHandler.itemStore.ComingWeekItems()
	if err != nil {
		itemHandler.logger.Println("Error retrieving items for the coming week:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	includeUnpublished := ok && user.IsAdmin

	res := ListItemsResponse{}
	for _, item := range items {
		if item.Published || includeUnpublished {
			res = append(res, *item)
		}
	}

	itemHandler.logger.Print("JSON RESPONSE: ")
	json.NewEncoder(itemHandler.logger.Writer()).Encode(res)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
