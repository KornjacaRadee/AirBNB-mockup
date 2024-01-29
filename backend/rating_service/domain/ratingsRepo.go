package domain

import (
	"fmt"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
)

type RatingsRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment and creates a keyspace
func New(logger *log.Logger) (*RatingsRepo, error) {
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
					}`, "ratings", 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	// Connect to reservations keyspace
	cluster.Keyspace = "ratings"
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &RatingsRepo{
		session: session,
		logger:  logger,
	}, nil
}

// Disconnect from database
func (rr *RatingsRepo) CloseSession() {
	rr.session.Close()
}

// Create tables
func (rr *RatingsRepo) CreateTables() {
	err := rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(host_id text, guest_id text, rating_id UUID, 
					time TIMESTAMP, rating int,
					PRIMARY KEY ((host_id), rating_id)) 
					WITH CLUSTERING ORDER BY (rating_id ASC)`,
			"host_ratings_by_host")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(host_id text, guest_id text, rating_id UUID, 
					time TIMESTAMP, rating int,
					PRIMARY KEY ((guest_id), rating_id)) 
					WITH CLUSTERING ORDER BY (rating_id ASC)`,
			"host_ratings_by_guest")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(host_id text, guest_id text, accommodation_id text, rating_id UUID, 
					time TIMESTAMP, rating int,
					PRIMARY KEY ((accommodation_id), rating_id)) 
					WITH CLUSTERING ORDER BY (rating_id ASC)`,
			"accommodation_ratings_by_accommodation")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(host_id text, guest_id text, accommodation_id text, rating_id UUID, 
					time TIMESTAMP, rating int,
					PRIMARY KEY ((host_id), rating_id)) 
					WITH CLUSTERING ORDER BY (rating_id ASC)`,
			"accommodation_ratings_by_host")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(host_id text, guest_id text, accommodation_id text, rating_id UUID, 
					time TIMESTAMP, rating int,
					PRIMARY KEY ((guest_id), rating_id)) 
					WITH CLUSTERING ORDER BY (rating_id ASC)`,
			"accommodation_ratings_by_guest")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}
}

func (rr *RatingsRepo) GetHostRatingsByHost(id string) (HostRatings, error) {
	scanner := rr.session.Query(`SELECT host_id, guest_id, rating_id, time, rating 
										FROM host_ratings_by_host WHERE host_id = ?`,
		id).Iter().Scanner()

	var ratings HostRatings
	for scanner.Next() {
		var rating HostRating
		var hostIdHex string
		var guestIdHex string

		err := scanner.Scan(&hostIdHex, &guestIdHex, &rating.Id, &rating.Time, &rating.Rating)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		hostId, err := primitive.ObjectIDFromHex(hostIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.HostId = hostId

		guestId, err := primitive.ObjectIDFromHex(guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.GuestId = guestId

		ratings = append(ratings, &rating)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return ratings, nil
}

func (rr *RatingsRepo) GetHostRatingsByGuest(id string) (HostRatings, error) {
	scanner := rr.session.Query(`SELECT host_id, guest_id, rating_id, time, rating 
										FROM host_ratings_by_guest WHERE guest_id = ?`,
		id).Iter().Scanner()

	var ratings HostRatings
	for scanner.Next() {
		var rating HostRating
		var hostIdHex string
		var guestIdHex string

		err := scanner.Scan(&hostIdHex, &guestIdHex, &rating.Id, &rating.Time, &rating.Rating)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		hostId, err := primitive.ObjectIDFromHex(hostIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.HostId = hostId

		guestId, err := primitive.ObjectIDFromHex(guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.GuestId = guestId

		ratings = append(ratings, &rating)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return ratings, nil
}

func (rr *RatingsRepo) GetHostRatingByIdAndHost(id string, hostId string) (*HostRating, error) {
	query := rr.session.Query(`SELECT host_id, guest_id, rating_id, time, rating 
										FROM host_ratings_by_host WHERE host_id = ? AND rating_id = ? LIMIT 1`, hostId, id)

	var rating HostRating
	var hostIdHex string
	var guestIdHex string

	err := query.Scan(&hostIdHex, &guestIdHex, &rating.Id, &rating.Time, &rating.Rating)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	newHostId, err := primitive.ObjectIDFromHex(hostIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	rating.HostId = newHostId

	guestId, err := primitive.ObjectIDFromHex(guestIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	rating.GuestId = guestId

	return &rating, nil
}

func (rr *RatingsRepo) GetHostRatingByIdAndGuest(id string, guestId string) (*HostRating, error) {
	query := rr.session.Query(`SELECT host_id, guest_id, rating_id, time, rating 
										FROM host_ratings_by_guest WHERE guest_id = ? AND rating_id = ? LIMIT 1`, guestId, id)

	var rating HostRating
	var hostIdHex string
	var guestIdHex string

	err := query.Scan(&hostIdHex, &guestIdHex, &rating.Id, &rating.Time, &rating.Rating)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	newHostId, err := primitive.ObjectIDFromHex(hostIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	rating.HostId = newHostId

	newGuestId, err := primitive.ObjectIDFromHex(guestIdHex)
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	rating.GuestId = newGuestId

	return &rating, nil
}

func (rr *RatingsRepo) InsertHostRating(rating *HostRating) error {
	id, _ := gocql.RandomUUID()
	err := rr.session.Query(
		`INSERT INTO host_ratings_by_guest (host_id, guest_id, rating_id, time, rating) 
		VALUES (?, ?, ?, ?, ?)`,
		rating.HostId.Hex(), rating.GuestId.Hex(), id, rating.Time, rating.Rating).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	err = rr.session.Query(
		`INSERT INTO host_ratings_by_host (host_id, guest_id, rating_id, time, rating) 
		VALUES (?, ?, ?, ?, ?)`,
		rating.HostId.Hex(), rating.GuestId.Hex(), id, rating.Time, rating.Rating).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
}

func (rr *RatingsRepo) GetAccommodationRatingsByAccommodation(id string) (AccommodationRatings, error) {
	scanner := rr.session.Query(`SELECT host_id, guest_id, accommodation_id, rating_id, time, rating 
										FROM accommodation_ratings_by_accommodation WHERE accommodation_id = ?`,
		id).Iter().Scanner()

	var ratings AccommodationRatings
	for scanner.Next() {
		var rating AccommodationRating
		var hostIdHex string
		var guestIdHex string
		var accommodationIdHex string

		err := scanner.Scan(&hostIdHex, &guestIdHex, &accommodationIdHex, &rating.Id, &rating.Time, &rating.Rating)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		hostId, err := primitive.ObjectIDFromHex(hostIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.HostId = hostId

		guestId, err := primitive.ObjectIDFromHex(guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.GuestId = guestId

		accommodationId, err := primitive.ObjectIDFromHex(accommodationIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.AccommodationId = accommodationId

		ratings = append(ratings, &rating)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return ratings, nil
}

func (rr *RatingsRepo) GetAccommodationRatingsByHost(id string) (AccommodationRatings, error) {
	scanner := rr.session.Query(`SELECT host_id, guest_id, accommodation_id, rating_id, time, rating 
										FROM accommodation_ratings_by_host WHERE host_id = ?`,
		id).Iter().Scanner()

	var ratings AccommodationRatings
	for scanner.Next() {
		var rating AccommodationRating
		var hostIdHex string
		var guestIdHex string
		var accommodationIdHex string

		err := scanner.Scan(&hostIdHex, &guestIdHex, &accommodationIdHex, &rating.Id, &rating.Time, &rating.Rating)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		hostId, err := primitive.ObjectIDFromHex(hostIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.HostId = hostId

		guestId, err := primitive.ObjectIDFromHex(guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.GuestId = guestId

		accommodationId, err := primitive.ObjectIDFromHex(accommodationIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.AccommodationId = accommodationId

		ratings = append(ratings, &rating)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return ratings, nil
}

func (rr *RatingsRepo) GetAccommodationRatingsByGuest(id string) (AccommodationRatings, error) {
	scanner := rr.session.Query(`SELECT host_id, guest_id, accommodation_id, rating_id, time, rating 
										FROM accommodation_ratings_by_guest WHERE guest_id = ?`,
		id).Iter().Scanner()

	var ratings AccommodationRatings
	for scanner.Next() {
		var rating AccommodationRating
		var hostIdHex string
		var guestIdHex string
		var accommodationIdHex string

		err := scanner.Scan(&hostIdHex, &guestIdHex, &accommodationIdHex, &rating.Id, &rating.Time, &rating.Rating)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		hostId, err := primitive.ObjectIDFromHex(hostIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.HostId = hostId

		guestId, err := primitive.ObjectIDFromHex(guestIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.GuestId = guestId

		accommodationId, err := primitive.ObjectIDFromHex(accommodationIdHex)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		rating.AccommodationId = accommodationId

		ratings = append(ratings, &rating)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return ratings, nil
}

func (rr *RatingsRepo) InsertAccommodationRating(rating *AccommodationRating) error {
	id, _ := gocql.RandomUUID()
	err := rr.session.Query(
		`INSERT INTO accommodation_ratings_by_accommodation (host_id, guest_id, accommodation_id, rating_id, time, rating) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		rating.HostId.Hex(), rating.GuestId.Hex(), rating.AccommodationId.Hex(), id, rating.Time, rating.Rating).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	err = rr.session.Query(
		`INSERT INTO accommodation_ratings_by_host (host_id, guest_id, accommodation_id, rating_id, time, rating) 
		VALUES (?, ?, ?, ?, ?)`,
		rating.HostId.Hex(), rating.GuestId.Hex(), rating.AccommodationId.Hex(), id, rating.Time, rating.Rating).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	err = rr.session.Query(
		`INSERT INTO accommodation_ratings_by_guest (host_id, guest_id, accommodation_id, rating_id, time, rating) 
		VALUES (?, ?, ?, ?, ?)`,
		rating.HostId.Hex(), rating.GuestId.Hex(), rating.AccommodationId.Hex(), id, rating.Time, rating.Rating).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
}
