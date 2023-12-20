package petHandler

import (
	"encoding/json"
	"net/http"
	"pet/service/petServ"
	"strconv"

	"github.com/gorilla/mux"
)

type PetHandler struct {
	service *petServ.PetService
}

func NewPetHandler(service *petServ.PetService) *PetHandler {
	return &PetHandler{service: service}
}

func RegisterPetHandlers(r *mux.Router, service *petServ.PetService) {
	handler := NewPetHandler(service)

	r.HandleFunc("/pets/{id}", mw.TokenAuthMiddleware(handler.GetPetByIDHandler)).Methods("GET")
	r.HandleFunc("/pets", mw.TokenAuthMiddleware(handler.GetPetByStatusHandler)).Methods("GET")
	r.HandleFunc("/pets/{id}", mw.TokenAuthMiddleware(handler.UpdatePetByIDHandler)).Methods("PUT")
	r.HandleFunc("/pets/{id}", mw.TokenAuthMiddleware(handler.DeletePetByIDHandler)).Methods("DELETE")
	r.HandleFunc("/pets/{id}", mw.TokenAuthMiddleware(handler.PostPetByIDHandler)).Methods("POST")
	r.HandleFunc("/pets/{id}/uploadImage", mw.TokenAuthMiddleware(handler.PostImageByIDHandler)).Methods("POST")
	r.HandleFunc("/pets", mw.TokenAuthMiddleware(handler.PostPetHandler)).Methods("POST")
}

func (h *PetHandler) GetPetByIDHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	pet, err := h.service.GetPetByID(ID)
	if err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}
	respondJSON(w, pet, http.StatusOK)
}

func (h *PetHandler) GetPetByStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	pet, err := h.service.GetPetByStatus(status)
	if err != nil {
		http.Error(w, "Pet not found", http.StatusBadRequest)
		return
	}
	respondJSON(w, pet, http.StatusOK)
}

func (h *PetHandler) UpdatePetByIDHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	var updatedPet petServ.Pet
	if err := json.NewDecoder(r.Body).Decode(&updatedPet); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePetByID(ID, updatedPet); err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PetHandler) DeletePetByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	ID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	_, err = h.service.DeletePetByID(ID)
	if err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
}

// PostPetHandler добавляет нового питомца.
// swagger:route POST /pets pets postPet
// Добавляет нового питомца.
// Responses:
//
//	201: petResponse
//	400: errorResponse
func (h *PetHandler) DeletePetByIDHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeletePetByID(ID); err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *PetHandler) PostPetByIDHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse data", http.StatusBadRequest)
		return
	}
	formData := repository.FormData{
		Name:   r.Form.Get("name"),
		Status: r.Form.Get("status"),
	}

	if err := h.service.PostPetByID(ID, formData); err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *PetHandler) PostImageByIDHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	image, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusBadRequest)
		return
	}
	if _, err := h.service.PostImageByID(ID, image); err != nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *PetHandler) PostPetHandler(w http.ResponseWriter, r *http.Request) {
	var newPet petRepo.Pet
	if err := json.NewDecoder(r.Body).Decode(&newPet); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if _, err := h.service.PostPet(newPet); err != nil {
		http.Error(w, "Failed to add pet", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
