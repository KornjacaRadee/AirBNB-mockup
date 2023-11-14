package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reservation-service/domain"
)

type KeyProduct struct{}

type AvailabilityPeriodsHandler struct {
	logger *log.Logger
	// NoSQL: injecting accommodation repository
	repo *domain.AvailabilityRepo
}

// NewAccommodationsHandler Injecting the logger makes this code much more testable.
func NewAvailabilityPeriodsHandler(l *log.Logger, r *domain.AvailabilityRepo) *AvailabilityPeriodsHandler {
	return &AvailabilityPeriodsHandler{l, r}
}

func (a *AvailabilityPeriodsHandler) GetAllAvailabilityPeriods(rw http.ResponseWriter, h *http.Request) {
	availabilityPeriods, err := a.repo.GetAll()
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	if availabilityPeriods == nil {
		return
	}

	err = availabilityPeriods.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (a *AvailabilityPeriodsHandler) PostAvailabilityPeriod(rw http.ResponseWriter, h *http.Request) {
	availabilityPeriod := h.Context().Value(KeyProduct{}).(*domain.AvailabilityPeriod)
	err := a.repo.Insert(availabilityPeriod)
	if err != nil {
		http.Error(rw, "Unable to create availabilityPeriods", http.StatusBadRequest)
		a.logger.Fatal(err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (a *AvailabilityPeriodsHandler) PatchAvailabilityPeriod(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	availabilityPeriod := h.Context().Value(KeyProduct{}).(*domain.AvailabilityPeriod)

	a.repo.Update(id, availabilityPeriod)
	rw.WriteHeader(http.StatusOK)
}

func (a *AvailabilityPeriodsHandler) DeleteAvailabilityPeriod(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	a.repo.Delete(id)
	rw.WriteHeader(http.StatusNoContent)
}

func (a *AvailabilityPeriodsHandler) MiddlewareAvailabilityPeriodDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		availabilityPeriod := &domain.AvailabilityPeriod{}
		err := availabilityPeriod.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, availabilityPeriod)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (a *AvailabilityPeriodsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		a.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
