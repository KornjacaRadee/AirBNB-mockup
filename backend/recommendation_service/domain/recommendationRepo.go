package domain

import (
	"context"
	// NoSQL: module containing Neo4J api client
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
)

type recommendationRepo struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func New(logger *log.Logger) (*recommendationRepo, error) {
	// Local instance
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &recommendationRepo{
		driver: driver,
		logger: logger,
	}, nil
}

// Check if connection is established
func (mr *recommendationRepo) CheckConnection() {
	ctx := context.Background()
	err := mr.driver.VerifyConnectivity(ctx)
	if err != nil {
		mr.logger.Panic(err)
		return
	}
	// Print Neo4J server address
	mr.logger.Printf(`Neo4J server address: %s`, mr.driver.Target().Host)
}
func (mr *recommendationRepo) CloseDriverConnection(ctx context.Context) {
	mr.driver.Close(ctx)
}
