package handlers

import (
	"accomm-service/domain"
	"accomm-service/repositories"
	"context"
	"log"
	"net/http"
)

type KeyProduct struct{}

type AccommodationsHandler struct {
	logger *log.Logger
	// NoSQL: injecting accommodation repository
	repo *repositories.AccommodationRepo
}

// Injecting the logger makes this code much more testable.
func NewAccommodationsHandler(l *log.Logger, r *repositories.AccommodationRepo) *AccommodationsHandler {
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
	a.repo.Insert(accommodation)
	rw.WriteHeader(http.StatusCreated)
}

func (p *AccommodationsHandler) MiddlewareAccommodationDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		accommodation := &domain.Accommodation{}
		err := accommodation.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
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
