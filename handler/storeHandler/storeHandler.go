package storeHandler

import (
	"encoding/json"
	"net/http"
	"pet/service/storeServ"
	"strconv"

	"github.com/gorilla/mux"
)

type StoreHandler struct {
	service *storeServ.StoreService
}

func NewStoreHandler(service *storeServ.StoreService) *StoreHandler {
	return &StoreHandler{service: service}
}

func RegisterStoreHandlers(r *mux.Router, service *storeServ.StoreService) {
	handler := NewStoreHandler(service)

	r.HandleFunc("/store/inventory", mw.TokenAuthMiddleware(handler.InventoryHandler)).Methods("GET")
	r.HandleFunc("/store/order", mw.TokenAuthMiddleware(handler.OrderHandler)).Methods("POST")
	r.HandleFunc("/store/order/{id}", mw.TokenAuthMiddleware(handler.GetOrderHandler)).Methods("GET")
	r.HandleFunc("/store/order/{id}", mw.TokenAuthMiddleware(handler.DeleteOrderHandler)).Methods("DELETE")
}

func (h *StoreHandler) InventoryHandler(w http.ResponseWriter, r *http.Request) {
	inventory, err := h.service.Inventory()
	if err != nil {
		http.Error(w, "Failed to get inventory", http.StatusInternalServerError)
		return
	}
	respondJSON(w, inventory, http.StatusOK)
}

func (h *StoreHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(ID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	respondJSON(w, order, http.StatusOK)
}

func (h *StoreHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteOrder(ID); err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StoreHandler) OrderHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder storeServ.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if _, err := h.service.Order(newOrder); err != nil {
		http.Error(w, "Failed to place order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
