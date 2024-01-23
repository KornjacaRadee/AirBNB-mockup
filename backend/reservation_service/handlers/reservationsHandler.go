package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reservation_service/client"
	"reservation_service/domain"
	"time"
)

type KeyProduct struct{}

type ReservationsHandler struct {
	logger              *log.Logger
	repo                *domain.ReservationsRepo
	accommodationClient client.AccommodationClient
	notificationClient  client.NotificationClient
}

func NewReservationsHandler(l *log.Logger, r *domain.ReservationsRepo, ac client.AccommodationClient, nc client.NotificationClient) *ReservationsHandler {
	return &ReservationsHandler{l, r, ac, nc}
}

func (r *ReservationsHandler) GetAvailabilityPeriodsByAccommodation(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	accommId := vars["id"]

	availabilityPeriodsByAccommodation, err := r.repo.GetAvailabilityPeriodsByAccommodation(accommId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
	}

	if availabilityPeriodsByAccommodation == nil {
		return
	}

	err = availabilityPeriodsByAccommodation.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (r *ReservationsHandler) GetReservationsByAvailabilityPeriod(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	availabilityPeriodId := vars["id"]

	reservationsByAvailabilityPeriod, err := r.repo.GetReservationsByAvailabilityPeriod(availabilityPeriodId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
	}

	if reservationsByAvailabilityPeriod == nil {
		return
	}

	err = reservationsByAvailabilityPeriod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (r *ReservationsHandler) GetReservationsByGuestId(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	guestId := vars["id"]

	reservationsByAvailabilityPeriod, err := r.repo.GetReservationsByUserId(guestId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
	}

	if reservationsByAvailabilityPeriod == nil {
		return
	}

	err = reservationsByAvailabilityPeriod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (r *ReservationsHandler) InsertAvailabilityPeriodByAccommodation(rw http.ResponseWriter, h *http.Request) {
	availabilityPeriodsByAccommodation := h.Context().Value(KeyProduct{}).(*domain.AvailabilityPeriodByAccommodation)
	accommodation, err := r.accommodationClient.GetAccommodation(h.Context(), availabilityPeriodsByAccommodation.AccommodationId)
	if err != nil {
		r.logger.Print("Cant get accommodation: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil || accommodation == nil {
		r.logger.Print("Accommodation does not exist")
		http.Error(rw, "Accommodation does not exist", http.StatusBadRequest)
		return
	}

	err = r.repo.InsertAvailabilityPeriodByAccommodation(availabilityPeriodsByAccommodation)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (r *ReservationsHandler) InsertReservationByAvailabilityPeriod(rw http.ResponseWriter, h *http.Request) {
	reservationByAvailabilityPeriod := h.Context().Value(KeyProduct{}).(*domain.ReservationByAvailabilityPeriod)
	err := r.repo.InsertReservationByAvailabilityPeriod(reservationByAvailabilityPeriod)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	accommodation, err := r.accommodationClient.GetAccommodation(h.Context(), reservationByAvailabilityPeriod.AccommodationId)
	if err != nil {
		r.logger.Print("Cant get accommodation: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	notification := client.NotificationData{
		Host: client.User{Id: accommodation.Owner.Id},
		Text: "Your accommodation " + accommodation.Name + " has been reserved (by " + reservationByAvailabilityPeriod.GuestId.Hex() + ")",
		Time: time.Now(),
	}
	// Call the profile service and handle fallback logic
	_, err = r.notificationClient.SendReservationNotification(h.Context(), notification)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		http.Error(rw, "Notification service not available, but reservation created", http.StatusCreated)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (a *ReservationsHandler) MiddlewareAvailabilityPeriodDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		availabilityPeriod := &domain.AvailabilityPeriodByAccommodation{}
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

func (a *ReservationsHandler) MiddlewareReservationDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		reservation := &domain.ReservationByAvailabilityPeriod{}
		err := reservation.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, reservation)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (a *ReservationsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		a.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
