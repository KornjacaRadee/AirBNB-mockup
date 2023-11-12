package handlers

import (
	"accomm-service/domain"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type KeyProduct struct{}

type AccommodationsHandler struct {
	logger *log.Logger
	// NoSQL: injecting accommodation repository
	repo *domain.AccommodationRepo
}

// NewAccommodationsHandler Injecting the logger makes this code much more testable.
func NewAccommodationsHandler(l *log.Logger, r *domain.AccommodationRepo) *AccommodationsHandler {
	return &AccommodationsHandler{l, r}
}

func (a *AccommodationsHandler) GetAllAccommodations(rw http.ResponseWriter, h *http.Request) {
	accommodations, err := a.repo.GetAll()
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	if accommodations == nil {
		return
	}

	err = accommodations.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (a *AccommodationsHandler) PostAccommodation(rw http.ResponseWriter, h *http.Request) {
	accommodation := h.Context().Value(KeyProduct{}).(*domain.Accommodation)
	err := a.repo.Insert(accommodation)
	if err != nil {
		http.Error(rw, "Unable to post accommodation", http.StatusBadRequest)
		a.logger.Fatal(err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (a *AccommodationsHandler) PatchAccommodation(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	accommodation := h.Context().Value(KeyProduct{}).(*domain.Accommodation)

	a.repo.Update(id, accommodation)
	rw.WriteHeader(http.StatusOK)
}

func (a *AccommodationsHandler) DeleteAccommodation(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	a.repo.Delete(id)
	rw.WriteHeader(http.StatusNoContent)
}

func (a *AccommodationsHandler) MiddlewareAccommodationDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		accommodation := &domain.Accommodation{}
		err := accommodation.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, accommodation)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (a *AccommodationsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		a.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
