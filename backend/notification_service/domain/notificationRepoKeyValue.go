package domain

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"notification_service/config"
	"os"
	"time"
)

type NotificationsRepo struct {
	cli    *api.Client
	logger *config.Logger
}

// Constructs Redis Client
func NewNotificationsRepo(logger *config.Logger) (*NotificationsRepo, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &NotificationsRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// NoSQL: Saves Product to DB
func (nr *NotificationsRepo) Insert(notification *Notification) error {
	kv := nr.cli.KV()

	notification.Time = time.Now()

	dbId, id := generateKey(notification.Host.Id.Hex())
	notification.Id = id

	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	productKeyValue := &api.KVPair{Key: dbId, Value: data}
	_, err = kv.Put(productKeyValue, nil)
	if err != nil {
		return err
	}

	return nil
}

func (nr *NotificationsRepo) GetAll() (Notifications, error) {
	kv := nr.cli.KV()
	data, _, err := kv.List(notificationsAll, nil)
	if err != nil {
		return nil, err
	}

	notifications := Notifications{}
	for _, pair := range data {
		notification := &Notification{}
		err = json.Unmarshal(pair.Value, notification)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (nr *NotificationsRepo) GetNotificationsByUserID(userID string) (Notifications, error) {
	kv := nr.cli.KV()
	data, _, err := kv.List(constructKey(userID, ""), nil)
	if err != nil {
		return nil, err
	}

	notifications := Notifications{}
	for _, pair := range data {
		notification := &Notification{}
		err = json.Unmarshal(pair.Value, notification)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
