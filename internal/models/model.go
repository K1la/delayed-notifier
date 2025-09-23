package models

import (
	"time"
)

type Notification struct {
	ID        string             `json:"id"`
	Message   string             `json:"message"`
	Channel   Channel            `json:"channel"`
	To        string             `json:"to"`
	SendAt    time.Time          `json:"send_at"`
	Status    NotificationStatus `json:"status"`
	Retries   int                `json:"retries"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type Channel string

const (
	ChannelEmail    Channel = "email"
	ChannelTelegram Channel = "telegram"
)

type NotificationStatus string

const (
	StatusPending  NotificationStatus = "pending"
	StatusSent     NotificationStatus = "sent"
	StatusCanceled NotificationStatus = "canceled"
	StatusFailed   NotificationStatus = "failed"
)
