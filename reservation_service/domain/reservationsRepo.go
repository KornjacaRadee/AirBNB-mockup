package domain

import (
	"fmt"
	"log"
	"os"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationsRepo struct {
	session *gocql.Session
	logger *log.Logger
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
		logger: logger,
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
					(accommodation_id text, availability_period_id UUID, start_date TIMESTAMP, end_date TIMESTAMP, price int, 
					PRIMARY KEY ((accommodation_id), availability_period_id)) 
					WITH CLUSTERING ORDER BY (availability_period_id ASC)`, 
					"availability_periods_by_accommodation")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(availability_period_id UUID, reservation_id UUID, start_date TIMESTAMP, end_date TIMESTAMP, guest_id text,
					PRIMARY KEY ((availability_period_id), reservation_id)) 
					WITH CLUSTERING ORDER BY (reservation_id ASC)`, 
					"reservations_by_availability_period")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}
}

func (rr *ReservationsRepo) GetAvailabilityPeriodsByAccommodation(id string) (AvailabilityPeriodsByAccommodation, error) {
	scanner := rr.session.Query(`SELECT accommodation_id, availability_period_id, start_date, end_date, price FROM availability_periods_by_accommodation WHERE accommodation_id = ?`,
					id).Iter().Scanner()
	
	var avaiabilityPeriods AvailabilityPeriodsByAccommodation
	for scanner.Next() {
		var avaiabilityPeriod AvailabilityPeriodByAccommodation
		var accommodationIdHex string
		err := scanner.Scan(&accommodationIdHex, &avaiabilityPeriod.Id, &avaiabilityPeriod.StartDate, &avaiabilityPeriod.EndDate, &avaiabilityPeriod.Price)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		accommodationId, err := primitive.ObjectIDFromHex(accommodationIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		avaiabilityPeriod.AccommodationId = accommodationId

		avaiabilityPeriods = append(avaiabilityPeriods, &avaiabilityPeriod)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return avaiabilityPeriods, nil
}

func (rr *ReservationsRepo) InsertAvailabilityPeriodByAccommodation(period *AvailabilityPeriodByAccommodation) (error) {
	err := rr.session.Query(
		`INSERT INTO availability_periods_by_accommodation (accommodation_id, availability_period_id, start_date, end_date, price) 
		VALUES (?, UUID(), ?, ?, ?)`,
		period.AccommodationId.Hex(), period.StartDate, period.EndDate, period.Price).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
}


func (rr *ReservationsRepo) GetReservationsByAvailabilityPeriod(id string) (ReservationsByAvailabilityPeriod, error) {
	scanner := rr.session.Query(`SELECT availability_period_id, reservation_id, start_date, end_date, guest_id FROM reservations_by_availability_period WHERE availability_period_id = ?`,
					id).Iter().Scanner()
	
	var reservations ReservationsByAvailabilityPeriod
	for scanner.Next() {
		var reservation ReservationByAvailabilityPeriod
		var guestIdHex string
		err := scanner.Scan(&reservation.AvailabilityPeriodId, &reservation.Id, &reservation.StartDate, &reservation.EndDate, &guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		guestId, err := primitive.ObjectIDFromHex(guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservation.GuestId = guestId

		reservations = append(reservations, &reservation)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservations, nil
}

func (rr *ReservationsRepo) InsertReservationByAvailabilityPeriod(reservation *ReservationByAvailabilityPeriod) (error) {
	err := rr.session.Query(
		`INSERT INTO reservations_by_availability_period (availability_period_id, reservation_id, start_date, end_date, guest_id) 
		VALUES (?, UUID(), ?, ?, ?)`,
		reservation.AvailabilityPeriodId, reservation.StartDate, reservation.EndDate, reservation.GuestId.Hex()).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
}