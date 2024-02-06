package domain

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	notification        = "host/%s/notification/%s"
	notificationForHost = "host/%s"
	notificationsAll    = "host"
)

func constructKey(hostId, notificationId string) string {
	if hostId != "" && notificationId != "" {
		return fmt.Sprintf(notification, hostId, notificationId)
	} else if hostId != "" && notificationId == "" {
		return fmt.Sprintf(notificationForHost, hostId)
	}
	return notificationsAll
}

func generateKey(hostId string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(notification, hostId, id), id
}
