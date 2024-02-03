package domain

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"
)

type ReservationsRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment and creates a keyspace
func New(logger *log.Logger) (*ReservationsRepo, error) {
	db := os.Getenv("CASS_DB")

	// Connect to default keyspace
	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Create 'reservations' keyspace
	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, "reservations", 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	// Connect to reservations keyspace
	cluster.Keyspace = "reservations"
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &ReservationsRepo{
		session: session,
		logger:  logger,
	}, nil
}

// Disconnect from database
func (rr *ReservationsRepo) CloseSession() {
	rr.session.Close()
}

// Create tables
func (rr *ReservationsRepo) CreateTables() {
	err := rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(accommodation_id text, host_id text, availability_period_id UUID, start_date TIMESTAMP, end_date TIMESTAMP, price int, is_price_per_guest BOOLEAN,
					PRIMARY KEY ((accommodation_id), availability_period_id)) 
					WITH CLUSTERING ORDER BY (availability_period_id ASC)`,
			"availability_periods_by_accommodation")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(availability_period_id UUID, reservation_id UUID, start_date TIMESTAMP, end_date TIMESTAMP, accommodation_id text, host_id text, guest_id text,  guest_num int, price int,
					PRIMARY KEY ((availability_period_id), reservation_id)) 
					WITH CLUSTERING ORDER BY (reservation_id ASC)`,
			"reservations_by_availability_period")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(guest_id text, availability_period_id UUID, reservation_id UUID, start_date TIMESTAMP, end_date TIMESTAMP, accommodation_id text, host_id text,  guest_num int, price int,
					PRIMARY KEY ((guest_id), reservation_id)) 
					WITH CLUSTERING ORDER BY (reservation_id ASC)`,
			"reservations_by_guest")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(availability_period_id UUID, reservation_id UUID, start_date TIMESTAMP, end_date TIMESTAMP,
					PRIMARY KEY ((availability_period_id), start_date, end_date, reservation_id))`,
			"reservations_by_dates")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}
}

func (rr *ReservationsRepo) GetAvailabilityPeriodsByAccommodation(id string) (AvailabilityPeriodsByAccommodation, error) {
	scanner := rr.session.Query(`SELECT accommodation_id, host_id, availability_period_id, start_date, end_date, price, is_price_per_guest FROM availability_periods_by_accommodation WHERE accommodation_id = ?`,
		id).Iter().Scanner()

	var availabilityPeriods AvailabilityPeriodsByAccommodation
	for scanner.Next() {
		var availabilityPeriod AvailabilityPeriodByAccommodation
		var accommodationIdHex string
		var hostIdHex string

		err := scanner.Scan(&accommodationIdHex, &hostIdHex, &availabilityPeriod.Id, &availabilityPeriod.StartDate, &availabilityPeriod.EndDate, &availabilityPeriod.Price, &availabilityPeriod.IsPricePerGuest)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		accommodationId, err := primitive.ObjectIDFromHex(accommodationIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		availabilityPeriod.AccommodationId = accommodationId

		hostId, err := primitive.ObjectIDFromHex(hostIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		availabilityPeriod.HostId = hostId

		availabilityPeriods = append(availabilityPeriods, &availabilityPeriod)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return availabilityPeriods, nil
}

func (rr *ReservationsRepo) GetAvailabilityPeriodById(id string, accommId string) (*AvailabilityPeriodByAccommodation, error) {
	query := rr.session.Query(`SELECT accommodation_id, host_id, availability_period_id, start_date, end_date, price, is_price_per_guest FROM availability_periods_by_accommodation 
                                                                            WHERE accommodation_id = ? AND availability_period_id = ? LIMIT 1`, accommId, id)

	var availabilityPeriod AvailabilityPeriodByAccommodation
	var accommodationIdHex string
	var hostIdHex string

	err := query.Scan(&accommodationIdHex, &hostIdHex, &availabilityPeriod.Id, &availabilityPeriod.StartDate, &availabilityPeriod.EndDate, &availabilityPeriod.Price, &availabilityPeriod.IsPricePerGuest)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	accommodationId, err := primitive.ObjectIDFromHex(accommodationIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	availabilityPeriod.AccommodationId = accommodationId

	hostId, err := primitive.ObjectIDFromHex(hostIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	availabilityPeriod.HostId = hostId

	return &availabilityPeriod, nil
}

func (rr *ReservationsRepo) InsertAvailabilityPeriodByAccommodation(period *AvailabilityPeriodByAccommodation) error {
	err := rr.session.Query(
		`INSERT INTO availability_periods_by_accommodation (accommodation_id, host_id, availability_period_id, start_date, end_date, price, is_price_per_guest) 
		VALUES (?, ?, UUID(), ?, ?, ?, ?)`,
		period.AccommodationId.Hex(), period.HostId.Hex(), period.StartDate, period.EndDate, period.Price, period.IsPricePerGuest).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
}

func (rr *ReservationsRepo) GetReservationsByAvailabilityPeriod(id string) (ReservationsByAvailabilityPeriod, error) {
	scanner := rr.session.Query(`SELECT availability_period_id, reservation_id, start_date, end_date, accommodation_id, host_id, guest_id, guest_num, price FROM reservations_by_availability_period WHERE availability_period_id = ?`,
		id).Iter().Scanner()

	var reservations ReservationsByAvailabilityPeriod
	for scanner.Next() {
		var reservation ReservationByAvailabilityPeriod
		var guestIdHex string
		var accommIdHex string
		var hostIdHex string

		err := scanner.Scan(&reservation.AvailabilityPeriodId, &reservation.Id, &reservation.StartDate, &reservation.EndDate, &accommIdHex, &hostIdHex, &guestIdHex, &reservation.GuestNum, &reservation.Price)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		accommId, err := primitive.ObjectIDFromHex(accommIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservation.AccommodationId = accommId

		guestId, err := primitive.ObjectIDFromHex(guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservation.GuestId = guestId

		hostId, err := primitive.ObjectIDFromHex(hostIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservation.HostId = hostId

		reservations = append(reservations, &reservation)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservations, nil
}

func (rr *ReservationsRepo) CheckActiveReservationsByAvailabilityPeriod(id gocql.UUID) (bool, error) {
	//check if reservation dates overlap with other reservation dates
	iter := rr.session.Query(`
        SELECT reservation_id FROM reservations_by_dates 
        WHERE availability_period_id = ? AND start_date > ?`,
		id, time.Now()).Iter()

	// Iterate over the result set to check if there are any rows
	for iter.Scan(nil) {
		return false, nil
	}

	if err := iter.Close(); err != nil {
		rr.logger.Println(err)
		return false, err
	}

	// If there are no rows, it means no reservations, so return true
	return true, nil
}

func (rr *ReservationsRepo) GetReservationsByUserId(id string) (ReservationsByAvailabilityPeriod, error) {
	scanner := rr.session.Query(`SELECT availability_period_id, reservation_id, start_date, end_date, accommodation_id, host_id, guest_id, guest_num, price FROM reservations_by_guest WHERE guest_id = ?`,
		id).Iter().Scanner()

	var reservations ReservationsByAvailabilityPeriod
	for scanner.Next() {
		var reservation ReservationByAvailabilityPeriod
		var guestIdHex string
		var accommIdHex string
		var hostIdHex string

		err := scanner.Scan(&reservation.AvailabilityPeriodId, &reservation.Id, &reservation.StartDate, &reservation.EndDate, &accommIdHex, &hostIdHex, &guestIdHex, &reservation.GuestNum, &reservation.Price)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		accommId, err := primitive.ObjectIDFromHex(accommIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservation.AccommodationId = accommId

		guestId, err := primitive.ObjectIDFromHex(guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservation.GuestId = guestId

		hostId, err := primitive.ObjectIDFromHex(hostIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservation.HostId = hostId

		reservations = append(reservations, &reservation)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservations, nil
}

func (rr *ReservationsRepo) GetReservationByIdAndGuestId(id string, guestId string) (*ReservationByAvailabilityPeriod, error) {
	query := rr.session.Query(`SELECT availability_period_id, reservation_id, start_date, end_date, accommodation_id, host_id, guest_id, guest_num, price 
											FROM reservations_by_guest WHERE reservation_id = ? AND guest_id = ?  LIMIT 1`, id, guestId)

	var reservation ReservationByAvailabilityPeriod
	var accommodationIdHex string
	var guestIdHex string
	var hostIdHex string

	err := query.Scan(&reservation.AvailabilityPeriodId, &reservation.Id, &reservation.StartDate, &reservation.EndDate, &accommodationIdHex, &hostIdHex, &guestIdHex,
		&reservation.GuestNum, &reservation.Price)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	accommodationId, err := primitive.ObjectIDFromHex(accommodationIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	resGuestId, err := primitive.ObjectIDFromHex(guestIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	hostId, err := primitive.ObjectIDFromHex(hostIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	reservation.HostId = hostId

	reservation.AccommodationId = accommodationId
	reservation.GuestId = resGuestId

	return &reservation, nil
}

func (rr *ReservationsRepo) InsertReservationByAvailabilityPeriod(reservation *ReservationByAvailabilityPeriod) error {
	checkReservationDates, err := rr.checkReservationDates(reservation)
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	if !checkReservationDates {
		err = errors.New("reservation dates not available")
		return err
	}

	//check if reservation dates are within period dates
	availabilityPeriod, err := rr.GetAvailabilityPeriodById(reservation.AvailabilityPeriodId.String(), reservation.AccommodationId.Hex())
	if err != nil {
		rr.logger.Println(err)
		return err
	}

	if !(reservation.StartDate.After(availabilityPeriod.StartDate) && reservation.EndDate.Before(availabilityPeriod.EndDate)) {
		err := errors.New("reservation dates not in period")
		rr.logger.Println(err)
		return err
	}
	id, _ := gocql.RandomUUID()
	price := rr.calculatePrice(reservation.StartDate, reservation.EndDate, availabilityPeriod.IsPricePerGuest, availabilityPeriod.Price, reservation.GuestNum)
	err = rr.session.Query(
		`INSERT INTO reservations_by_availability_period (availability_period_id, reservation_id, start_date, end_date, accommodation_id, host_id, guest_id, guest_num, price) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		reservation.AvailabilityPeriodId, id, reservation.StartDate, reservation.EndDate, reservation.AccommodationId.Hex(), reservation.HostId.Hex(), reservation.GuestId.Hex(), reservation.GuestNum, price).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	err = rr.session.Query(
		`INSERT INTO reservations_by_guest (availability_period_id, reservation_id, start_date, end_date, accommodation_id, host_id, guest_id, guest_num, price) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		reservation.AvailabilityPeriodId, id, reservation.StartDate, reservation.EndDate, reservation.AccommodationId.Hex(), reservation.HostId.Hex(), reservation.GuestId.Hex(), reservation.GuestNum, price).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}

	err = rr.session.Query(
		`INSERT INTO reservations_by_dates (availability_period_id, reservation_id, start_date, end_date) 
		VALUES (?, ?, ?, ?)`,
		reservation.AvailabilityPeriodId, id, reservation.StartDate, reservation.EndDate).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
}

// function that checks if reservation dates are available, returns true if they are
func (rr *ReservationsRepo) checkReservationDates(reservation *ReservationByAvailabilityPeriod) (bool, error) {

	//check if reservation dates overlap with other reservation dates
	iter := rr.session.Query(`
        SELECT reservation_id FROM reservations_by_dates 
        WHERE availability_period_id = ? AND start_date < ? AND end_date > ? ALLOW FILTERING`,
		reservation.AvailabilityPeriodId, reservation.EndDate, reservation.StartDate).Iter()

	// Iterate over the result set to check if there are any rows
	for iter.Scan(nil) {
		return false, nil
	}

	if err := iter.Close(); err != nil {
		rr.logger.Println(err)
		return false, err
	}

	// If there are no rows, it means no overlap, so return true
	return true, nil
}

func (rr *ReservationsRepo) DeleteReservationByIdAndGuestId(id, guestId string) (*ReservationByAvailabilityPeriod, error) {
	reservation, err := rr.GetReservationByIdAndGuestId(id, guestId)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	if time.Now().After(reservation.StartDate) {
		return nil, errors.New("cannot delete reservation after start date has passed")
	}

	if err := rr.session.Query(`DELETE FROM reservations_by_availability_period WHERE reservation_id = ? AND availability_period_id = ?`, id, reservation.AvailabilityPeriodId).Exec(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	if err := rr.session.Query(`DELETE FROM reservations_by_guest WHERE reservation_id = ? AND guest_id = ?`, id, reservation.GuestId.Hex()).Exec(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	if err := rr.session.Query(`DELETE FROM reservations_by_dates WHERE reservation_id = ? AND availability_period_id = ? AND
                                        start_date = ? AND end_date = ?`, id, reservation.AvailabilityPeriodId, reservation.StartDate, reservation.EndDate).Exec(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservation, nil
}

func (rr *ReservationsRepo) calculatePrice(startDate time.Time, endDate time.Time, isPricePerGuest bool, price int, numberOfGuest int) int {
	reservationDuration := endDate.Sub(startDate)

	reservationDurationInDays := int(reservationDuration.Hours()) / 24

	if isPricePerGuest {
		return reservationDurationInDays * price * numberOfGuest
	}

	return reservationDurationInDays * price
}
